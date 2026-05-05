package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/rakuten.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "p", "PlayReady")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   season := maya.StringFlag(&c.season, "s", "season ID")
   //----------------------------------------------------------
   audio_language := maya.StringFlag(&c.AudioLanguage, "A", "audio language")
   episode := maya.StringFlag(&c.Episode, "e", "episode ID")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return c.run(c.do_season)
   case audio_language.IsSet:
      return c.run(c.do_language)
   case dash.IsSet:
      return c.run(c.do_language)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {address},
      {season},
      {audio_language, episode},
      {dash},
   })
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

func (c *client) do_language() error {
   var err error
   switch {
   case c.Url.IsMovie():
      c.StreamInfo, err = rakuten.FetchMovieStreaming(
         c.Url.ContentId, c.Start.Profile.Classification, c.AudioLanguage,
      )
   case c.Url.IsTvShow():
      c.StreamInfo, err = rakuten.FetchEpisodeStreaming(
         c.Episode, c.Start.Profile.Classification, c.AudioLanguage,
      )
   }
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.StreamInfo.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.StreamInfo.FetchLicense)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

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

func (c *client) do_address() error {
   var err error
   c.Url, err = rakuten.ParseUrl(c.address)
   if err != nil {
      return err
   }
   c.Start, err = rakuten.FetchStart(c.Url.MarketCode)
   if err != nil {
      return err
   }
   switch {
   case c.Url.IsMovie():
      movie, err := rakuten.FetchMovie(
         c.Url.ContentId, c.Start.Profile.Classification, c.Start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case c.Url.IsTvShow():
      show, err := rakuten.FetchTvShow(
         c.Url.ContentId, c.Start.Profile.Classification, c.Start.Market,
      )
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return cache.Write(c)
}

type client struct {
   // cache
   AudioLanguage string
   Dash          *maya.Dash
   Episode       string
   Job           maya.Job
   Start         *rakuten.Start
   StreamInfo    *rakuten.StreamInfo
   Url           *rakuten.ParsedUrl
   // flags
   address string
   season  string
   // state
   cache_err error
}
