package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "log"
   "os"
)

func (c *client) do_dash() error {
   var (
      entitlement rtbf.Entitlement
      manifest    maya.Manifest
   )
   err := c.cache.Decode(&entitlement, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: entitlement.FetchWidevine,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine maya.FlagString

   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
}

func (c *client) do_email_password() error {
   account, err := rtbf.FetchAccount(string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(account)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rtbf"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "rtbf")
}

func (c *client) do_address() error {
   var account rtbf.Account
   err := c.cache.Decode(&account)
   if err != nil {
      return err
   }
   path, err := rtbf.GetPath(string(c.address))
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
