package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "path"
)

type WidevineFolder string

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: pluto.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

func (c *client) do_show() error {
   series, err := pluto.FetchSeries(path.Base(c.show.Value))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Encode(series)
}

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie.Value))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(series.GetMovieUrl())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_episode() error {
   var series pluto.Series
   err := c.cache.Decode(&series)
   if err != nil {
      return err
   }
   episode, err := series.GetEpisodeUrl(c.episode.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(episode)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   dash     maya.Flag
   episode  maya.Flag
   movie    maya.Flag
   show     maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.AddValue(&c.movie, "m", "movie URL")
   c.flag.AddValue(&c.show, "s", "show URL")
   c.flag.AddValue(&c.episode, "e", "episode ID")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.widevine.Set:
      return c.cache.Encode(WidevineFolder(c.widevine.Value))
   case c.movie.Set:
      return c.do_movie()
   case c.show.Set:
      return c.do_show()
   case c.episode.Set:
      return c.do_episode()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
