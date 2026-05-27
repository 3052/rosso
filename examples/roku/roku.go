package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
   "os"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      playback roku.Playback
   )
   err := c.cache.Decode(&manifest, &playback)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: playback.LicenseWidevine,
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

   account_activation maya.FlagBool
   activation_status  maya.FlagBool
   dash               maya.FlagString
   roku_id            maya.FlagString
   use_account        maya.FlagBool

   cache maya.Cache
}

func (c *client) do_account_activation() error {
   account_token, err := roku.GetAccountToken(nil)
   if err != nil {
      return err
   }
   account_activation, err := roku.CreateAccountActivation(account_token)
   if err != nil {
      return err
   }
   fmt.Println(account_activation)
   return c.cache.Encode(account_activation, account_token)
}

func (c *client) do_activation_status() error {
   account_activation := &roku.AccountActivation{}
   account_token := &roku.AccountToken{}
   err := c.cache.Decode(account_activation, account_token)
   if err != nil {
      return err
   }
   activation_status, err := roku.GetActivationStatus(
      account_token, account_activation,
   )
   if err != nil {
      return err
   }
   return c.cache.Encode(activation_status)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/roku"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "account-activation", Value: &c.account_activation},
      {Name: "activation-status", Value: &c.activation_status},
      {Name: "roku-id", Value: &c.roku_id},
      {Name: "use-account", Value: &c.use_account, Needs: "roku-id"},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.account_activation {
      return c.do_account_activation()
   }
   if c.activation_status {
      return c.do_activation_status()
   }
   if c.roku_id != "" {
      return c.do_roku_id()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "roku")
}

func (c *client) do_roku_id() error {
   var status *roku.ActivationStatus
   if c.use_account {
      status = &roku.ActivationStatus{}
      err := c.cache.Decode(&status)
      if err != nil {
         return err
      }
   }
   account_token, err := roku.GetAccountToken(status)
   if err != nil {
      return err
   }
   playback, err := roku.GetPlayback(account_token, string(c.roku_id))
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&playback.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(account_token, manifest, playback)
}
