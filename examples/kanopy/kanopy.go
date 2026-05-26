package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/kanopy"
   "fmt"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/kanopy"); err != nil {
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
   return flags.Usage(os.Stderr, "kanopy")
}

func (c *client) do_dash() error {
   var (
      login         kanopy.Login
      manifest      kanopy.Manifest
      maya_manifest maya.Manifest
   )
   err := c.cache.Decode(&login, &manifest, &maya_manifest)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return kanopy.CreateLicense(&login, &manifest, body)
   }
   return maya.DownloadDash(string(c.dash), &maya_manifest, &maya.Options{
      Device:  string(c.Widevine),
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

type client struct {
   Widevine maya.FlagString

   address  maya.FlagString
   dash     maya.FlagString
   email    maya.FlagString
   password maya.FlagString

   cache maya.Cache
}

func (c *client) do_email_password() error {
   login, err := kanopy.LoginUser(string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(login)
}

func (c *client) do_address() error {
   login := &kanopy.Login{}
   err := c.cache.Decode(login)
   if err != nil {
      return err
   }
   video, err := kanopy.ParseVideo(string(c.address))
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
