package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/peacock"
   "log"
   "os"
   "path"
)

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      playout  peacock.Playout
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &playout, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: playout.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   id_session, err := peacock.FetchIdSession(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(id_session)
}

func (c *client) do_address() error {
   id_session := &peacock.Cookie{}
   err := c.cache.Decode(id_session)
   if err != nil {
      return err
   }
   token, err := peacock.FetchToken(id_session)
   if err != nil {
      return err
   }
   playout, err := token.FetchPlayout(path.Base(c.Address.Value))
   if err != nil {
      return err
   }
   endpoint, err := playout.GetFastly()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(endpoint)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playout)
}

type WidevineFolder maya.Flag[string]

type client struct {
   cache          maya.Cache
   WidevineFolder WidevineFolder
   Email          maya.Flag[string] `depends:"Password"`
   Password       maya.Flag[string] `depends:"Email"`
   Address        maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/peacock"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.WidevineFolder.Set {
      return c.cache.Encode(c.WidevineFolder)
   }
   if c.Email.Set {
      if c.Password.Set {
         return c.do_email_password()
      }
   }
   if c.Address.Set {
      return c.do_address()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "peacock", c)
}
