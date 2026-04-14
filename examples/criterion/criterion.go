package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "log"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/criterion.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
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
   if widevine.IsSet {
      return cache.Write(c)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return with_cache(c.do_address)
   }
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},
      {address},
      {dash},
   })
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.File.FetchWidevine)
}

func (c *client) do_email_password() error {
   var err error
   c.Token, err = criterion.FetchToken(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   err := c.Token.Refresh()
   if err != nil {
      return err
   }
   files_href, err := criterion.FetchFilesHref(
      c.Token.AccessToken, path.Base(c.address),
   )
   if err != nil {
      return err
   }
   files, err := criterion.FetchFiles(c.Token.AccessToken, files_href)
   if err != nil {
      return err
   }
   c.File, err = criterion.GetDash(files)
   if err != nil {
      return err
   }
   dash, err := c.File.ParseDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

type client struct {
   File  *criterion.File
   Token *criterion.Token
   Dash  *maya.Dash
   //------------------------
   Job maya.Job
   //------------------------
   email    string
   password string
   //------------------------
   address string
}
