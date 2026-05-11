package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   email := c.flag.String(&c.email, "e", "email")
   password := c.flag.String(&c.password, "p", "password")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   playReady := c.flag.String(&c.playReady, "P", "PlayReady")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(device(c.playReady))
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
      {playReady},
      {email, password},
      {address},
      {dash},
   })
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
   deep_link, err := device.DeepLink(hulu.ParseId(c.address))
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
      manifest  maya.Manifest
      playReady device
      playlist  hulu.Playlist
   )
   err := c.cache.Decode(&manifest, &playReady, &playlist)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(playReady),
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
   device, err := hulu.FetchDevice(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(device)
}

type client struct {
   address   string
   cache     maya.Cache
   dash      string
   email     string
   flag      maya.FlagSet
   password  string
   playReady string
}

type device string
