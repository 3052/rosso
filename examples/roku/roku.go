package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
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

///

type WidevineFolder maya.Flag[string]

type client struct {
   cache             maya.Cache
   WidevineFolder    WidevineFolder
   AccountActivation maya.Flag[bool]
   ActivationStatus  maya.Flag[bool]
   RokuId            maya.Flag[string]
   UseAccount        maya.Flag[bool] `depends:"RokuId"`
   DashId            maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/roku"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.WidevineFolder.Set {
      return c.cache.Encode(c.WidevineFolder)
   }
   if c.AccountActivation.Set {
      return c.do_account_activation()
   }
   if c.ActivationStatus.Set {
      return c.do_activation_status()
   }
   if c.RokuId.Set {
      return c.do_roku_id()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "roku", c)
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

func (c *client) do_roku_id() error {
   var status *roku.ActivationStatus
   if c.UseAccount.Set {
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
   playback, err := roku.GetPlayback(account_token, c.RokuId.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&playback.Url.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(account_token, manifest, playback)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      playback roku.Playback
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &playback, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: playback.LicenseWidevine,
   })
}
