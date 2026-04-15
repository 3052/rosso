package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
   "net/url"
)

func (c *client) do_address() error {
   var err error
   c.User, err = plex.FetchUser()
   if err != nil {
      return err
   }
   path, err := plex.ParsePath(c.address)
   if err != nil {
      return err
   }
   metadata, err := plex.FetchMatch(c.User.AuthToken, path)
   if err != nil {
      return err
   }
   metadata, err = metadata.Fetch(c.User.AuthToken)
   if err != nil {
      return err
   }
   c.Part, err = metadata.GetDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return c.Part.GetManifest(c.User.AuthToken), nil
   })
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

type client struct {
   User *plex.User
   Part *plex.Part
   //------------------
   Dash *maya.Dash
   Job  maya.Job
   //------------------
   address string
}

func (c *client) do() error {
   err := cache.Setup("rosso/plex.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {address},
      {dash},
   })
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job,
      func(data []byte) ([]byte, error) {
         return c.Part.FetchWidevine(c.User.AuthToken, data)
      },
   )
}
