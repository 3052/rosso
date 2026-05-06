package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/canal"
   "fmt"
   "log"
   "net/url"
   "os"
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
   cache    maya.Cache
   email    string
   err      error
   job      maya.Job
   password string
   query    string
   season   int
   tracking string
}

func get(address string) error {
   target, err := url.Parse(address)
   if err != nil {
      return err
   }
   resp, err := maya.Get(target, nil)
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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/canal.xml"); err != nil {
      return err
   }
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   query := maya.StringFlag(&c.query, "q", "query")
   refresh := maya.BoolFlag("r", "refresh")
   season := maya.IntFlag(&c.season, "s", "season")
   subtitles := maya.BoolFlag("S", "subtitles")
   tracking := maya.StringFlag(&c.tracking, "t", "tracking")
   c.err = c.cache.Decode(&c.job)
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if refresh.IsSet {
      return c.do_refresh()
   }
   if query.IsSet {
      return c.do_query()
   }
   if tracking.IsSet {
      if season.IsSet {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   if subtitles.IsSet {
      return c.do_subtitles()
   }
   if dash.IsSet {
      return c.do_dash()
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

///

func (c *client) do_dash() error {
   return c.Dash.Download(&c.job, c.Player.FetchWidevine)
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
   return c.cache.Write(c)
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
   return c.cache.Write(c)
}

func (c *client) do_refresh() error {
   var err error
   c.Session, err = canal.FetchSession(c.Session.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Write(c)
}
