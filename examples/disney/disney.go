package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "fmt"
   "log"
)

type client struct {
   cache           maya.Cache
   PlayReadyFolder maya.Flag[string]
   Email           maya.Flag[string]
   Passcode        maya.Flag[string]
   ProfileId       maya.Flag[string]
   Refresh         maya.Flag[bool]
   Address         maya.Flag[string]
   SeasonId        maya.Flag[string]
   MediaId         maya.Flag[string]
   HlsId           maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/disney"); err != nil {
      return err
   }

   c.flag.AddValue(&c.PlayReadyFolder, "PR", "PlayReady")
   c.flag.AddValue(&c.Email, "e", "email")
   c.flag.AddValue(&c.Passcode, "p", "passcode")
   c.flag.AddValue(&c.ProfileId, "P", "profile ID")
   c.flag.Add(&c.Refresh, "r", "refresh")
   c.flag.AddValue(&c.Address, "a", "address")
   c.flag.AddValue(&c.SeasonId, "s", "season ID")
   c.flag.AddValue(&c.MediaId, "m", "media ID")
   c.flag.AddValue(&c.HlsId, "h", "HLS ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.PlayReadyFolder.Set:
      return c.cache.Encode(PlayReadyFolder(c.PlayReadyFolder.Value))
   case c.Email.Set:
      return c.do_email()
   case c.Passcode.Set:
      return c.do_passcode()
   case c.ProfileId.Set:
      return c.do_profile()
   case c.Refresh.Set:
      return c.do_refresh()
   case c.Address.Set:
      return c.do_address()
   case c.SeasonId.Set:
      return c.do_season()
   case c.MediaId.Set:
      return c.do_media()
   case c.HlsId.Set:
      return c.do_hls()
   }
   fmt.Println(c.flag)
   return nil
}

type PlayReadyFolder string

func (c *client) do_address() error {
   entity_id, err := disney.GetEntity(c.Address.Value)
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
   season, err := token.FetchSeason(c.SeasonId.Value)
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

func (c *client) do_hls() error {
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
      Device:  string(playReady),
      Drm:     maya.DrmPlayReady,
      License: token.FetchPlayReady,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type EmailString string

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
   return c.cache.Encode(EmailString(c.Email.Value), token)
}

func (c *client) do_passcode() error {
   var (
      email EmailString
      token disney.Token
   )
   err := c.cache.Decode(&email, &token)
   if err != nil {
      return err
   }
   otp, err := token.AuthenticateWithOtp(string(email), c.Passcode.Value)
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
