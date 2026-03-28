package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "path"
)

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
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
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
   case dash_id.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      movie,
      show,
      episode,
      dash_id,
   }})
}

func (c *client) do_episode() error {
   url, err := c.Series.GetEpisodeUrl(c.episode)
   if err != nil {
      return err
   }
   c.Dash, err = pluto.FetchDash(url)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, pluto.Widevine)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie))
   if err != nil {
      return err
   }
   c.Dash, err = pluto.FetchDash(series.GetMovieUrl())
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

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
   Dash   *pluto.Dash
   Series *pluto.Series
   //------------------
   Job maya.Job
   //------------------
   movie string
   //------------------
   show string
   //------------------
   episode string
   //------------------
   dash_id string
}
