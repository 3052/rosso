package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "flag"
   "log"
   "path"
)

type client struct {
   Dash      *criterion.Dash
   File *criterion.File
   Token     *criterion.Token
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
func (c *client) do() error {
   err := cache.Setup("rosso/criterion.xml")
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
      c.Dash.Body, c.Dash.Url, c.dash_id, c.File.Widevine,
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
   item, err := c.Token.Item(path.Base(c.address))
   if err != nil {
      return err
   }
   files, err := c.Token.Files(item.Links.Files.Href)
   if err != nil {
      return err
   }
   c.File, err = files.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.File.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
