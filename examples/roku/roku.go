package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/roku"
   "fmt"
   "log"
)

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
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
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
   if dash.IsSet {
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {token},
      {set_code},
      {roku_id, c.get_code},
      {dash},
   })
}

func (c *client) do_dash_id() error {
   return c.Dash.Download(&c.Job, func(data []byte) ([]byte, error) {
      return roku.FetchWidevine(c.Playback.Drm.Widevine.LicenseServer, data)
   })
}

type client struct {
   Activation *roku.Activation
   Code       *roku.Code
   Token      *roku.Token
   Playback   *roku.Playback
   Dash       *maya.Dash
   //--------------------
   Job maya.Job
   //--------------------
   roku_id  string
   get_code *maya.Flag
}

func (c *client) do_roku_id() error {
   var code_token string
   if c.get_code.IsSet {
      code_token = c.Code.Token
   }
   var err error
   c.Token, err = roku.FetchToken(code_token)
   if err != nil {
      return err
   }
   c.Playback, err = roku.FetchPlayback(c.Token.AuthToken, c.roku_id)
   if err != nil {
      return err
   }
   dash, err := roku.ParseDash(c.Playback.Url)
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(dash)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_set_code() error {
   var err error
   c.Code, err = roku.FetchCode(c.Token.AuthToken, c.Activation.Code)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_token() error {
   var err error
   c.Token, err = roku.FetchToken("")
   if err != nil {
      return err
   }
   c.Activation, err = roku.FetchActivation(c.Token.AuthToken)
   if err != nil {
      return err
   }
   fmt.Println(roku.FormatActivation(c.Activation.Code))
   return cache.Write(c)
}

var cache maya.Cache

func main() {
   maya.SetProxy("", "*.mp4")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
