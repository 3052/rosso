package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "os"
   "path"
)

type client struct {
   cache maya.Cache

   address   maya.FlagString
   dash      maya.FlagString
   link_code maya.FlagBool
   proxy     maya.FlagString
   season    maya.FlagInt
   session   maya.FlagBool
   widevine  maya.FlagString
   mubi      maya.FlagInt
}

///

func (c *client) do() error {
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.widevine},
      {Name: "proxy", Value: &c.proxy},
      {Name: "link-code", Value: &c.link_code},
      {Name: "session", Value: &c.session},
      {Name: "address", Value: &c.address},
      {Name: "season", Value: &c.season, Needs: "address"},
      {Name: "mubi-id", Value: &c.mubi},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   if flags.IsSet(&c.widevine) {
      return c.cache.Encode(c.WidevineFolder)
   }
   if flags.IsSet(&c.proxy) {
      return c.cache.Encode(SetProxy(c.SetProxy.Value))
   }
   if c.link_code {
      return c.do_link_code()
   }
   if c.session {
      return c.do_session()
   }
   if flags.IsSet(&c.address) {
      if flags.IsSet(&c.season) {
         return c.do_address_season()
      }
      return c.do_address()
   }
   if flags.IsSet(&c.mubi) {
      return c.do_mubi()
   }
   if flags.IsSet(&c.dash) {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "mubi")
}

type SetProxy string

///

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
   slug := path.Base(c.Address.Value)
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
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &session, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: session.FetchWidevine,
   })
}

func (c *client) do_mubi() error {
   var session mubi.Session
   err := c.cache.Decode(&session)
   if err != nil {
      return err
   }
   err = session.FetchViewing(c.MubiId.Value)
   if err != nil {
      return err
   }
   secure_url, err := session.FetchSecureUrl(c.MubiId.Value)
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
   slug := path.Base(c.Address.Value)
   episodes, err := mubi.FetchEpisodes(slug, c.Season.Value)
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

func (c *client) do_use_proxy() error {
   var proxy SetProxy
   err := c.cache.Decode(&proxy)
   if err != nil {
      return err
   }
   return maya.SetProxy(string(proxy))
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
