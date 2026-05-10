package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

func (c *client) do_audio() error {
   var (
      address     rakuten.Address
      start       rakuten.Start
      stream_info *rakuten.StreamInfo
   )
   err := c.cache.Decode(&address, &start)
   if err != nil {
      return err
   }
   switch {
   case address.IsMovie():
      stream_info, err = rakuten.FetchMovieStreaming(
         address.ContentId, start.Profile.Classification, c.audio,
      )
   case address.IsTvShow():
      stream_info, err = rakuten.FetchEpisodeStreaming(
         c.episode, start.Profile.Classification, c.audio,
      )
   }
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(&stream_info.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, stream_info)
}

func (c *client) do_dash() error {
   var (
      dash        maya.Dash
      stream_info rakuten.StreamInfo
   )
   err := c.cache.Decode(&c.job, &dash, &stream_info)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, stream_info.FetchLicense)
}

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

func (c *client) do_season() error {
   var start rakuten.Start
   err := c.cache.Decode(&start)
   if err != nil {
      return err
   }
   season, err := rakuten.FetchSeason(
      c.season, start.Profile.Classification, start.Market,
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
      return c.do_audio()
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
