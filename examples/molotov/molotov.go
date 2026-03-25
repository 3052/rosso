package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "flag"
   "log"
)

func (c *client) do_address() error {
   url, err := molotov.ParseUrl(c.address)
   if err != nil {
      return err
   }
   err = c.Login.Refresh()
   if err != nil {
      return err
   }
   program, err := url.FetchProgram(c.Login.Auth.AccessToken)
   if err != nil {
      return err
   }
   c.Asset, err = program.Asset(c.Login.Auth.AccessToken)
   if err != nil {
      return err
   }
   c.Dash, err = c.Asset.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Asset.Widevine,
   )
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do() error {
   err := cache.Setup("rosso/molotov.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[widevine] {
      return cache.Write(c)
   }
   if set[email] {
      if set[password] {
         return c.do_email_password()
      }
   }
   if set[address] {
      return with_cache(c.do_address)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {email, password},
      {address},
      {dash_id},
   })
}

type client struct {
   Asset *molotov.Asset
   Dash  *molotov.Dash
   Login *molotov.Login
   //------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
   //------------
   dash_id string
}

func (c *client) do_email_password() error {
   var err error
   c.Login, err = molotov.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
