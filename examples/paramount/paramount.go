// WE COULD CACHE THE APP SECRET BUT THEN TO CHANGE LOCATION YOU WOULD NEED TO
// SET PROXY AND APP SECRET
package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "log"
   "net/http"
)

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.m4s,*.mp4")
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
   username string
   password string
   //--------------------
   cookie *maya.Flag
   //--------------------
   ParamountId string
   //--------------------
   dash_id string
}

///

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
   token, err := paramount.FetchStreamingUrl(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   c.Dash, err = token.Dash()
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
   token, err := paramount.FetchPlayReady(at, c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   return c.Job.DownloadDash(c.Dash.Body, c.Dash.Url, c.dash_id, token.Send)
}

func (c *client) do() error {
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
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
   if playReady.IsSet {
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
      {username, password},
      {paramount_id, c.cookie},
      {dash_id, c.cookie},
   })
}
