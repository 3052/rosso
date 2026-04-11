package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "fmt"
   "log"
)

func (c *client) do_address() error {
   entity, err := disney.ParseEntity(c.address)
   if err != nil {
      return err
   }
   page, err := c.Token.FetchPage(entity)
   if err != nil {
      return err
   }
   fmt.Println(page)
   return nil
}

func (c *client) do_season_id() error {
   season, err := c.Token.FetchSeason(c.season)
   if err != nil {
      return err
   }
   fmt.Println(season)
   return nil
}

func (c *client) do_media_id() error {
   stream, err := c.Token.FetchStream(c.media)
   if err != nil {
      return err
   }
   c.Hls, err = stream.FetchHls()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListHls(c.Hls.Body, c.Hls.Url)
}

func (c *client) do_refresh() error {
   err := c.Token.Refresh()
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_hls_id() error {
   return c.Job.DownloadHls(
      c.Hls.Body, c.Hls.Url, c.hls_id, c.Token.FetchPlayReady,
   )
}

func main() {
   maya.SetProxy("", "*.mp4", "*.mp4a")
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do() error {
   err := cache.Setup("rosso/disney.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   playReady := maya.StringFlag(&c.Job.PlayReady, "PR", "PlayReady")
   //--------------------------------------------------------------
   email := maya.StringFlag(&c.Email, "e", "email")
   //--------------------------------------------------------------
   passcode := maya.StringFlag(&c.passcode, "p", "passcode")
   //--------------------------------------------------------------
   profile := maya.StringFlag(&c.profile, "P", "profile ID")
   //--------------------------------------------------------------
   refresh := maya.BoolFlag("r", "refresh")
   //--------------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //--------------------------------------------------------------
   season := maya.StringFlag(&c.season, "s", "season ID")
   //--------------------------------------------------------------
   media := maya.StringFlag(&c.media, "m", "media ID")
   //--------------------------------------------------------------
   hls_id := maya.IntFlag(&c.hls_id, "h", "HLS ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return cache.Write(c)
   case email.IsSet:
      return c.do_email()
   case passcode.IsSet:
      return with_cache(c.do_passcode)
   case profile.IsSet:
      return with_cache(c.do_profile_id)
   case refresh.IsSet:
      return with_cache(c.do_refresh)
   case address.IsSet:
      return with_cache(c.do_address)
   case season.IsSet:
      return with_cache(c.do_season_id)
   case media.IsSet:
      return with_cache(c.do_media_id)
   case hls_id.IsSet:
      return with_cache(c.do_hls_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      playReady,
      email,
      passcode,
      profile,
      refresh,
      address,
      season,
      media,
      hls_id,
   }})
}

var cache maya.Cache

type client struct {
   Hls   *disney.Hls
   Token *disney.Token
   //-----------------
   Job maya.Job
   //-----------------
   Email string
   //-----------------
   passcode string
   //-----------------
   profile string
   //-----------------
   address string
   //-----------------
   season string
   //-----------------
   media string
   //-----------------
   hls_id int
}

func (c *client) do_email() error {
   var err error
   c.Token, err = disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := c.Token.RequestOtp(c.Email)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return cache.Write(c)
}

func (c *client) do_passcode() error {
   otp, err := c.Token.AuthenticateWithOtp(c.Email, c.passcode)
   if err != nil {
      return err
   }
   login, err := c.Token.LoginWithActionGrant(otp.ActionGrant)
   if err != nil {
      return err
   }
   for i, profile := range login.Account.Profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&profile)
   }
   return cache.Write(c)
}

func (c *client) do_profile_id() error {
   err := c.Token.SwitchProfile(c.profile)
   if err != nil {
      return err
   }
   return cache.Write(c)
}
