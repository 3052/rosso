package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "log"
)

func (c *client) do_address() error {
   err := c.Device.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := c.Device.DeepLink(hulu.ParseId(c.address))
   if err != nil {
      return err
   }
   c.Playlist, err = c.Device.Playlist(deep_link.EabId)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playlist.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   var err error
   c.Device, err = hulu.FetchDevice(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Device   *hulu.Device
   //--------------------
   Job maya.Job
   //--------------------
   email    string
   password string
   //--------------------
   address string
   //--------------------
   dash_id string
}

var cache maya.Cache

func (c *client) do() error {
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "P", "PlayReady")
   //-------------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if playReady.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return with_cache(c.do_address)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
   )
}
