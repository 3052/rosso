package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

type client struct {
   cache           maya.Cache
   PlayReadyFolder maya.Flag[string]
   Username        maya.Flag[string] `depends:"Password"`
   Password        maya.Flag[string] `depends:"Username"`
   ProfileId       maya.Flag[string]
   Address         maya.Flag[string]
   DashId          maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/crave"); err != nil {
      return err
   }

   c.flag.AddValue(&c.PlayReadyFolder, "PR", "PlayReady")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.Password, "p", "password")
   c.flag.AddValue(&c.Username, "u", "username")
   c.flag = append(c.flag, nil)
   c.flag.AddValue(&c.ProfileId, "P", "profile")
   c.flag.AddValue(&c.Address, "a", "address")
   c.flag.AddValue(&c.DashId, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.PlayReadyFolder.Set {
      return c.cache.Encode(PlayReadyFolder(c.PlayReadyFolder.Value))
   }
   if c.Username.Set {
      if c.Password.Set {
         return c.do_username_password()
      }
   }
   if c.ProfileId.Set {
      return c.do_profile()
   }
   if c.Address.Set {
      return c.do_address()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   fmt.Println(c.flag)
   return nil
}

type PlayReadyFolder string

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_dash_id() error {
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
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: license,
   })
}

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
   if err = c.cache.Decode(profile_token); err != nil {
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
