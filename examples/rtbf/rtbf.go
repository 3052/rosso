package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/rtbf"
   "log"
)

func (c *client) do_dash() error {
   var (
      dash        maya.Dash
      entitlement rtbf.Entitlement
   )
   err := c.cache.Decode(&c.job, &dash, &entitlement)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, entitlement.FetchWidevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   job      maya.Job
   password string
}

func (c *client) do_email_password() error {
   account, err := rtbf.FetchAccount(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(account)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/rtbf"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
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
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   var account rtbf.Account
   err := c.cache.Decode(&account)
   if err != nil {
      return err
   }
   path, err := rtbf.GetPath(c.address)
   if err != nil {
      return err
   }
   asset_id, err := rtbf.FetchAssetId(path)
   if err != nil {
      return err
   }
   identity, err := account.Identity()
   if err != nil {
      return err
   }
   session, err := identity.Session()
   if err != nil {
      return err
   }
   entitlement, err := session.Entitlement(asset_id)
   if err != nil {
      return err
   }
   media, err := entitlement.GetDash()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(media)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, entitlement)
}
