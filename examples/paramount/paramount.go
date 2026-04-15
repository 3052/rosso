package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/paramount"
   "fmt"
   "log"
   "net/http"
)

func (c *client) do_dash() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   var cbs_com *http.Cookie
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   session, err := app.FetchPlayReady(c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   return c.Dash.Download(&c.Job, session.Fetch)
}

var cache maya.Cache

type client struct {
   CbsCom *http.Cookie
   Dash   *maya.Dash
   //--------------------
   Job maya.Job
   //--------------------
   App string
   //--------------------
   username string
   password string
   //--------------------
   cookie *maya.Flag
   //--------------------
   ParamountId string
}

func (c *client) do_username_password() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   c.CbsCom, err = app.FetchCbsCom(c.username, c.password)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do() error {
   err := cache.Setup("rosso/paramount.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   app := maya.StringFlag(&c.App, "a", fmt.Sprint(paramount.GetAppKeys()))
   //--------------------------------------------------------------
   username := maya.StringFlag(&c.username, "U", "username")
   password := maya.StringFlag(&c.password, "P", "password")
   //--------------------------------------------------------------
   paramount_id := maya.StringFlag(&c.ParamountId, "p", "paramount ID")
   c.cookie = maya.BoolFlag("c", "cookie")
   //--------------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if playReady.IsSet {
      return cache.Write(c)
   }
   if app.IsSet {
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
   if dash.IsSet {
      return with_cache(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {playReady},
      {app},
      {username, password},
      {paramount_id, c.cookie},
      {dash, c.cookie},
   })
}

func main() {
   maya.SetProxy("", "*.m4s", "*.mp4")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_paramount() error {
   app, err := paramount.GetApp(c.App)
   if err != nil {
      return err
   }
   var cbs_com *http.Cookie
   if c.cookie.IsSet {
      cbs_com = c.CbsCom
   }
   session, err := app.FetchStreamingUrl(c.ParamountId, cbs_com)
   if err != nil {
      return err
   }
   dash, err := session.ParseDash()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
