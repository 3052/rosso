package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/molotov"); err != nil {
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
   return maya.FormatFlags(os.Stderr, "molotov", c)
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

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   auth, err := molotov.FetchAuth(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(auth)
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

func (c *client) do_dash_id() error {
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
