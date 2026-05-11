package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

///

func (c *client) do_dash() error {
   var (
      dash     maya.Dash
      login    kanopy.Login
      manifest kanopy.Manifest
   )
   err := c.cache.Decode(&c.job, &dash, &login, &manifest)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, func(data []byte) ([]byte, error) {
      return kanopy.CreateLicense(&login, &manifest, data)
   })
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

func (c *client) do_email_password() error {
   login, err := kanopy.LoginUser(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/kanopy"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
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
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   login := &kanopy.Login{}
   err := c.cache.Decode(login)
   if err != nil {
      return err
   }
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = kanopy.GetVideo(login, video.Alias)
      if err != nil {
         return err
      }
   }
   memberships, err := kanopy.GetMemberships(login)
   if err != nil {
      return err
   }
   play, err := kanopy.CreatePlay(login, &memberships[0], video)
   if err != nil {
      return err
   }
   for _, caption := range play.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   manifest, err := play.GetDash()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(&manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, manifest)
}
