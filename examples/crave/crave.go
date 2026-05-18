package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
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

type PlayReadyFolder string

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
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: license,
   })
}

///

func (c *client) do_username_password() error {
   account_token, err := crave.PerformLogin(c.username.Value, c.password.Value)
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
   profile_token, err := crave.SwitchProfile(account_token, c.profile.Value)
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

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   address   maya.Flag
   dash      maya.Flag
   password  maya.Flag
   playReady maya.Flag
   profile   maya.Flag
   username  maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/crave"); err != nil {
      return err
   }
   c.flag.AddValue(&c.playReady, "PR", "PlayReady")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.password, "p", "password")
   c.flag.AddValue(&c.username, "u", "username")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.profile, "P", "profile")
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.playReady.Set {
      return c.cache.Encode(PlayReadyFolder(c.playReady.Value))
   }
   if c.username.Set {
      if c.password.Set {
         return c.do_username_password()
      }
   }
   if c.profile.Set {
      return c.do_profile()
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

func (c *client) do_address() error {
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   profile_token := &crave.ProfileToken{}
   if err = c.cache.Decode(profile_token); err != nil {
      return err
   }
   media, err := crave.ParseMedia(address)
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
