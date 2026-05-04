package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

func (c *client) do() error {
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   cache_err := cache.Read(c)
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringFlag(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   code := maya.BoolFlag("c", "link code")
   //----------------------------------------------------------
   session := maya.BoolFlag("S", "session")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   season := maya.IntFlag(&c.season, "s", "season")
   //----------------------------------------------------------
   mubi_id := maya.IntFlag(&c.mubi_id, "m", "Mubi ID")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   var (
      action    func() error
      use_cache = true
   )
   switch {
   case widevine.IsSet, proxy.IsSet:
      action = c.do_write_cache
      use_cache = false
   case code.IsSet:
      action = c.do_code
      use_cache = false
   case session.IsSet:
      action = c.do_session
   case address.IsSet:
      action = c.do_address
      use_cache = false
   case mubi_id.IsSet:
      action = c.do_mubi_id
   case dash.IsSet:
      action = c.do_dash
   }
   if action != nil {
      if use_cache && cache_err != nil {
         return cache_err
      }
      if err := maya.SetProxy(c.Proxy); err != nil {
         return err
      }
      return action()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {proxy},
      {code},
      {session},
      {address, season},
      {mubi_id},
      {dash},
   })
}

func (c *client) do_write_cache() error {
   return cache.Write(c)
}

var cache maya.Cache

func (c *client) do_mubi_id() error {
   err := c.Session.FetchViewing(c.mubi_id)
   if err != nil {
      return err
   }
   secure_url, err := c.Session.FetchSecureUrl(c.mubi_id)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(secure_url.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_code() error {
   var err error
   c.LinkCode, err = mubi.FetchLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(c.LinkCode)
   return cache.Write(c)
}

func (c *client) do_session() error {
   var err error
   c.Session, err = c.LinkCode.FetchSession()
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_address() error {
   slug := path.Base(c.address)
   if c.season >= 1 {
      episodes, err := mubi.FetchEpisodes(slug, c.season)
      if err != nil {
         return err
      }
      for i, episode := range episodes {
         if i >= 1 {
            fmt.Println()
         }
         fmt.Println(&episode)
      }
   } else {
      film, err := mubi.FetchFilm(slug)
      if err != nil {
         return err
      }
      fmt.Println(film)
   }
   return nil
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.Session.FetchWidevine)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash     *maya.Dash
   LinkCode *mubi.LinkCode
   Session  *mubi.Session
   //--------------------
   Job maya.Job
   //--------------------
   Proxy string
   //--------------------
   address string
   season  int
   //--------------------
   mubi_id int
}
