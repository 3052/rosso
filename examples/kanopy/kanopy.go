package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
)

func (c *client) do_dash() error {
   var (
      login         kanopy.Login
      manifest      kanopy.Manifest
      maya_manifest maya.Manifest
      widevine      device
   )
   err := c.cache.Decode(&login, &manifest, &maya_manifest, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return kanopy.CreateLicense(&login, &manifest, body)
   }
   return maya.DownloadDash(c.dash, &maya_manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: license,
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
   login, err := kanopy.LoginUser(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   flag     maya.FlagSet
   password string
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/kanopy"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   email := c.flag.String(&c.email, "e", "email")
   password := c.flag.String(&c.password, "p", "password")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(device(c.widevine))
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
   maya_manifest, err := maya.ListDash(&manifest.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, maya_manifest)
}
