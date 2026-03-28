// WE COULD CACHE THE APP SECRET BUT THEN TO CHANGE LOCATION YOU WOULD NEED TO
// SET PROXY AND APP SECRET
package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
   "net/http"
)

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
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   token, err := paramount.PlayReady(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, token.Send)
}

func (c *client) do_username_password() error {
   app_secret, err := paramount.FetchAppSecret()
   if err != nil {
      return err
   }
   at, err := paramount.GetAt(app_secret)
   if err != nil {
      return err
   }
   c.CbsCom, err = paramount.FetchCbsCom(at, c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
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
   cookie      *maya.Flag
   //--------------------
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   proxy := maya.StringFlag(&c.Proxy, "x", "proxy")
   //--------------------------------------------------------------
   username := maya.StringFlag(&c.username, "U", "username")
   password := maya.StringFlag(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringFlag(&c.ParamountId, "p", "paramount ID")
   c.cookie = maya.BoolFlag("c", "cookie")
   //--------------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   err = maya.SetProxy(c.Proxy, "*.m4s,*.mp4")
   if err != nil {
      return err
   }
   if playReady.IsSet {
      return cache.Write(c)
   }
   if proxy.IsSet {
      return cache.Write(c)
   }
   if username.IsSet {
      if password.IsSet {
         return with_cache(c.do_username_password)
      }
   }
   if paramount_id.IsSet {
      return with_cache(c.do_paramount)
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {proxy},
      {username, password},
      {paramount_id, c.cookie},
      {dash_id, c.cookie},
   })
}

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
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   item, err := paramount.FetchItem(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   c.Dash, err = item.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}
