package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
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
      {Name: "threads", Value: &c.threads, Needs: "dash-id"},
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

func (c *client) do_dash() error {
   var (
      asset    molotov.Asset
      manifest maya.Manifest
   )
   err := c.cache.Decode(&asset, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: asset.FetchWidevine,
      Threads: int(c.threads),
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
   auth, err := molotov.FetchAuth(string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(auth)
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
   program, err := molotov.ParseProgram(string(c.address))
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

type client struct {
   Widevine maya.FlagString
   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString
   threads  maya.FlagInt

   cache maya.Cache
}
