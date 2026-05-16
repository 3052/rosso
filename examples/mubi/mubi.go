package main

import (
   "fmt"
   "log"
   "path"

   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(command).start()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) start() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   c.flag.AddValue(&c.proxy, "p", "Proxy")
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.Add(&c.code, "c", "link code")
   c.flag.Add(&c.session, "S", "session")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.season, "s", "season")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.mubi_id, "m", "Mubi ID")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }

   if c.proxy.Set {
      return c.cache.Encode(proxy_value(c.proxy.Value))
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine_value(c.widevine.Value))
   }

   if run := c.run(); run != nil {
      var proxy proxy_value
      if err := c.cache.Decode(&proxy); err != nil {
         return err
      }
      if err := maya.SetProxy(string(proxy)); err != nil {
         return err
      }
      return run()
   }

   fmt.Println(c.flag)
   return nil
}

func (c *command) run() func() error {
   if c.code.Set {
      return c.run_code
   }
   if c.session.Set {
      return c.run_session
   }
   if c.address.Set {
      if c.season.Set {
         return c.run_address_season
      }
      return c.run_address
   }
   if c.mubi_id.Set {
      return c.run_mubi_id
   }
   if c.dash.Set {
      return c.run_dash
   }
   return nil
}

type command struct {
   cache maya.Cache
   flag  maya.FlagSet

   address  maya.Flag
   code     maya.Flag
   dash     maya.Flag
   mubi_id  maya.Flag
   proxy    maya.Flag
   season   maya.Flag
   session  maya.Flag
   widevine maya.Flag
}

type (
   proxy_value    string
   widevine_value string
)

func (c *command) run_mubi_id() error {
   mubi_id, err := c.mubi_id.ParseInt()
   if err != nil {
      return err
   }
   var session mubi.Session
   if err := c.cache.Decode(&session); err != nil {
      return err
   }
   if err := session.FetchViewing(mubi_id); err != nil {
      return err
   }
   secure_url, err := session.FetchSecureUrl(mubi_id)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(secure_url.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *command) run_dash() error {
   var (
      manifest maya.Manifest
      session  mubi.Session
      widevine widevine_value
   )
   err := c.cache.Decode(&manifest, &session, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: session.FetchWidevine,
   })
}

func (c *command) run_code() error {
   link_code, err := mubi.FetchLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(link_code)
   return c.cache.Encode(link_code)
}

func (c *command) run_session() error {
   var link_code mubi.LinkCode
   if err := c.cache.Decode(&link_code); err != nil {
      return err
   }
   session, err := link_code.FetchSession()
   if err != nil {
      return err
   }
   return c.cache.Encode(session)
}

func (c *command) run_address() error {
   film, err := mubi.FetchFilm(path.Base(c.address.Value))
   if err != nil {
      return err
   }
   fmt.Println(film)
   return nil
}

func (c *command) run_address_season() error {
   season, err := c.season.ParseInt()
   if err != nil {
      return err
   }
   episodes, err := mubi.FetchEpisodes(path.Base(c.address.Value), season)
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i > 0 {
         fmt.Println()
      }
      fmt.Println(episode)
   }
   return nil
}
