package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/plex.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   xff := maya.StringFlag(&c.xff, "x", "x-forwarded-for")
   //----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case dash_id.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {address, xff},
      {dash_id},
   })
}

func (c *client) do_address() error {
   var err error
   c.User, err = plex.FetchUser()
   if err != nil {
      return err
   }
   address, err := plex.GetPath(c.address)
   if err != nil {
      return err
   }
   metadata, err := c.User.RatingKey(address)
   if err != nil {
      return err
   }
   metadata, err = c.User.Media(metadata, c.xff)
   if err != nil {
      return err
   }
   c.Part, err = metadata.Dash()
   if err != nil {
      return err
   }
   c.Dash, err = c.User.Dash(c.Part, c.xff)
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
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

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.User.Widevine(c.Part, data)
      },
   )
}

type client struct {
   Dash      *plex.Dash
   Part *plex.Part
   User      *plex.User
   //------------------
   Job maya.Job
   //------------------
   address string
   xff     string
   //------------------
   dash_id string
}
