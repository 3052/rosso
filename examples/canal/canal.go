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

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   dash      maya.Flag
   email     maya.Flag
   password  maya.Flag
   query     maya.Flag
   season    maya.Flag
   tracking  maya.Flag
   widevine  maya.Flag
   refresh   maya.Flag
   subtitles maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/canal"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.email, "e", "email")
   c.flag.AddValue(&c.password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.refresh, "r", "refresh")
   c.flag.AddValue(&c.query, "q", "query")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.tracking, "t", "tracking")
   c.flag.AddValue(&c.season, "s", "season")
   c.flag = append(c.flag, nil)
   c.flag.Add(&c.subtitles, "S", "subtitles")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine_device(c.widevine.Value))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if c.refresh.Set {
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
   if c.subtitles.Set {
      return c.do_subtitles()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_dash() error {
   var (
      device   widevine_device
      manifest maya.Manifest
      player   canal.Player
   )
   err := c.cache.Decode(&device, &manifest, &player)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Drm:     maya.DrmWidevine,
      Device:  string(device),
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

func (c *client) do_tracking_season() error {
   season, err := c.season.ParseInt()
   if err != nil {
      return err
   }
   var session canal.Session
   if err = c.cache.Decode(&session); err != nil {
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

type widevine_device string

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
