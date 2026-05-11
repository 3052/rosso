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

func (c *client) do() error {
   if err := c.cache.Setup("rosso/canal"); err != nil {
      return err
   }
   email := c.flag.String(&c.email, "e", "email")
   password := c.flag.String(&c.password, "p", "password")
   query := c.flag.String(&c.query, "q", "query")
   refresh := c.flag.Bool("r", "refresh")
   season := c.flag.Int(&c.season, "s", "season")
   subtitles := c.flag.Bool("S", "subtitles")
   tracking := c.flag.String(&c.tracking, "t", "tracking")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(widevine_folder(c.widevine))
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
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {email, password},
      {refresh},
      {query},
      {tracking, season},
      {subtitles},
      {dash},
   })
}

func (c *client) do_tracking() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   player, err := session.Player(c.tracking)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&player.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, player)
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      player   canal.Player
      widevine widevine_folder
   )
   err := c.cache.Decode(&manifest, &player, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Drm:     maya.DrmWidevine,
      Device:  string(widevine),
      License: player.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
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

func (c *client) do_query() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   collections, err := session.Search(c.query)
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

func (c *client) do_tracking_season() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   episodes, err := session.Episodes(c.tracking, c.season)
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
   var player canal.Player
   err := c.cache.Decode(&player)
   if err != nil {
      return err
   }
   for _, subtitles := range player.Subtitles {
      err := get(subtitles.Url)
      if err != nil {
         return err
      }
   }
   return nil
}

type client struct {
   cache    maya.Cache
   dash     string
   email    string
   flag     maya.FlagSet
   password string
   query    string
   season   int
   tracking string
   widevine string
}

type widevine_folder string
