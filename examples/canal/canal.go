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

func (c *client) do() error {
   err := cache.Setup("rosso/canal.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   //----------------------------------------------------------
   threads := maya.IntFlag(&c.Job.Threads, "T", "threads")
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   refresh := maya.BoolFlag("r", "refresh")
   //---------------------------------------------------
   tracking := maya.StringFlag(&c.tracking, "t", "tracking")
   season := maya.IntFlag(&c.season, "s", "season")
   //----------------------------------------------------
   subtitles := maya.BoolFlag("S", "subtitles")
   //----------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if threads.IsSet {
      return cache.Write(c)
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
   if tracking.IsSet {
      if season.IsSet {
         return with_cache(c.do_tracking_season)
      }
      return with_cache(c.do_tracking)
   }
   if subtitles.IsSet {
      return with_cache(c.do_subtitles)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {threads},
      {widevine},
      {email, password},
      {refresh},
      {tracking, season},
      {subtitles},
      {dash_id},
   })
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

func (c *client) do_tracking() error {
   var err error
   c.Player, err = c.Session.Player(c.tracking)
   if err != nil {
      return err
   }
   c.Dash, err = c.Player.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Player.Widevine,
   )
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
func main() {
   maya.SetProxy("", "*.dash")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   Dash    *canal.Dash
   Player  *canal.Player
   Session *canal.Session
   //--------------------
   Job maya.Job
   //--------------------
   email    string
   password string
   //--------------------
   tracking string
   season   int
   //--------------------
   dash_id string
}
