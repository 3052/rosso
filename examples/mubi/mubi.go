package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      session  mubi.Session
      widevine device
   )
   err := c.cache.Decode(&manifest, &session, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(widevine),
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
   slug := path.Base(c.address)
   if c.season >= 1 {
      episodes, err := mubi.FetchEpisodes(slug, c.season)
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

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   mubi_id  int
   season   int
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   code := c.flag.Bool("c", "link code")
   mubi_id := c.flag.Int(&c.mubi_id, "m", "Mubi ID")
   season := c.flag.Int(&c.season, "s", "season")
   session := c.flag.Bool("S", "session")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(device(c.widevine))
   }
   if code.IsSet {
      return c.do_code()
   }
   if session.IsSet {
      return c.do_session()
   }
   if address.IsSet {
      return c.do_address()
   }
   if mubi_id.IsSet {
      return c.do_mubi_id()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {code},
      {session},
      {address, season},
      {mubi_id},
      {dash},
   })
}

func (c *client) do_mubi_id() error {
   var session mubi.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   err = session.FetchViewing(c.mubi_id)
   if err != nil {
      return err
   }
   secure_url, err := session.FetchSecureUrl(c.mubi_id)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(secure_url.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
