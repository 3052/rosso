package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/kanopy.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return c.run(c.do_address)
   }
   if dash.IsSet {
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, func(data []byte) ([]byte, error) {
      return kanopy.CreateLicense(c.Login, c.Manifest, data)
   })
}

func (c *client) do_email_password() error {
   var err error
   c.Login, err = kanopy.LoginUser(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   video, err := kanopy.ParseVideo(c.address)
   if err != nil {
      return err
   }
   if video.VideoId == 0 {
      video, err = kanopy.GetVideo(c.Login, video.Alias)
      if err != nil {
         return err
      }
   }
   memberships, err := kanopy.GetMemberships(c.Login)
   if err != nil {
      return err
   }
   play, err := kanopy.CreatePlay(c.Login, &memberships[0], video)
   if err != nil {
      return err
   }
   for _, caption := range play.Captions {
      for _, file := range caption.Files {
         fmt.Println(file.Url)
      }
   }
   c.Manifest, err = play.GetDashManifest()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Manifest.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   // cache
   Dash     *maya.Dash
   Job      maya.Job
   Login    *kanopy.LoginResponse
   Manifest *kanopy.Manifest
   // flags
   address  string
   email    string
   password string
   // state
   cache_err error
}
