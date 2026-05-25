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
   file, err := os.Create(path.Base(address.Path))
   if err != nil {
      return err
   }
   defer file.Close()
   _, err = file.ReadFrom(resp.Body)
   return err
}

type client struct {
   Widevine maya.FlagString

   email     maya.FlagString
   password  maya.FlagString
   refresh   maya.FlagBool
   query     maya.FlagString
   tracking  maya.FlagString
   season    maya.FlagInt
   subtitles maya.FlagBool
   dash      maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/canal"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "refresh", Value: &c.refresh},
      {Name: "query", Value: &c.query},
      {Name: "tracking", Value: &c.tracking},
      {Name: "season", Value: &c.season, Needs: "tracking"},
      {Name: "subtitles", Value: &c.subtitles},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.query != "" {
      return c.do_query()
   }
   if c.tracking != "" {
      if c.season >= 1 {
         return c.do_tracking_season()
      }
      return c.do_tracking()
   }
   if c.subtitles {
      return c.do_subtitles()
   }
   if c.dash != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "canal")
}

///

func (c *client) do_email_password() error {
   ticket, err := canal.FetchTicket()
   if err != nil {
      return err
   }
   login, err := ticket.Login(c.Email.Value, c.Password.Value)
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
   collections, err := session.Search(c.Query.Value)
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
   player, err := session.Player(c.Tracking.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&player.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, player)
}

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

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      player   canal.Player
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &player, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: player.FetchWidevine,
   })
}

func (c *client) do_tracking_season() error {
   var session canal.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   episodes, err := session.Episodes(c.Tracking.Value, c.Season.Value)
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
