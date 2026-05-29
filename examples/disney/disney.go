package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "fmt"
   "log"
   "os"
)

type client struct {
   Email     maya.FlagString
   PlayReady maya.FlagString
   Proxy     maya.FlagString
   address   maya.FlagString
   hls       maya.FlagString
   media     maya.FlagString
   passcode  maya.FlagString
   profile   maya.FlagString
   refresh   maya.FlagBool
   season    maya.FlagString
   threads   maya.FlagInt

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/disney"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "proxy", Value: &c.Proxy},
      {Name: "email", Value: &c.Email},
      {Name: "passcode", Value: &c.passcode},
      {Name: "profile-id", Value: &c.profile},
      {Name: "refresh", Value: &c.refresh},
      {Name: "address", Value: &c.address},
      {Name: "season-id", Value: &c.season},
      {Name: "media-id", Value: &c.media},
      {Name: "hls-id", Value: &c.hls},
      {Name: "threads", Value: &c.threads, Needs: "hls-id"},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.PlayReady) || flags.IsSet(&c.Proxy) {
      return c.cache.Encode(c)
   }
   if err := maya.SetProxy(string(c.Proxy)); err != nil {
      return err
   }
   if flags.IsSet(&c.Email) {
      return c.do_email()
   }
   if c.passcode != "" {
      return c.do_passcode()
   }
   if c.profile != "" {
      return c.do_profile()
   }
   if bool(c.refresh) {
      return c.do_refresh()
   }
   if c.address != "" {
      return c.do_address()
   }
   if c.season != "" {
      return c.do_season()
   }
   if c.media != "" {
      return c.do_media()
   }
   if c.hls != "" {
      return c.do_hls()
   }
   return flags.Usage(os.Stderr, "disney")
}

func (c *client) do_hls() error {
   var (
      manifest maya.Manifest
      token    disney.Token
   )
   err := c.cache.Decode(&manifest, &token)
   if err != nil {
      return err
   }
   return maya.DownloadHls(string(c.hls), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: token.FetchPlayReady,
      Threads: int(c.threads),
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_email() error {
   token, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := token.RequestOtp(string(c.Email))
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return c.cache.Encode(c, token)
}

func (c *client) do_passcode() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   otp, err := token.AuthenticateWithOtp(string(c.Email), string(c.passcode))
   if err != nil {
      return err
   }
   login, err := token.LoginWithActionGrant(otp.ActionGrant)
   if err != nil {
      return err
   }
   for i, profile := range login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return c.cache.Encode(token)
}

func (c *client) do_profile() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   err = token.SwitchProfile(string(c.profile))
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}

func (c *client) do_refresh() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   err = token.Refresh()
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}

func (c *client) do_address() error {
   entity_id, err := disney.GetEntityId(string(c.address))
   if err != nil {
      return err
   }
   entity, err := disney.GetEntityId(entity_id)
   if err != nil {
      return err
   }
   var token disney.Token
   if err = c.cache.Decode(&token); err != nil {
      return err
   }
   page, err := token.FetchPage(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *client) do_season() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   season, err := token.FetchSeason(string(c.season))
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_media() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   stream, err := token.FetchStream(string(c.media))
   if err != nil {
      return err
   }
   manifest, err := maya.ListHls(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
