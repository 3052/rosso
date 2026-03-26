package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   film, err := mubi.ParseFilm(c.address)
   if err != nil {
      return err
   }
   if film.Id == 0 {
      err = film.FetchId()
      if err != nil {
         return err
      }
   }
   err = c.Session.Viewing(film.Id)
   if err != nil {
      return err
   }
   secure_url, err := c.Session.SecureUrl(film.Id)
   if err != nil {
      return err
   }
   c.Dash, err = secure_url.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Session.Widevine,
   )
}

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.dash")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash     *mubi.Dash
   LinkCode *mubi.LinkCode
   Session  *mubi.Session
   //--------------------
   Job maya.Job
   //--------------------
   address string
   //--------------------
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   code := maya.BoolVar(new(bool), "c", "link code")
   //----------------------------------------------------------
   session := maya.BoolVar(new(bool), "s", "session")
   //----------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[code]:
      return c.do_code()
   case set[session]:
      return with_cache(c.do_session)
   case set[address]:
      return with_cache(c.do_address)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {code},
      {session},
      {address},
      {dash_id},
   })
}

var cache maya.Cache

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
   c.Session, err = c.LinkCode.Session()
   if err != nil {
      return err
   }
   return cache.Write(c)
}
