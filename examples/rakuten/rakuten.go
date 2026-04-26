package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   var err error
   c.Content, err = rakuten.ParseContent(c.address)
   if err != nil {
      return err
   }
   c.Classification, err = c.Content.FetchClassification()
   if err != nil {
      return err
   }
   switch {
   case c.Content.IsMovie():
      movie, err := c.Content.Movie(c.Classification)
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case c.Content.IsTvShow():
      show, err := c.Content.TvShow(c.Classification)
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return cache.Write(c)
}

func (c *client) do_season() error {
   season, err := c.Content.Season(c.Classification, c.season)
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

type client struct {
   Content        *rakuten.Content
   Dash           *maya.Dash
   StreamInfo     *rakuten.StreamInfo
   Classification *rakuten.Classification
   //-------------------
   Job maya.Job
   //-------------------
   address string
   //-------------------
   season string
   //-------------------
   Language string
   Episode  string
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.StreamInfo.FetchPlayReady)
}

func (c *client) do() error {
   err := cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "p", "PlayReady")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   season := maya.StringFlag(&c.season, "s", "season ID")
   //----------------------------------------------------------
   language := maya.StringFlag(&c.Language, "A", "audio language")
   episode := maya.StringFlag(&c.Episode, "e", "episode ID")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return with_cache(c.do_season)
   case language.IsSet:
      return with_cache(c.do_language)
   case dash.IsSet:
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {address},
      {season},
      {language, episode},
      {dash},
   })
}

func main() {
   maya.SetProxy("", "*.isma", "*.ismv")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_language() error {
   var err error
   c.StreamInfo, err = c.Content.FetchStreamInfo(
      c.Classification, c.Episode, c.Language, rakuten.PlayReady, rakuten.Uhd,
   )
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.StreamInfo.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache
