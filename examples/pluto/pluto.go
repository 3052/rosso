package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "os"
   "path"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_show_address() error {
   series, err := pluto.FetchSeries(path.Base(c.ShowAddress.Value))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Encode(series)
}

func (c *client) do_movie_address() error {
   series, err := pluto.FetchSeries(path.Base(c.MovieAddress.Value))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(series.GetMovieUrl())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_episode_id() error {
   var series pluto.Series
   err := c.cache.Decode(&series)
   if err != nil {
      return err
   }
   episode, err := series.GetEpisodeUrl(c.EpisodeId.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(episode)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

type WidevineFolder maya.Flag[string]

type client struct {
   cache          maya.Cache
   WidevineFolder WidevineFolder
   MovieAddress   maya.Flag[string]
   ShowAddress    maya.Flag[string]
   EpisodeId      maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.WidevineFolder.Set:
      return c.cache.Encode(c.WidevineFolder)
   case c.MovieAddress.Set:
      return c.do_movie_address()
   case c.ShowAddress.Set:
      return c.do_show_address()
   case c.EpisodeId.Set:
      return c.do_episode_id()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "pluto", c)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: pluto.FetchWidevine,
   })
}
