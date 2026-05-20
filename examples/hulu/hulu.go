package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
   "os"
)

func (c *client) do_dash_id() error {
   var (
      manifest  maya.Manifest
      playReady PlayReadyFolder
      playlist  hulu.Playlist
   )
   err := c.cache.Decode(&manifest, &playReady, &playlist)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
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
   deep_link, err := device.DeepLink(hulu.ParseId(c.Address.Value))
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

func (c *client) do_email_password() error {
   device, err := hulu.FetchDevice(c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(device)
}

type PlayReadyFolder maya.Flag[string]

type client struct {
   cache           maya.Cache
   PlayReadyFolder PlayReadyFolder
   Email           maya.Flag[string] `depends:"Password"`
   Password        maya.Flag[string] `depends:"Email"`
   Address         maya.Flag[string]
   DashId          maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.PlayReadyFolder.Set {
      return c.cache.Encode(c.PlayReadyFolder)
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
   return maya.FormatFlags(os.Stderr, "hulu", c)
}
