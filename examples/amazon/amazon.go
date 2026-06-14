package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
   "fmt"
   "log"
   "os"
)

// 6 license_test.go

func (c *client) do_title_id() error {
   var actor_token amazon.ActorToken
   err := c.cache.Decode(&actor_token)
   if err != nil {
      return err
   }
   // Calling the updated function which returns an *ItemDetails
   itemDetails, err := amazon.GetItemDetails(
      actor_token.Token, string(c.title_id),
   )
   if err != nil {
      return fmt.Errorf("Failed to get item details (playback envelope): %v", err)
   }
   mpdUrl, err := amazon.GetVodPlaybackResources(
      actor_token.Token, string(c.title_id), itemDetails.PlaybackEnvelope,
   )
   if err != nil {
      return fmt.Errorf("Failed to get VOD playback resources: %v", err)
   }
   log.Println("mpdUrl", mpdUrl)
   // Map the properties of the returned struct into your local test struct
   return c.cache.Encode(itemDetails)
}

type client struct {
   Widevine       maya.FlagString
   actor_token    maya.FlagBool
   complete_login maya.FlagBool
   initiate_login maya.FlagBool
   title_id       maya.FlagString

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
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "initiate-login", Value: &c.initiate_login},
      {Name: "complete-login", Value: &c.complete_login},
      {Name: "actor-token", Value: &c.actor_token},
      {
         Name:  "title-id",
         Value: &c.title_id,
         Usage: "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624",
      },
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.Widevine):
      return c.cache.Encode(c)
   case bool(c.initiate_login):
      return c.do_initiate_login()
   case bool(c.complete_login):
      return c.do_complete_login()
   case bool(c.actor_token):
      return c.do_actor_token()
   case c.title_id != "":
      return c.do_title_id()
   }
   return flags.Usage(os.Stderr, "amazon")
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
   // Call the updated function which now returns a *CodePair
   codes, err := amazon.CreateCodePair()
   if err != nil {
      return fmt.Errorf("Failed to create code pair: %v", err)
   }
   // Access the properties using dot notation
   err = amazon.InitiateMDSO(codes.PublicCode)
   if err != nil {
      return fmt.Errorf("Failed to initiate MDSO: %v", err)
   }
   return c.cache.Encode(codes)
}

func (c *client) do_actor_token() error {
   var token_pair amazon.TokenPair
   err := c.cache.Decode(&token_pair)
   if err != nil {
      return err
   }
   // Updated to receive a *Profile
   profile, err := amazon.GetPrimaryProfile(token_pair.AccessToken)
   if err != nil {
      return fmt.Errorf("Failed to get primary profile: %v", err)
   }
   // Pass the extracted string to GetActorToken and receive an *ActorToken
   actorToken, err := amazon.GetActorToken(
      token_pair.RefreshToken, profile.ProfileID,
   )
   if err != nil {
      return fmt.Errorf("Failed to get actor token: %v", err)
   }
   // Map the properties of the returned structs into your local test struct
   return c.cache.Encode(actorToken, profile)
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
