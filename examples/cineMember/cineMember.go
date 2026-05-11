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
   flag     maya.FlagSet
   password string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/cineMember"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   password := c.flag.String(&c.password, "p", "password")
   email := c.flag.String(&c.email, "e", "email")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
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
   return maya.PrintFlags([]maya.FlagSet{
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
