package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/rtbf.xml")
   if err != nil {
      return err
   }
   cache_err := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   var (
      action    func() error
      use_cache = true
   )
   switch {
   case widevine.IsSet:
      action = c.do_write
      use_cache = false
   case email.IsSet && password.IsSet:
      action = c.do_email_password
      use_cache = false
   case address.IsSet:
      action = c.do_address
   case dash.IsSet:
      action = c.do_dash_id
   }
   if action != nil {
      if use_cache && cache_err != nil {
         return cache_err
      }
      return action()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_write() error {
   return cache.Write(c)
}

func (c *client) do_address() error {
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   identity, err := c.Account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   c.Entitlement, err = session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   format, err := c.Entitlement.GetDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(format.GetManifest)
   if err != nil {
      return err
   }
   return c.do_write()
}

func (c *client) do_email_password() error {
   var err error
   c.Account, err = rtbf.FetchAccount(c.email, c.password)
   if err != nil {
      return err
   }
   return c.do_write()
}

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, c.Entitlement.FetchWidevine)
}

func main() {
   maya.SetProxy("")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   Account     *rtbf.Account
   Dash        *maya.Dash
   Entitlement *rtbf.Entitlement
   //---------------------------
   Job maya.Job
   //---------------------------
   email    string
   password string
   //---------------------------
   address string
}
