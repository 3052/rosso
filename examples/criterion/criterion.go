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
   threads := maya.IntFlag(&c.Job.Threads, "t", "threads")
   //----------------------------------------------------------
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if threads.IsSet {
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
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {threads},
      {email, password},
      {address},
      {dash_id},
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.File.FetchWidevine,
   )
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
   item, err := c.Token.FetchItem(path.Base(c.address))
   if err != nil {
      return err
   }
   files, err := c.Token.FetchFiles(item.Links.Files.Href)
   if err != nil {
      return err
   }
   c.File, err = files.GetDash()
   if err != nil {
      return err
   }
   c.Dash, err = c.File.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Dash  *criterion.Dash
   File  *criterion.File
   Token *criterion.Token
   //------------------------
   Job maya.Job
   //------------------------
   email    string
   password string
   //------------------------
   address string
   //------------------------
   dash_id string
}
