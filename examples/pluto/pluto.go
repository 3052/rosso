package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "net/url"
   "path"
)

func main() {
   maya.SetProxy("", "*.m4s")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_show() error {
   var err error
   c.Series, err = pluto.FetchSeries(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&c.Series.Vod[0])
   return cache.Write(c)
}

type client struct {
   Series *pluto.Series
   Dash   *maya.Dash
   //------------------
   Job maya.Job
   //------------------
   movie string
   //------------------
   show string
   //------------------
   episode string
}

func (c *client) do() error {
   err := cache.Setup("rosso/pluto.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   movie := maya.StringFlag(&c.movie, "m", "movie URL")
   //----------------------------------------------------------
   show := maya.StringFlag(&c.show, "s", "show URL")
   //----------------------------------------------------------
   episode := maya.StringFlag(&c.episode, "e", "episode ID")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case movie.IsSet:
      return c.do_movie()
   case show.IsSet:
      return c.do_show()
   case episode.IsSet:
      return with_cache(c.do_episode)
   case dash.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      movie,
      show,
      episode,
      dash,
   }})
}

func (c *client) do_episode() error {
   var err error
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return c.Series.GetEpisodeUrl(c.episode)
   })
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie))
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return series.GetMovieUrl(), nil
   })
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, pluto.FetchWidevine)
}
