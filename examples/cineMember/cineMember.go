package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
   "net/http"
)

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
   c.Dash, err = maya.ListDash(link.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func main() {
   maya.SetProxy("", "*.m4s")
   log.SetFlags(log.Ltime)
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

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, nil)
}

type client struct {
   Dash    *maya.Dash
   Session *http.Cookie
   //---------------------
   Job maya.Job
   //-------------
   email    string
   password string
   //-------------
   address string
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
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
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
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {email, password},
      {address},
      {dash},
   })
}
