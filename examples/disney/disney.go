package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   entity, err := disney.GetEntityId(address)
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
   season, err := token.FetchSeason(c.season.Value)
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
   stream, err := token.FetchStream(c.media.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListHls(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}

func (c *client) do_hls() error {
   var (
      device   playReady_device
      manifest maya.Manifest
      token    disney.Token
   )
   err := c.cache.Decode(&device, &manifest, &token)
   if err != nil {
      return err
   }
   return maya.DownloadHls(c.hls.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmPlayReady,
      License: token.FetchPlayReady,
   })
}

type client struct {
   cache     maya.Cache
   flag      maya.FlagSet
   address   maya.Flag
   email     maya.Flag
   hls       maya.Flag
   media     maya.Flag
   passcode  maya.Flag
   playReady maya.Flag
   profile   maya.Flag
   season    maya.Flag
   refresh   maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/disney"); err != nil {
      return err
   }
   c.flag.AddValue(&c.playReady, "PR", "PlayReady")
   c.flag.AddValue(&c.email, "e", "email")
   c.flag.AddValue(&c.passcode, "p", "passcode")
   c.flag.AddValue(&c.profile, "P", "profile ID")
   c.flag.Add(&c.refresh, "r", "refresh")
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.season, "s", "season ID")
   c.flag.AddValue(&c.media, "m", "media ID")
   c.flag.AddValue(&c.hls, "h", "HLS ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.playReady.Set:
      return c.cache.Encode(playReady_device(c.playReady.Value))
   case c.email.Set:
      return c.do_email()
   case c.passcode.Set:
      return c.do_passcode()
   case c.profile.Set:
      return c.do_profile()
   case c.refresh.Set:
      return c.do_refresh()
   case c.address.Set:
      return c.do_address()
   case c.season.Set:
      return c.do_season()
   case c.media.Set:
      return c.do_media()
   case c.hls.Set:
      return c.do_hls()
   }
   fmt.Println(c.flag)
   return nil
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type disney_email string

type playReady_device string

func (c *client) do_email() error {
   token, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := token.RequestOtp(c.email.Value)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return c.cache.Encode(disney_email(c.email.Value), token)
}

func (c *client) do_passcode() error {
   var (
      email disney_email
      token disney.Token
   )
   err := c.cache.Decode(&email, &token)
   if err != nil {
      return err
   }
   otp, err := token.AuthenticateWithOtp(string(email), c.passcode.Value)
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
   err = token.SwitchProfile(c.profile.Value)
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
