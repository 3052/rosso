package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "net/url"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/pluto.xml")
   if err != nil {
      return err
   }
   cache_err := cache.Read(c)
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
   var (
      action    func() error
      use_cache = true
   )
   switch {
   case widevine.IsSet:
      action = c.do_write
      use_cache = false
   case movie.IsSet:
      action = c.do_movie
      use_cache = false
   case show.IsSet:
      action = c.do_show
      use_cache = false
   case episode.IsSet:
      action = c.do_episode
   case dash.IsSet:
      action = c.do_dash_id
   }
   if action != nil {
      if use_cache && cache_err != nil {
         return cache_err
      }
      return action()
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      movie,
      show,
      episode,
      dash,
   }})
}

func (c *client) do_write() error {
   return cache.Write(c)
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

func main() {
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
