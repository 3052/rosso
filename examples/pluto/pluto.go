package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "path"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      widevine device
   )
   err := c.cache.Decode(&manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
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

func (c *client) do_show() error {
   series, err := pluto.FetchSeries(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Encode(series)
}

type client struct {
   cache    maya.Cache
   dash     string
   episode  string
   flag     maya.FlagSet
   movie    string
   show     string
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto"); err != nil {
      return err
   }
   episode := c.flag.String(&c.episode, "e", "episode ID")
   movie := c.flag.String(&c.movie, "m", "movie URL")
   show := c.flag.String(&c.show, "s", "show URL")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(device(c.widevine))
   case movie.IsSet:
      return c.do_movie()
   case show.IsSet:
      return c.do_show()
   case episode.IsSet:
      return c.do_episode()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{{
      widevine,
      movie,
      show,
      episode,
      dash,
   }})
}

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie))
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
   episode, err := series.GetEpisodeUrl(c.episode)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(episode)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
