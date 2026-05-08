package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/disney"
   "fmt"
   "log"
)

func (c *client) do_hls() error {
   var (
      hls   maya.Hls
      token disney.Token
   )
   err := c.cache.Decode(&c.job, &hls, &token)
   if err != nil {
      return err
   }
   return hls.Download(c.hls, &c.job, token.FetchPlayReady)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type disney_email string

func (c *client) do_email() error {
   token, err := disney.RegisterDevice()
   if err != nil {
      return err
   }
   request_otp, err := token.RequestOtp(c.email)
   if err != nil {
      return err
   }
   fmt.Println(request_otp)
   return c.cache.Encode(disney_email(c.email), token)
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
   otp, err := token.AuthenticateWithOtp(string(email), c.passcode)
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
   err = token.SwitchProfile(c.profile)
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
   var token disney.Token
   err := c.cache.Decode(&token)
   if err != nil {
      return err
   }
   entity, err := disney.ParseEntity(c.address)
   if err != nil {
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
   season, err := token.FetchSeason(c.season)
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
   stream, err := token.FetchStream(c.media)
   if err != nil {
      return err
   }
   hls, err := maya.ListHls(stream)
   if err != nil {
      return err
   }
   return c.cache.Encode(hls)
}

type client struct {
   address  string
   cache    maya.Cache
   email    string
   hls      int
   job      maya.Job
   media    string
   passcode string
   profile  string
   season   string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/disney"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   media := maya.StringFlag(&c.media, "m", "media ID")
   passcode := maya.StringFlag(&c.passcode, "p", "passcode")
   profile := maya.StringFlag(&c.profile, "P", "profile ID")
   refresh := maya.BoolFlag("r", "refresh")
   season := maya.StringFlag(&c.season, "s", "season ID")
   email := maya.StringFlag(&c.email, "e", "email")
   playReady := maya.StringFlag(&c.job.PlayReady, "PR", "PlayReady")
   hls := maya.IntFlag(&c.hls, "h", "HLS ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case playReady.IsSet:
      return c.cache.Encode(c.job)
   case email.IsSet:
      return c.do_email()
   case passcode.IsSet:
      return c.do_passcode()
   case profile.IsSet:
      return c.do_profile()
   case refresh.IsSet:
      return c.do_refresh()
   case address.IsSet:
      return c.do_address()
   case season.IsSet:
      return c.do_season()
   case media.IsSet:
      return c.do_media()
   case hls.IsSet:
      return c.do_hls()
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
      hls,
   }})
}
