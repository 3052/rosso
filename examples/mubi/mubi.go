package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

type client struct {
   cache maya.Cache
   flag  maya.FlagSet

   address  maya.Flag
   code     maya.Flag
   dash     maya.Flag
   mubi_id  maya.Flag
   season   maya.Flag
   session  maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
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
   if c.widevine.Set {
      return c.cache.Encode(widevine_device(c.widevine.Value))
   }
   if c.code.Set {
      return c.do_code()
   }
   if c.session.Set {
      return c.do_session()
   }
   if c.address.Set {
      return c.do_address()
   }
   if c.mubi_id.Set {
      return c.do_mubi_id()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_mubi_id() error {
   mubi_id, err := c.mubi_id.ParseInt()
   if err != nil {
      return err
   }
   var session mubi.Session
   if err = c.cache.Decode(&session); err != nil {
      return err
   }
   err = session.FetchViewing(mubi_id)
   if err != nil {
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

type widevine_device string

func (c *client) do_dash() error {
   var (
      device   widevine_device
      manifest maya.Manifest
      session  mubi.Session
   )
   err := c.cache.Decode(&device, &manifest, &session)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: session.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_code() error {
   link_code, err := mubi.FetchLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(link_code)
   return c.cache.Encode(link_code)
}

func (c *client) do_session() error {
   var link_code mubi.LinkCode
   err := c.cache.Decode(&link_code)
   if err != nil {
      return err
   }
   session, err := link_code.FetchSession()
   if err != nil {
      return err
   }
   return c.cache.Encode(session)
}

func (c *client) do_address() error {
   season, err := c.season.ParseInt()
   if err != nil {
      return err
   }
   slug := path.Base(c.address.Value)
   if season >= 1 {
      episodes, err := mubi.FetchEpisodes(slug, season)
      if err != nil {
         return err
      }
      for i, episode := range episodes {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(&episode)
      }
   } else {
      film, err := mubi.FetchFilm(slug)
      if err != nil {
         return err
      }
      fmt.Println(film)
   }
   return nil
}
