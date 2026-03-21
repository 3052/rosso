package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/hulu"
   "flag"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/hulu.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringVar(&c.Job.PlayReady, "P", "PlayReady")
   //-------------------------------------------------------------
   threads := maya.IntVar(&c.Job.Threads, "t", "threads")
   //-------------------------------------------------------------
   email := maya.StringVar(&c.email, "e", "email")
   password := maya.StringVar(&c.password, "p", "password")
   //------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //---------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   if set[playReady] {
      return cache.Write(c)
   }
   if set[threads] {
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
      {playReady},
      {threads},
      {email, password},
      {address},
      {dash_id},
   })
}

func (c *client) do_address() error {
   err := c.Session.TokenRefresh()
   if err != nil {
      return err
   }
   deep_link, err := c.Session.DeepLink(hulu.Id(c.address))
   if err != nil {
      return err
   }
   c.Playlist, err = c.Session.Playlist(deep_link)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playlist.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Dash     *hulu.Dash
   Playlist *hulu.Playlist
   Session  *hulu.Session
   //--------------------
   Job maya.Job
   //--------------------
   email    string
   password string
   //--------------------
   address string
   //--------------------
   dash_id string
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4,*.mp4a")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email_password() error {
   var err error
   c.Session, err = hulu.FetchSession(c.email, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playlist.PlayReady,
   )
}
