package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   auth, err := molotov.FetchAuth(c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(auth)
}

type widevine_device string

func (c *client) do_dash() error {
   var (
      asset    molotov.Asset
      manifest maya.Manifest
      device   widevine_device
   )
   err := c.cache.Decode(&asset, &manifest, &device)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: asset.FetchWidevine,
   })
}

func (c *client) do_address() error {
   input, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   address, err := molotov.ParseAddress(input)
   if err != nil {
      return err
   }
   var auth molotov.Auth
   err := c.cache.Decode(&auth)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   play, err := auth.FetchPlay(address)
   if err != nil {
      return err
   }
   asset, err := auth.FetchAsset(play)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(asset.GetManifest())
   if err != nil {
      return err
   }
   return c.cache.Encode(asset, auth, manifest)
}

///

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   flag     maya.FlagSet
   password string
   widevine string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/molotov"); err != nil {
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
      return c.cache.Encode(widevine_device(c.widevine))
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
