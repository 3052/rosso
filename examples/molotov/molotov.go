package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "fmt"
   "log"
)

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
   if err = c.cache.Decode(&auth); err != nil {
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
   if err := c.cache.Setup("rosso/molotov"); err != nil {
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
      return c.cache.Encode(widevine_device(c.widevine.Value))
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
