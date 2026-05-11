package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := c.cache.Setup("rosso/crave"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   password := c.flag.String(&c.password, "p", "password")
   profile := c.flag.String(&c.profile, "P", "profile")
   username := c.flag.String(&c.username, "u", "username")
   playReady := c.flag.String(&c.playReady, "PR", "PlayReady")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if playReady.IsSet {
      return c.cache.Encode(device(c.playReady))
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   if profile.IsSet {
      return c.do_profile()
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {playReady},
      {username, password},
      {profile},
      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   profile_token := &crave.ProfileToken{}
   err := c.cache.Decode(profile_token)
   if err != nil {
      return err
   }
   media, err := crave.ParseMedia(c.address)
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
      playReady     device
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
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
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

func (c *client) do_username_password() error {
   account_token, err := crave.PerformLogin(c.username, c.password)
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
   profile_token, err := crave.SwitchProfile(account_token, c.profile)
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
   address   string
   cache     maya.Cache
   dash      string
   flag      maya.FlagSet
   password  string
   playReady string
   profile   string
   username  string
}

type device string
