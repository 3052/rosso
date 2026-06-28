package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
   "fmt"
   "log"
   "os"
)

func (c *client) do_title_id() error {
   var token_pair amazon.TokenPair
   err := c.cache.Decode(&token_pair)
   if err != nil {
      return err
   }
   if err = token_pair.Refresh(); err != nil {
      return err
   }
   profile, err := amazon.GetPrimaryProfile(
      token_pair.AccessToken, string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("failed to get primary profile: %v", err)
   }
   actor_token, err := amazon.GetActorToken(
      token_pair.RefreshToken, profile.ProfileID, string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("failed to get actor token: %v", err)
   }
   item_details, err := amazon.GetItemDetails(
      actor_token.Token, string(c.title_id), string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("failed to get item details (playback envelope): %v", err)
   }
   playback := amazon.VodPlaybackParams{
      ActorAccessToken:   actor_token.Token,
      BitrateAdaptation:  string(c.bitrate_adaptation),
      DeviceTypeID: string(c.DeviceTypeId),
      DRMType:            "PlayReady",
      DynamicRangeFormat: string(c.dynamic_range),
      MaxVideoResolution: "2160p",
      PlaybackEnvelope:   item_details.PlaybackEnvelope,
      TitleId:            string(c.title_id),
      VideoCodec:         string(c.video_codec),
   }
   resources, err := playback.Fetch()
   if err != nil {
      return fmt.Errorf("failed to get VOD playback resources: %v", err)
   }
   clean, err := resources.Clean()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(clean)
   if err != nil {
      return err
   }
   return c.cache.Encode(actor_token, item_details, manifest)
}
func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   DeviceTypeId       maya.FlagString
   PlayReady          maya.FlagString
   bitrate_adaptation maya.FlagString
   complete_login     maya.FlagBool
   dash_id            maya.FlagString
   dynamic_range      maya.FlagString
   initiate_login     maya.FlagBool
   title_id           maya.FlagString
   video_codec        maya.FlagString

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/amazon/client"
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }

   c.bitrate_adaptation = "CVBR"
   c.dynamic_range = "None"
   c.video_codec = "H265"

   flags := maya.FlagSet{
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "device-type-id", Value: &c.DeviceTypeId},
      {Name: "initiate-login", Value: &c.initiate_login},
      {Name: "complete-login", Value: &c.complete_login},
      {
         Name:  "title-id",
         Value: &c.title_id,
         Usage: "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624",
      },
      {
         Name:  "video-codec",
         Value: &c.video_codec,
         Usage: "H264 H265",
         Needs: "title-id",
      },
      {
         Name:  "bitrate-adaptation",
         Value: &c.bitrate_adaptation,
         Usage: "CVBR CBR",
         Needs: "title-id",
      },
      {
         Name:  "dynamic-range",
         Value: &c.dynamic_range,
         Usage: "None HDR10 DolbyVision",
         Needs: "title-id",
      },
      {Name: "dash-id", Value: &c.dash_id},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.PlayReady):
      return c.cache.Encode(c)
   case flags.IsSet(&c.DeviceTypeId):
      return c.cache.Encode(c)
   case bool(c.initiate_login):
      return c.do_initiate_login()
   case bool(c.complete_login):
      return c.do_complete_login()
   case c.title_id != "":
      return c.do_title_id()
   case c.dash_id != "":
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "amazon")
}

func (c *client) do_complete_login() error {
   var code_pair amazon.CodePair
   err := c.cache.Decode(&code_pair)
   if err != nil {
      return err
   }
   tokenPair, err := amazon.PollRegister(
      code_pair.PublicCode, code_pair.PrivateCode, string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("login incomplete or failed: %v", err)
   }
   return c.cache.Encode(tokenPair)
}

func (c *client) do_dash_id() error {
   var (
      actor_token  amazon.ActorToken
      item_details amazon.ItemDetails
      manifest     maya.Manifest
   )
   err := c.cache.Decode(&actor_token, &item_details, &manifest)
   if err != nil {
      return err
   }
   // Fetch the license from Amazon
   license := func(signedRequest []byte) ([]byte, error) {
      return amazon.GetPlayReadyLicense(
         actor_token.Token,
         item_details.PlaybackEnvelope,
         signedRequest,
         string(c.DeviceTypeId),
      )
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.PlayReady),
      Drm:     maya.DrmPlayReady,
      License: license,
   })
}

func (c *client) do_initiate_login() error {
   codes, err := amazon.CreateCodePair(string(c.DeviceTypeId))
   if err != nil {
      return fmt.Errorf("failed to create code pair: %v", err)
   }
   log.Print(codes)
   return c.cache.Encode(codes)
}

