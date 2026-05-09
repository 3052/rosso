package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

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
   dash, err := maya.ListDash(&playlist.StreamUrl.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, playlist)
}

func (c *client) do_dash() error {
   var (
      dash     maya.Dash
      playlist hulu.Playlist
   )
   err := c.cache.Decode(&c.job, &dash, &playlist)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, playlist.FetchPlayReady)
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
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/hulu"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   playReady := maya.StringFlag(&c.job.PlayReady, "P", "PlayReady")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(c.job)
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
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {email, password},
      {address},
      {dash},
   })
}
