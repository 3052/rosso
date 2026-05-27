package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "os"
   "path"
)

func (c *client) do_dash() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
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

type client struct {
   Widevine maya.FlagString

   movie   maya.FlagString
   show    maya.FlagString
   episode maya.FlagString
   dash    maya.FlagString

   cache maya.Cache
}

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(
      path.Base(string(c.movie)),
   )
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(series.GetMovieUrl())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_show() error {
   series, err := pluto.FetchSeries(
      path.Base(string(c.show)),
   )
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Encode(series)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "movie-address", Value: &c.movie},
      {Name: "show-address", Value: &c.show},
      {Name: "episode-id", Value: &c.episode},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.Widevine):
      return c.cache.Encode(c)
   case c.movie != "":
      return c.do_movie()
   case c.show != "":
      return c.do_show()
   case c.episode != "":
      return c.do_episode()
   case c.dash != "":
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "pluto")
}

func (c *client) do_episode() error {
   var series pluto.Series
   err := c.cache.Decode(&series)
   if err != nil {
      return err
   }
   episode, err := series.GetEpisodeUrl(string(c.episode))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(episode)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
