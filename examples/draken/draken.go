package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/draken"
   "flag"
   "log"
   "path"
)

func (c *client) do_address() error {
   viewer, err := draken.FetchViewer(path.Base(c.address))
   if err != nil {
      return err
   }
   entitlement, err := viewer.Entitlement(c.Login.Token)
   if err != nil {
      return err
   }
   c.Playback, err = viewer.Playback(c.Login.Token, entitlement.Token)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do() error {
   err := cache.Setup("rosso/draken.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   //----------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //----------------------------------------------
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

func (c *client) do_email_password() error {
   var err error
   c.Login, err = draken.FetchLogin(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id,
      func(data []byte) ([]byte, error) {
         return c.Playback.Widevine(c.Login.Token, data)
      },
   )
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

type client struct {
   Dash     *draken.Dash
   Login    *draken.Login
   Playback *draken.Playback
   //-----------------------
   Job maya.Job
   //-----------------------
   email    string
   password string
   //-----------------------
   address string
   //-----------------------
   dash_id string
}
