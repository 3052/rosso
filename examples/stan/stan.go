package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/stan"
   "fmt"
   "log"
   "os"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   PlayReady  maya.FlagString
   activation maya.FlagBool
   token      maya.FlagBool
   program_id maya.FlagInt
   quality    maya.FlagString
   dash_id    maya.FlagString

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/stan/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "activation", Value: &c.activation},
      {Name: "token", Value: &c.token},

      {Name: "program-id", Value: &c.program_id, Needs: "quality"},
      {
         Name:  "quality",
         Value: &c.quality,
         Needs: "program-id",
         Usage: "high ultra",
      },
      {Name: "dash-id", Value: &c.dash_id},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.PlayReady) {
      return c.cache.Encode(c)
   }
   if c.activation {
      return c.do_activation()
   }
   if c.token {
      return c.do_token()
   }
   if c.program_id >= 1 {
      if c.quality != "" {
         return c.do_program_id()
      }
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "stan")
}

func (c *client) do_activation() error {
   code, err := stan.FetchActivationCode()
   if err != nil {
      return err
   }
   fmt.Println(code)
   return c.cache.Encode(code)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      media    stan.Media
   )
   err := c.cache.Decode(&manifest, &media)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: media.LicensePlayReady,
   })
}

func (c *client) do_program_id() error {
   var token stan.WebToken
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   session, err := token.FetchSession(true)
   if err != nil {
      return err
   }
   media, err := session.FetchMedia(int(c.program_id), string(c.quality), "playready")
   if err != nil {
      return err
   }
   address, err := media.BaseUrl(stan.BaseUrl[0])
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(address)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, media)
}

func (c *client) do_token() error {
   var code stan.ActivationCode
   err := c.cache.Decode(&code)
   if err != nil {
      return err
   }
   token, err := code.FetchToken()
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}
