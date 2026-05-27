package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

type client struct {
   WidevineFolder maya.Flag[string]
   Email          maya.Flag[string] `depends:"Password"`
   Password       maya.Flag[string] `depends:"Email"`
   Address        maya.Flag[string]
   DashId         maya.Flag[string]

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rtbf"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
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
   return maya.FormatFlags(os.Stderr, "rtbf", c)
}

func (c *client) do_email_password() error {
   account, err := rtbf.FetchAccount(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(account)
}

func (c *client) do_address() error {
   var account rtbf.Account
   err := c.cache.Decode(&account)
   if err != nil {
      return err
   }
   path, err := rtbf.GetPath(c.Address.Value)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   identity, err := account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   entitlement, err := session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   media, err := entitlement.GetDash()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(media)
   if err != nil {
      return err
   }
   return c.cache.Encode(entitlement, manifest)
}

func (c *client) do_dash_id() error {
   var (
      entitlement rtbf.Entitlement
      manifest    maya.Manifest
      widevine    WidevineFolder
   )
   err := c.cache.Decode(&entitlement, &manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: entitlement.FetchWidevine,
   })
}
