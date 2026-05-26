package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "os"
   "path"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "proxy", Value: &c.Proxy},
      {Name: "link-code", Value: &c.link_code},
      {Name: "session", Value: &c.session},
      {Name: "address", Value: &c.address, Usage: "film or series URL"},
      {Name: "season", Value: &c.season, Needs: "address"},
      {Name: "mubi-id", Value: &c.mubi_id},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if flags.IsSet(&c.Proxy) {
      return c.cache.Encode(c)
   }
   if err := maya.SetProxy(string(c.Proxy)); err != nil {
      return err
   }
   if c.link_code {
      return c.do_link_code()
   }
   if c.session {
      return c.do_session()
   }
   if c.address != "" {
      if c.season >= 1 {
         return c.do_address_season()
      }
      return c.do_address()
   }
   if c.mubi_id >= 1 {
      return c.do_mubi()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "mubi")
}

func (c *client) do_link_code() error {
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
   slug := path.Base(string(c.address))
   film, err := mubi.FetchFilm(slug)
   if err != nil {
      return err
   }
   fmt.Println(film)
   return nil
}

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      session  mubi.Session
   )
   err := c.cache.Decode(&manifest, &session)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Drm:     maya.DrmWidevine,
      License: session.FetchWidevine,
      Device:  string(c.Widevine),
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_mubi() error {
   var session mubi.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   err = session.FetchViewing(int(c.mubi_id))
   if err != nil {
      return err
   }
   secure_url, err := session.FetchSecureUrl(int(c.mubi_id))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(secure_url.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_address_season() error {
   slug := path.Base(string(c.address))
   episodes, err := mubi.FetchEpisodes(slug, int(c.season))
   if err != nil {
      return err
   }
   for i, episode := range episodes {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(episode)
   }
   return nil
}

type client struct {
   Proxy    maya.FlagString
   Widevine maya.FlagString

   address   maya.FlagString
   dash      maya.FlagString
   link_code maya.FlagBool
   mubi_id   maya.FlagInt
   season    maya.FlagInt
   session   maya.FlagBool

   cache maya.Cache
}
