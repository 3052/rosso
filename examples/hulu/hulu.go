package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
   "os"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playlist hulu.Playlist
   )
   err := c.cache.Decode(&manifest, &playlist)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: playlist.FetchPlayReady,
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
   PlayReady maya.FlagString

   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
}

func (c *client) do_email_password() error {
   device, err := hulu.FetchDevice(string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(device)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.PlayReady) {
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
   return flags.Usage(os.Stderr, "hulu")
}

func (c *client) do_address() error {
   var device hulu.Device
   err := c.cache.Decode(&device)
   if err != nil {
      return err
   }
   err = device.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := device.DeepLink(
      hulu.ParseId(string(c.address)),
   )
   if err != nil {
      return err
   }
   playlist, err := device.Playlist(deep_link.EabId)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&playlist.StreamUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, playlist)
}
