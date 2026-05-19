package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "fmt"
   "log"
)

type WidevineFolder string

func (c *client) do_dash() error {
   var (
      entitlement rtbf.Entitlement
      manifest    maya.Manifest
      widevine    WidevineFolder
   )
   err := c.cache.Decode(&entitlement, &manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(widevine),
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

///

func (c *client) do_email_password() error {
   account, err := rtbf.FetchAccount(c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(account)
}

func (c *client) do_address() error {
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   var account rtbf.Account
   if err = c.cache.Decode(&account); err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(address.Path)
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

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   address  maya.Flag
   dash     maya.Flag
   email    maya.Flag
   password maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rtbf"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.email, "e", "email")
   c.flag.AddValue(&c.password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.widevine.Set {
      return c.cache.Encode(WidevineFolder(c.widevine.Value))
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if c.address.Set {
      return c.do_address()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
