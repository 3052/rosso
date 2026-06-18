package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
   "fmt"
   "log"
   "os"
)

func (c *client) do_title_id() error {
   token_pair := &amazon.TokenPair{}
   err := c.cache.Decode(token_pair)
   if err != nil {
      return err
   }
   token_pair, err = amazon.RefreshToken(token_pair.RefreshToken)
   if err != nil {
      return err
   }
   // Updated to receive a *Profile
   profile, err := amazon.GetPrimaryProfile(token_pair.AccessToken)
   if err != nil {
      return fmt.Errorf("Failed to get primary profile: %v", err)
   }
   actor_token, err := amazon.GetActorToken(
      token_pair.RefreshToken, profile.ProfileID,
   )
   if err != nil {
      return fmt.Errorf("Failed to get actor token: %v", err)
   }
   item_details, err := amazon.GetItemDetails(
      actor_token.Token, string(c.TitleId),
   )
   if err != nil {
      return fmt.Errorf("Failed to get item details (playback envelope): %v", err)
   }
   playback, err := amazon.GetVodPlaybackResources(
      actor_token.Token, string(c.TitleId), item_details.PlaybackEnvelope,
   )
   if err != nil {
      return fmt.Errorf("Failed to get VOD playback resources: %v", err)
   }
   url, err := playback.GetUrl()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(url)
   if err != nil {
      return err
   }
   return c.cache.Encode(actor_token, c, item_details, manifest, playback)
}

func (c *client) do_dash_id() error {
   var (
      actor_token  amazon.ActorToken
      item_details amazon.ItemDetails
      manifest     maya.Manifest
      playback amazon.PlaybackResource
   )
   err := c.cache.Decode(&actor_token, &item_details, &manifest, &playback)
   if err != nil {
      return err
   }
   // Fetch the license from Amazon
   license := func(signedRequest []byte) ([]byte, error) {
      return amazon.GetPlayReadyLicense(
         actor_token.Token,
         string(c.TitleId),
         item_details.PlaybackEnvelope,
         playback.SessionHandoffToken,
         signedRequest,
      )
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: license,
   })
}

func (*client) CachePath() string {
   return "rosso/examples/amazon/client"
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_complete_login() error {
   var code_pair amazon.CodePair
   err := c.cache.Decode(&code_pair)
   if err != nil {
      return err
   }
   // Call the updated function which now returns a *TokenPair
   tokenPair, err := amazon.PollRegister(
      code_pair.PublicCode, code_pair.PrivateCode,
   )
   if err != nil {
      return fmt.Errorf("Login incomplete or failed: %v", err)
   }
   // Map the properties of the returned struct into your local test struct
   return c.cache.Encode(tokenPair)
}

func (c *client) do_initiate_login() error {
   codes, err := amazon.CreateCodePair()
   if err != nil {
      return fmt.Errorf("Failed to create code pair: %v", err)
   }
   // Access the properties using dot notation
   err = amazon.InitiateMDSO(codes.PublicCode)
   if err != nil {
      return fmt.Errorf("Failed to initiate MDSO: %v", err)
   }
   log.Print(codes)
   return c.cache.Encode(codes)
}
type client struct {
   TitleId        maya.FlagString
   PlayReady      maya.FlagString
   complete_login maya.FlagBool
   dash_id        maya.FlagString
   initiate_login maya.FlagBool

   cache maya.Cache
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
      {Name: "initiate-login", Value: &c.initiate_login},
      {Name: "complete-login", Value: &c.complete_login},
      {
         Name:  "title-id",
         Value: &c.TitleId,
         Usage: "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624",
      },
      {Name: "dash-id", Value: &c.dash_id},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.PlayReady):
      return c.cache.Encode(c)
   case bool(c.initiate_login):
      return c.do_initiate_login()
   case bool(c.complete_login):
      return c.do_complete_login()
   case flags.IsSet(&c.TitleId):
      return c.do_title_id()
   case c.dash_id != "":
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "amazon")
}
