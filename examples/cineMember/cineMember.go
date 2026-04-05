package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
   "net/http"
)

type client struct {
   Dash    *cineMember.Dash
   Session *http.Cookie
   //---------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
   //------------
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/cineMember.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
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
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_address() error {
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   stream, err := cineMember.FetchStream(c.Session, id)
   if err != nil {
      return err
   }
   link, err := stream.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = link.Dash()
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
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, nil)
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

func (c *client) do_email_password() error {
   var err error
   c.Session, err = cineMember.FetchSession()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(c.Session, c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
