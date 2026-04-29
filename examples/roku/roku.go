package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

func main() {
   maya.SetProxy("", "*.mp4")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_account_activation() error {
   var err error
   c.AccountToken, err = roku.GetAccountToken(nil)
   if err != nil {
      return err
   }
   c.AccountActivation, err = roku.CreateAccountActivation(c.AccountToken)
   if err != nil {
      return err
   }
   fmt.Println(c.AccountActivation)
   return cache.Write(c)
}

func (c *client) do_activation_status() error {
   var err error
   c.ActivationStatus, err = roku.GetActivationStatus(
      c.AccountToken, c.AccountActivation,
   )
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   Playback          *roku.Playback
   AccountToken      *roku.AccountToken
   ActivationStatus  *roku.ActivationStatus
   AccountActivation *roku.AccountActivation
   Dash              *maya.Dash
   //--------------------
   Job maya.Job
   //--------------------
   roku_id     string
   use_account *maya.Flag
}

func (c *client) do_roku_id() error {
   var status *roku.ActivationStatus
   if c.use_account.IsSet {
      status = c.ActivationStatus
   }
   var err error
   c.AccountToken, err = roku.GetAccountToken(status)
   if err != nil {
      return err
   }
   c.Playback, err = roku.GetPlayback(c.AccountToken, c.roku_id)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(c.Playback.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Playback.GetWidevineLicense)
}

func (c *client) do() error {
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   account_activation := maya.BoolFlag("a", "account activation")
   //----------------------------------------------------------
   activation_status := maya.BoolFlag("A", "activation status")
   //----------------------------------------------------------
   roku_id := maya.StringFlag(&c.roku_id, "r", "Roku ID")
   c.use_account = maya.BoolFlag("u", "use account")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if account_activation.IsSet {
      return c.do_account_activation()
   }
   if activation_status.IsSet {
      return with_cache(c.do_activation_status)
   }
   if roku_id.IsSet {
      if c.use_account.IsSet {
         return with_cache(c.do_roku_id)
      }
      return c.do_roku_id()
   }
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {account_activation},
      {activation_status},
      {roku_id, c.use_account},
      {dash},
   })
}
