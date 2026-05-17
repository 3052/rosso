package main

import (
   "fmt"
   "log"
   "path"

   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
)

type client struct {
   cache maya.Cache
   flag maya.FlagSet

   address   *maya.Flag
   dash      *maya.Flag
   link_code      *maya.Flag
   mubi_id   *maya.Flag
   set_proxy     *maya.Flag
   season    *maya.Flag
   session   *maya.Flag
   use_proxy *maya.Flag
   widevine  *maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }

   c.widevine = c.flag.AddGroup("widevine-folder", true, 1)
   c.set_proxy = c.flag.AddGroup("set-proxy", true, 1)
   c.use_proxy = c.flag.AddGroup("use-proxy", false, 1)
   c.link_code = c.flag.AddGroup("link-code", false, 1)
   c.session = c.flag.AddGroup("session", false, 1)
   
   c.address = c.flag.AddGroup("address", true, 2)
   c.season = c.flag.AddGroup("season", true, 2)
   
   c.mubi_id = c.flag.AddGroup("mubi-id", true, 3)
   c.dash = c.flag.AddGroup("dash-id", true, 3)

   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(widevine_folder(c.widevine.Value))
   }
   if c.set_proxy.Set {
      return c.cache.Encode(set_proxy(c.set_proxy.Value))
   }
   if c.use_proxy.Set {
      if err := c.do_use_proxy(); err != nil {
         return err
      }
   }
   if c.link_code.Set {
      return c.do_link_code()
   }
   if c.session.Set {
      return c.do_session()
   }
   if c.address.Set {
      if c.season.Set {
         return c.do_address_season()
      }
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

func (c *client) do_use_proxy() error {
   var proxy set_proxy
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

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      session  mubi.Session
      widevine widevine_folder
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
   slug := path.Base(c.address.Value)
   film, err := mubi.FetchFilm(slug)
   if err != nil {
      return err
   }
   fmt.Println(film)
   return nil
}

func (c *client) do_address_season() error {
   season, err := c.season.ParseInt()
   if err != nil {
      return err
   }
   slug := path.Base(c.address.Value)
   episodes, err := mubi.FetchEpisodes(slug, season)
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

type widevine_folder string

type set_proxy string
