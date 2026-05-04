package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

var cache maya.Cache

type client struct {
   // cache
   Job      maya.Job
   Dash     *maya.Dash
   Session  *mubi.Session
   LinkCode *mubi.LinkCode
   Proxy    string
   // flags
   address string
   season  int
   mubi_id int
   // state
   cache_err error
}

func (c *client) do() error {
   if err := cache.Setup("rosso/mubi.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
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
   if err := maya.ParseFlags(); err != nil {
      return err
   }

   if widevine.IsSet || proxy.IsSet {
      return cache.Write(c)
   }
   if code.IsSet {
      c.cache_err = nil
      return c.run(c.do_code)
   }
   if session.IsSet {
      return c.run(c.do_session)
   }
   if address.IsSet {
      c.cache_err = nil
      return c.run(c.do_address)
   }
   if mubi_id.IsSet {
      return c.run(c.do_mubi_id)
   }
   if dash.IsSet {
      return c.run(c.do_dash)
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

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   if err := maya.SetProxy(c.Proxy); err != nil {
      return err
   }
   return action()
}

// ----------------------------------------------------------------------
// Command Handlers
// ----------------------------------------------------------------------

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

// ----------------------------------------------------------------------
// Main
// ----------------------------------------------------------------------

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
