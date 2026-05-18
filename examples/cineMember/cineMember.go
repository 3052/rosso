package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/cineMember"
   "log"
   "os"
)

type client struct {
   cache    maya.Cache
   Email    maya.Flag[string] `depends:"Password"`
   Password maya.Flag[string] `depends:"Email"`
   Address  maya.Flag[string]
   DashId   maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/cineMember"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   if c.Email.Set {
      if c.Password.Set {
         return c.do_email_password()
      }
   }
   if c.Address.Set {
      return c.do_address()
   }
   if c.DashId.Set {
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "cineMember", c)
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
   err = cineMember.FetchLogin(phpSessId, c.Email.Value, c.Password.Value)
   if err != nil {
      return err
   }
   return c.cache.Encode(phpSessId)
}

func (c *client) do_dash_id() error {
   var manifest maya.Manifest
   err := c.cache.Decode(&manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, nil)
}

func (c *client) do_address() error {
   phpSessId := &cineMember.Cookie{}
   err := c.cache.Decode(phpSessId)
   if err != nil {
      return err
   }
   id, err := cineMember.FetchId(c.Address.Value)
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
