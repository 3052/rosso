package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
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

type client struct {
   Widevine maya.FlagString

   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/molotov"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
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
   return flags.Usage(os.Stderr, "molotov")
}

///

func (c *client) do_email_password() error {
   auth, err := molotov.FetchAuth(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(auth)
}

func (c *client) do_dash() error {
   var (
      asset    molotov.Asset
      manifest maya.Manifest
      widevine WidevineFolder
   )
   err := c.cache.Decode(&asset, &manifest, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: asset.FetchWidevine,
   })
}

func (c *client) do_address() error {
   var auth molotov.Auth
   err := c.cache.Decode(&auth)
   if err != nil {
      return err
   }
   err = auth.Refresh()
   if err != nil {
      return err
   }
   program, err := molotov.ParseProgram(c.Address.Value)
   if err != nil {
      return err
   }
   play, err := auth.FetchPlay(program)
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
