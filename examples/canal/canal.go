package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/canal"
   "fmt"
   "log"
   "net/http"
   "os"
   "path"
)

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Player.FetchWidevine)
}

func (c *client) do_tracking() error {
   var err error
   c.Player, err = c.Session.Player(c.tracking)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Player.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_tracking_season() error {
   episodes, err := c.Session.Episodes(c.tracking, c.season)
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&episode)
   }
   return nil
}

func (c *client) do_subtitles() error {
   for _, subtitles := range c.Player.Subtitles {
      err := get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

func get(address string) error {
   resp, err := http.Get(address)
   if err != nil {
      return err
   }
   defer resp.Body.Close()
   file, err := os.Create(path.Base(address))
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   return err
}

var cache maya.Cache

type client struct {
   Dash    *maya.Dash
   Player  *canal.Player
   Session *canal.Session
   //--------------------
   Job maya.Job
   //--------------------
   email    string
   password string
   //--------------------
   query string
   //--------------------
   tracking string
   season   int
}

func (c *client) do() error {
   err := cache.Setup("rosso/canal.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   refresh := maya.BoolFlag("r", "refresh")
   //---------------------------------------------------
   query := maya.StringFlag(&c.query, "q", "query")
   //---------------------------------------------------
   tracking := maya.StringFlag(&c.tracking, "t", "tracking")
   season := maya.IntFlag(&c.season, "s", "season")
   //----------------------------------------------------
   subtitles := maya.BoolFlag("S", "subtitles")
   //----------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if refresh.IsSet {
      return with_cache(c.do_refresh)
   }
   if query.IsSet {
      return with_cache(c.do_query)
   }
   if tracking.IsSet {
      if season.IsSet {
         return with_cache(c.do_tracking_season)
      }
      return with_cache(c.do_tracking)
   }
   if subtitles.IsSet {
      return with_cache(c.do_subtitles)
   }
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {refresh},
      {query},
      {tracking, season},
      {subtitles},
      {dash},
   })
}

func main() {
   maya.SetProxy("", "*.dash")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_query() error {
   collections, err := c.Session.Search(c.query)
   if err != nil {
      return err
   }
   var line bool
   for _, collection := range collections {
      for _, asset := range collection.Assets {
         if line {
            fmt.Println()
         } else {
            line = true
         }
         fmt.Println(&asset)
      }
   }
   return nil
}

func (c *client) do_email_password() error {
   ticket, err := canal.FetchTicket()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.email, c.password)
   if err != nil {
      return err
   }
   c.Session, err = canal.FetchSession(login.SsoToken)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_refresh() error {
   var err error
   c.Session, err = canal.FetchSession(c.Session.SsoToken)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
