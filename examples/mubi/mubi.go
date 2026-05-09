package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "fmt"
   "log"
   "path"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address string
   cache   maya.Cache
   dash    string
   job     maya.Job
   mubi_id int
   season  int
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/mubi"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   code := maya.BoolFlag("c", "link code")
   mubi_id := maya.IntFlag(&c.mubi_id, "m", "Mubi ID")
   season := maya.IntFlag(&c.season, "s", "season")
   session := maya.BoolFlag("S", "session")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
   }
   if code.IsSet {
      return c.do_code()
   }
   if session.IsSet {
      return c.do_session()
   }
   if address.IsSet {
      return c.do_address()
   }
   if mubi_id.IsSet {
      return c.do_mubi_id()
   }
   if dash.IsSet {
      return c.do_dash()
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

///

func (c *client) do_code() error {
   var err error
   c.LinkCode, err = mubi.FetchLinkCode()
   if err != nil {
      return err
   }
   fmt.Println(c.LinkCode)
   return c.cache.Write(c)
}

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
   return c.cache.Write(c)
}

func (c *client) do_session() error {
   var err error
   c.Session, err = c.LinkCode.FetchSession()
   if err != nil {
      return err
   }
   return c.cache.Write(c)
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
