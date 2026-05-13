package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
)

func (c *client) do_address() error {
   address, err := c.address.Url()
   if err != nil {
      return err
   }
   phpSessId := &cineMember.Cookie{}
   if err = c.cache.Decode(phpSessId); err != nil {
      return err
   }
   id, err := cineMember.FetchId(address)
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

func (c *client) do_dash() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, nil)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   address  *maya.Flag
   dash     *maya.Flag
   email    *maya.Flag
   password *maya.Flag
}

func (c *client) do_email_password() error {
   phpSessId, err := cineMember.GetPhpSessId()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(phpSessId, c.email.Value, c.password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(phpSessId)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/cineMember"); err != nil {
      return err
   }
   c.address = c.flag.AddValue("a", "address")
   c.password = c.flag.AddValue("p", "password")
   c.email = c.flag.AddValue("e", "email")
   c.dash = c.flag.AddValue("d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   if c.email.Set {
      if c.password.Set {
         return c.do_email_password()
      }
   }
   if c.address.Set {
      return c.do_address()
   }
   if c.dash.Set {
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {c.email, c.password},
      {c.address},
      {c.dash},
   })
}
