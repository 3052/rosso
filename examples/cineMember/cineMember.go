package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   email    maya.FlagString
   password maya.FlagString
   address  maya.FlagString
   dash     maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   flags := maya.FlagSet{
      {Name: "email", Value: &c.email, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "email"},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if c.email != "" {
      if c.password != "" {
         return c.do_email_password()
      }
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "cineMember")
}

func (c *client) do_address() error {
   phpSessId := &cineMember.Cookie{}
   err := c.cache.Decode(phpSessId)
   if err != nil {
      return err
   }
   id, err := cineMember.FetchId(string(c.address))
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
   return maya.DownloadDash(string(c.dash), &manifest, nil)
}

func (c *client) do_email_password() error {
   phpSessId, err := cineMember.GetPhpSessId()
   if err != nil {
      return err
   }
   err = cineMember.FetchLogin(phpSessId, string(c.email), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(phpSessId)
}
