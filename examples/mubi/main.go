package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/mubi"
   "flag"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/mubi.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringVar(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   code := maya.BoolVar(new(bool), "c", "link code")
   //----------------------------------------------------------
   session := maya.BoolVar(new(bool), "S", "session")
   //----------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   season := maya.IntVar(&c.season, "s", "season")
   //----------------------------------------------------------
   mubi_id := maya.IntVar(&c.mubi_id, "m", "Mubi ID")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   err = maya.SetProxy(c.Proxy, "*.dash")
   if err != nil {
      return err
   }
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[proxy]:
      return cache.Write(c)
   case set[code]:
      return c.do_code()
   case set[session]:
      return with_cache(c.do_session)
   case set[address]:
      return c.do_address()
   case set[mubi_id]:
      return with_cache(c.do_mubi_id)
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {proxy},
      {code},
      {session},
      {address, season},
      {mubi_id},
      {dash_id},
   })
}

var cache maya.Cache

type client struct {
   Dash     *mubi.Dash
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
   //--------------------
   dash_id string
}
