package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address string
   audio   string
   cache   maya.Cache
   dash    string
   episode string
   job     maya.Job
   season  string
}

func (c *client) do_address() error {
   address, err := rakuten.ParseAddress(c.address)
   if err != nil {
      return err
   }
   start, err := rakuten.FetchStart(address.MarketCode)
   if err != nil {
      return err
   }
   switch {
   case address.IsMovie():
      movie, err := rakuten.FetchMovie(
         address.ContentId, start.Profile.Classification, start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case address.IsTvShow():
      show, err := rakuten.FetchTvShow(
         address.ContentId, start.Profile.Classification, start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return c.cache.Encode(address, start)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rakuten"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   audio := maya.StringFlag(&c.audio, "A", "audio language")
   episode := maya.StringFlag(&c.episode, "e", "episode ID")
   season := maya.StringFlag(&c.season, "s", "season ID")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   playReady := maya.StringFlag(&c.job.PlayReady, "p", "PlayReady")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return c.cache.Encode(c.job)
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return c.do_season()
   case audio.IsSet:
      return c.do_language()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {address},

      {season},
      {audio, episode},
      {dash},
   })
}

///

func (c *client) do_season() error {
   season, err := rakuten.FetchSeason(
      c.season, c.Start.Profile.Classification, c.Start.Market,
   )
   if err != nil {
      return err
   }
   for i, episode := range season.Episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

func (c *client) do_language() error {
   var err error
   switch {
   case c.Url.IsMovie():
      c.StreamInfo, err = rakuten.FetchMovieStreaming(
         c.Url.ContentId, c.Start.Profile.Classification, c.audio,
      )
   case c.Url.IsTvShow():
      c.StreamInfo, err = rakuten.FetchEpisodeStreaming(
         c.episode, c.Start.Profile.Classification, c.audio,
      )
   }
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.StreamInfo.GetManifest)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.StreamInfo.FetchLicense)
}
