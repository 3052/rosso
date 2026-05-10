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
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache   maya.Cache
   episode string
   job     maya.Job
   movie   string
   show    string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/pluto.xml"); err != nil {
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

///

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
   return c.cache.Write(c)
}

func (c *client) do_episode() error {
   var err error
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return c.Series.GetEpisodeUrl(c.episode)
   })
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, pluto.FetchWidevine)
}

func (c *client) do_show() error {
   var err error
   c.Series, err = pluto.FetchSeries(path.Base(c.show))
   if err != nil {
      return err
   }
   fmt.Println(&c.Series.Vod[0])
   return c.cache.Write(c)
}
