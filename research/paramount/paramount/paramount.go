package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/research/paramount"
   "flag"
   "log"
   "net/http"
)

func (c *client) do_paramount() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   var cbs_com *http.Cookie
   if c.get_cookie {
      cbs_com = c.CbsCom
   }
   session, err := paramount.PlayReady(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   c.Dash, err = session.Dash()
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
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

type client struct {
   CbsCom *http.Cookie
   Dash   *paramount.Dash
   //--------------------
   Job maya.Job
   //--------------------
   Proxy string
   //--------------------
   username string
   password string
   //--------------------
   ParamountId string
   //--------------------
   dash_id    string
   get_cookie bool
}

func (c *client) do() error {
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringVar(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   threads := maya.IntVar(&c.Job.Threads, "t", "threads")
   //--------------------------------------------------------------
   proxy := maya.StringVar(&c.Proxy, "x", "proxy")
   //--------------------------------------------------------------
   username := maya.StringVar(&c.username, "U", "username")
   password := maya.StringVar(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringVar(&c.ParamountId, "p", "paramount ID")
   get_cookie := maya.BoolVar(&c.get_cookie, "c", "get cookie")
   //--------------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   err = maya.SetProxy(c.Proxy, "*.m4s,*.mp4")
   if err != nil {
      return err
   }
   if set[playReady] {
      return cache.Write(c)
   }
   if set[threads] {
      return cache.Write(c)
   }
   if set[proxy] {
      return cache.Write(c)
   }
   if set[paramount_id] {
      return with_cache(c.do_paramount)
   }
   if set[dash_id] {
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {playReady},
      {threads},
      {proxy},
      {username, password},
      {paramount_id, get_cookie},
      {dash_id, get_cookie},
   })
}

func (c *client) do_dash_id() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   var cbs_com *http.Cookie
   if c.get_cookie {
      cbs_com = c.CbsCom
   }
   token, err := paramount.PlayReady(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, token.Send)
}
