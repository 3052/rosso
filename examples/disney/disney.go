package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
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

///

type client struct {
   PlayReadyFolder maya.Flag[string]
   Email           maya.Flag[string]
   Passcode        maya.Flag[string]
   ProfileId       maya.Flag[string]
   Refresh         maya.Flag[bool]
   Address         maya.Flag[string]
   SeasonId        maya.Flag[string]
   MediaId         maya.Flag[string]
   HlsId           maya.Flag[string]

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/disney"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.PlayReadyFolder.Set:
      return c.cache.Encode(c.PlayReadyFolder)
   case c.Email.Set:
      return c.do_email()
   case c.Passcode.Set:
      return c.do_passcode()
   case c.ProfileId.Set:
      return c.do_profile_id()
   case c.Refresh.Set:
      return c.do_refresh()
   case c.Address.Set:
      return c.do_address()
   case c.SeasonId.Set:
      return c.do_season_id()
   case c.MediaId.Set:
      return c.do_media_id()
   case c.HlsId.Set:
      return c.do_hls_id()
   }
   return maya.FormatFlags(os.Stderr, "disney", c)
}

func (c *client) do_email() error {
   token, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := token.RequestOtp(c.Email.Value)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return c.cache.Encode(c.Email, token)
}

func (c *client) do_passcode() error {
   var (
      email_data Email
      token      disney.Token
   )
   err := c.cache.Decode(&email_data, &token)
   if err != nil {
      return err
   }
   otp, err := token.AuthenticateWithOtp(email_data.Value, c.Passcode.Value)
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

func (c *client) do_hls_id() error {
   var (
      manifest  maya.Manifest
      playReady PlayReadyFolder
      token     disney.Token
   )
   err := c.cache.Decode(&manifest, &playReady, &token)
   if err != nil {
      return err
   }
   return maya.DownloadHls(c.HlsId.Value, &manifest, &maya.Options{
      Device:  playReady.Value,
      Drm:     maya.DrmPlayReady,
      License: token.FetchPlayReady,
   })
}

func (c *client) do_profile_id() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   err = token.SwitchProfile(c.ProfileId.Value)
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
   entity_id, err := disney.GetEntityId(c.Address.Value)
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

func (c *client) do_season_id() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   season, err := token.FetchSeason(c.SeasonId.Value)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_media_id() error {
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   stream, err := token.FetchStream(c.MediaId.Value)
   if err != nil {
      return err
   }
   manifest, err := maya.ListHls(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest)
}
