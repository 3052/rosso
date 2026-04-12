package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rakuten"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/rakuten.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   threads := maya.IntFlag(&c.Job.Threads, "t", "threads")
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringFlag(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   season := maya.StringFlag(&c.season, "s", "season ID")
   //----------------------------------------------------------
   language := maya.StringFlag(&c.Language, "A", "audio language")
   episode := maya.StringFlag(&c.Episode, "e", "episode ID")
   //----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   err = maya.SetProxy(c.Proxy, "*.isma", "*.ismv")
   if err != nil {
      return err
   }
   switch {
   case threads.IsSet:
      return cache.Write(c)
   case widevine.IsSet:
      return cache.Write(c)
   case proxy.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return with_cache(c.do_season)
   case language.IsSet:
      return with_cache(c.do_language)
   case dash_id.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {threads},
      {widevine},
      {proxy},
      {address},
      {season},
      {language, episode},
      {dash_id},
   })
}

var cache maya.Cache

func (c *client) do_address() error {
   var err error
   c.Content, err = rakuten.ParseContent(c.address)
   if err != nil {
      return err
   }
   switch {
   case c.Content.IsMovie():
      movie, err := c.Content.Movie()
      if err != nil {
         return err
      }
      fmt.Println(movie)
   case c.Content.IsTvShow():
      show, err := c.Content.TvShow()
      if err != nil {
         return err
      }
      fmt.Println(show)
   }
   return cache.Write(c)
}

func (c *client) do_season() error {
   season, err := c.Content.Season(c.season)
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
   stream, err := c.Content.Stream(
      c.Episode, c.Language, rakuten.Widevine, rakuten.Fhd,
   )
   if err != nil {
      return err
   }
   c.Dash, err = stream.Dash()
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
   stream, err := c.Content.Stream(
      c.Episode, c.Language, rakuten.Widevine, rakuten.Hd,
   )
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, stream.Widevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Content *rakuten.Content
   Dash    *rakuten.Dash
   //-------------------
   Job maya.Job
   //-------------------
   Proxy string
   //-------------------
   address string
   //-------------------
   season string
   //-------------------
   Language string
   Episode  string
   //-------------------
   dash_id string
}
