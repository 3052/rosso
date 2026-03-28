package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

func (c *client) do_roku_id() error {
   var code *roku.Code
   if c.get_code.IsSet {
      code = c.Code
   }
   var err error
   c.Token, err = roku.FetchToken(code)
   if err != nil {
      return err
   }
   c.Playback, err = c.Token.Playback(c.roku_id)
   if err != nil {
      return err
   }
   c.Dash, err = c.Playback.Dash()
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
      c.Dash.Body, c.Dash.Url, c.dash_id, c.Playback.Widevine,
   )
}

func (c *client) do_set_code() error {
   var err error
   c.Code, err = c.Token.Code(c.Activation)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_token() error {
   var err error
   c.Token, err = roku.FetchToken(nil)
   if err != nil {
      return err
   }
   c.Activation, err = c.Token.Activation()
   if err != nil {
      return err
   }
   fmt.Println(c.Activation)
   return cache.Write(c)
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "*.mp4")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Activation *roku.Activation
   Code       *roku.Code
   Dash       *roku.Dash
   Playback   *roku.Playback
   Token      *roku.Token
   //--------------------
   Job maya.Job
   //--------------------
   roku_id  string
   get_code *maya.Flag
   //--------------------
   dash_id string
}

func (c *client) do() error {
   err := cache.Setup("rosso/roku.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   token := maya.BoolFlag("t", "token")
   //----------------------------------------------------------
   set_code := maya.BoolFlag("s", "set code")
   //----------------------------------------------------------
   roku_id := maya.StringFlag(&c.roku_id, "r", "Roku ID")
   c.get_code = maya.BoolFlag("g", "get code")
   //----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if token.IsSet {
      return c.do_token()
   }
   if set_code.IsSet {
      return with_cache(c.do_set_code)
   }
   if roku_id.IsSet {
      if c.get_code.IsSet {
         return with_cache(c.do_roku_id)
      }
      return c.do_roku_id()
   }
   if dash_id.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {token},
      {set_code},
      {roku_id, c.get_code},
      {dash_id},
   })
}
