package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
)

func (c *client) do_dash() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   phpSessId, err := cineMember.GetPhpSessId()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(phpSessId, c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(phpSessId)
}

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   email    string
   password string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/cineMember"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   password := maya.StringFlag(&c.password, "p", "password")
   email := maya.StringFlag(&c.email, "e", "email")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
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
      {email, password},
      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   phpSessId := &cineMember.Cookie{}
   err := c.cache.Decode(phpSessId)
   if err != nil {
      return err
   }
   id, err := cineMember.FetchId(c.address)
   if err != nil {
      return err
   }
   stream, err := cineMember.FetchStream(phpSessId, id)
   if err != nil {
      return err
   }
   dash, err := stream.GetDash()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(dash)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
