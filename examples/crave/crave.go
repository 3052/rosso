package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
   "os"
)

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
   password maya.FlagString
   profile  maya.FlagString
   username maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/crave"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "username", Value: &c.username, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "username"},
      {Name: "profile-id", Value: &c.profile},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.PlayReady) {
      return c.cache.Encode(c)
   }
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.profile != "" {
      return c.do_profile()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "crave")
}

///

func (c *client) do_username_password() error {
   account_token, err := crave.PerformLogin(c.Username.Value, c.Password.Value)
   if err != nil {
      return err
   }
   profiles, err := crave.GetProfiles(account_token)
   if err != nil {
      return err
   }
   for i, profile := range profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return c.cache.Encode(account_token)
}

func (c *client) do_profile() error {
   account_token := &crave.AccountToken{}
   err := c.cache.Decode(account_token)
   if err != nil {
      return err
   }
   profile_token, err := crave.SwitchProfile(account_token, c.ProfileId.Value)
   if err != nil {
      return err
   }
   subs, err := crave.GetSubscriptions(profile_token)
   if err != nil {
      return err
   }
   for i, sub := range subs {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&sub)
   }
   return c.cache.Encode(profile_token)
}

func (c *client) do_address() error {
   profile_token := &crave.ProfileToken{}
   err := c.cache.Decode(profile_token)
   if err != nil {
      return err
   }
   media, err := crave.ParseMedia(c.Address.Value)
   if err != nil {
      return err
   }
   if media.FirstContent.Id == 0 {
      media, err = crave.GetMedia(media.Id)
      if err != nil {
         return err
      }
   }
   playback, err := crave.GetPlayback(profile_token, media)
   if err != nil {
      return err
   }
   stream, err := crave.GetStream(profile_token, playback)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, media, playback)
}

func (c *client) do_dash() error {
   var (
      manifest      maya.Manifest
      playReady     PlayReadyFolder
      playback      crave.Playback
      profile_token crave.ProfileToken
   )
   err := c.cache.Decode(&manifest, &playReady, &playback, &profile_token)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return crave.AcquireLicense(body, &profile_token, &playback)
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
      Drm:     maya.DrmPlayReady,
      License: license,
   })
}
