package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/pluto"
   "fmt"
   "log"
   "path"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto"); err != nil {
      return err
   }
   episode := maya.StringFlag(&c.episode, "e", "episode ID")
   movie := maya.StringFlag(&c.movie, "m", "movie URL")
   show := maya.StringFlag(&c.show, "s", "show URL")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(c.job)
   case movie.IsSet:
      return c.do_movie()
   case show.IsSet:
      return c.do_show()
   case episode.IsSet:
      return c.do_episode()
   case dash.IsSet:
      return c.do_dash()
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
   var series pluto.Series
   err := c.cache.Decode(&series)
   if err != nil {
      return err
   }
   episode, err := series.GetEpisodeUrl(c.episode)
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(episode)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash)
}

func (c *client) do_dash() error {
   var dash maya.Dash
   err := c.cache.Decode(&c.job, &dash)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, pluto.FetchWidevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache   maya.Cache
   dash    string
   episode string
   job     maya.Job
   movie   string
   show    string
}

func (c *client) do_movie() error {
   series, err := pluto.FetchSeries(path.Base(c.movie))
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(series.GetMovieUrl())
   if err != nil {
      return err
   }
   return c.cache.Encode(dash)
}

func (c *client) do_show() error {
   series, err := pluto.FetchSeries(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&series.Vod[0])
   return c.cache.Encode(series)
}
