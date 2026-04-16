package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

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
   maya.SetProxy("", "*.dash")
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
   address string
   season  int
   //--------------------
   mubi_id int
}

func (c *client) do() error {
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   //----------------------------------------------------------
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
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
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case code.IsSet:
      return c.do_code()
   case session.IsSet:
      return with_cache(c.do_session)
   case address.IsSet:
      return c.do_address()
   case mubi_id.IsSet:
      return with_cache(c.do_mubi_id)
   case dash.IsSet:
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {code},
      {session},
      {address, season},
      {mubi_id},
      {dash},
   })
}

var cache maya.Cache
