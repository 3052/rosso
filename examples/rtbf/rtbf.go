package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "log"
)

func (c *client) do_dash() error {
   var (
      entitlement rtbf.Entitlement
      manifest    maya.Manifest
      widevine    device
   )
   err := c.cache.Decode(&entitlement, &manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
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

func (c *client) do_email_password() error {
   account, err := rtbf.FetchAccount(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(account)
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   flag     maya.FlagSet
   password string
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rtbf"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   email := c.flag.String(&c.email, "e", "email")
   password := c.flag.String(&c.password, "p", "password")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(device(c.widevine))
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   var account rtbf.Account
   err := c.cache.Decode(&account)
   if err != nil {
      return err
   }
   path, err := rtbf.GetPath(c.address)
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
