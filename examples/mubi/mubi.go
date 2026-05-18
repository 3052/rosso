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
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   SetProxy       maya.Flag[string]
   LinkCode       maya.Flag[bool]
   Session        maya.Flag[bool]
   Address        maya.Flag[string]
   Season         maya.Flag[int] `depends:"Address"`
   MubiId         maya.Flag[int]
   UseProxy       maya.Flag[bool] `depends:"MubiId"`
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.WidevineFolder.Set {
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   }
   if c.SetProxy.Set {
      return c.cache.Encode(SetProxy(c.SetProxy.Value))
   }
   if c.UseProxy.Set {
      if err := c.do_use_proxy(); err != nil {
         return err
      }
   }
   if c.LinkCode.Set {
      return c.do_link_code()
   }
   if c.Session.Set {
      return c.do_session()
   }
   if c.Address.Set {
      if c.Season.Set {
         return c.do_address_season()
      }
      return c.do_address()
   }
   if c.MubiId.Set {
      return c.do_mubi_id()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "mubi", c)
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

func (c *client) do_dash_id() error {
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
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: session.FetchWidevine,
   })
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
   slug := path.Base(c.Address.Value)
   film, err := mubi.FetchFilm(slug)
   if err != nil {
      return err
   }
   fmt.Println(film)
   return nil
}

func (c *client) do_mubi_id() error {
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

type WidevineFolder string

type SetProxy string
