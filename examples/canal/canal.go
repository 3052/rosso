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

func get(address *url.URL) error {
   resp, err := maya.Get(address, nil)
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

type client struct {
   cache    maya.Cache
   dash     *maya.Flag
   email    *maya.Flag
   password *maya.Flag
   query    *maya.Flag
   season   *maya.Flag
   tracking *maya.Flag
   widevine *maya.Flag
   flag     maya.FlagSet
}

type widevine string

func (c *client) do_email_password() error {
   ticket, err := canal.FetchTicket()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   session, err := canal.FetchSession(login.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Encode(session)
}

func (c *client) do_refresh() error {
   session := &canal.Session{}
   err := c.cache.Decode(session)
   if err != nil {
      return err
   }
   session, err = canal.FetchSession(session.SsoToken)
   if err != nil {
      return err
   }
   return c.cache.Encode(session)
}

func (c *client) do_query() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   collections, err := session.Search(c.query.Value)
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

func (c *client) do_tracking() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   player, err := session.Player(c.tracking.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&player.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, player)
}

func (c *client) do_tracking_season() error {
   season, err := c.season.Int()
   if err != nil {
      return err
   }
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   episodes, err := session.Episodes(c.tracking.Value, season)
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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/canal"); err != nil {
      return err
   }
   c.dash = c.flag.AddValue("d", "DASH ID")
   c.email = c.flag.AddValue("e", "email")
   c.password = c.flag.AddValue("p", "password")
   c.query = c.flag.AddValue("q", "query")
   c.season = c.flag.AddValue("s", "season")
   c.tracking = c.flag.AddValue("t", "tracking")
   c.widevine = c.flag.AddValue("w", "Widevine")
   refresh := c.flag.Add("r", "refresh")
   subtitles := c.flag.Add("S", "subtitles")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine(c.widevine.Value))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if refresh.Set {
      return c.do_refresh()
   }
   if c.query.Set {
      return c.do_query()
   }
   if c.tracking.Set {
      if c.season.Set {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   if subtitles.Set {
      return c.do_subtitles()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {c.widevine},
      {c.email, c.password},
      {refresh},
      {c.query},
      {c.tracking, c.season},

      {subtitles},
      {c.dash},
   })
}

///

func (c *client) do_subtitles() error {
   var player canal.Player
   err := c.cache.Decode(&player)
   if err != nil {
      return err
   }
   for _, subtitles := range player.Subtitles {
      err := get(&subtitles.Url.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      player   canal.Player
      device   widevine
   )
   err := c.cache.Decode(&manifest, &player, &device)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Drm:     maya.DrmWidevine,
      Device:  string(device),
      License: player.FetchWidevine,
   })
}
