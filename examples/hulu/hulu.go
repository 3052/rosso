package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "fmt"
   "log"
)

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   address   maya.Flag
   dash      maya.Flag
   email     maya.Flag
   password  maya.Flag
   playReady maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   c.flag.AddValue(&c.playReady, "P", "PlayReady")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.email, "e", "email")
   c.flag.AddValue(&c.password, "p", "password")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.playReady.Set {
      return c.cache.Encode(playReady_device(c.playReady.Value))
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

type playReady_device string

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
   deep_link, err := device.DeepLink(hulu.ParseId(c.address.Value))
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

func (c *client) do_dash() error {
   var (
      device   playReady_device
      manifest maya.Manifest
      playlist hulu.Playlist
   )
   err := c.cache.Decode(&device, &manifest, &playlist)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
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

func (c *client) do_email_password() error {
   device, err := hulu.FetchDevice(c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(device)
}
