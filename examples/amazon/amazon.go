package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
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
   DeviceTypeId       maya.FlagString
   PlayReady          maya.FlagString
   TitleId            maya.FlagString
   bitrate_adaptation maya.FlagString
   complete_login     maya.FlagBool
   dash_id            maya.FlagString
   dynamic_range      maya.FlagString
   initiate_login     maya.FlagBool
   playback           maya.FlagBool
   video_codec        maya.FlagString
   min_bitrate        maya.FlagInt

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
   c.video_codec = "H264"

   flags := maya.FlagSet{
      {Name: "device-type-id", Value: &c.DeviceTypeId},
      {Name: "playReady-folder", Value: &c.PlayReady},
      {Name: "initiate-login", Value: &c.initiate_login},
      {Name: "complete-login", Value: &c.complete_login},
      {Name: "title-id", Value: &c.TitleId},
      {Name: "playback", Value: &c.playback},
      {
         Name:  "bitrate-adaptation",
         Value: &c.bitrate_adaptation,
         Usage: "CVBR CBR",
         Needs: "playback",
      },
      {
         Name:  "dynamic-range",
         Value: &c.dynamic_range,
         Usage: "None HDR10 DolbyVision",
         Needs: "playback",
      },
      {
         Name:  "video-codec",
         Value: &c.video_codec,
         Usage: "H264 H265",
         Needs: "playback",
      },
      {Name: "dash-id", Value: &c.dash_id},
      {Name: "min-bitrate", Value: &c.min_bitrate, Needs: "dash-id"},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.DeviceTypeId):
      return c.cache.Encode(c)
   case flags.IsSet(&c.PlayReady):
      return c.cache.Encode(c)
   case bool(c.initiate_login):
      return c.do_initiate_login()
   case bool(c.complete_login):
      return c.do_complete_login()
   case flags.IsSet(&c.TitleId):
      return c.do_title_id()
   case bool(c.playback):
      return c.do_playback()
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
      &code_pair, string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("login incomplete or failed: %v", err)
   }
   return c.cache.Encode(tokenPair)
}
func (c *client) do_dash_id() error {
   var (
      actor_token amazon.ActorToken
      manifest    maya.Manifest
      metadata    amazon.PlaybackExperienceMetadata
   )
   err := c.cache.Decode(&actor_token, &manifest, &metadata)
   if err != nil {
      return err
   }
   // Fetch the license from Amazon
   license := func(signedRequest []byte) ([]byte, error) {
      return amazon.GetPlayReadyLicense(
         &actor_token,
         &metadata,
         signedRequest,
         string(c.DeviceTypeId),
      )
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:     string(c.PlayReady),
      Drm:        maya.DrmPlayReady,
      License:    license,
      MinBitrate: int(c.min_bitrate),
   })
}

func (c *client) do_initiate_login() error {
   codes, err := amazon.CreateCodePair(string(c.DeviceTypeId))
   if err != nil {
      return fmt.Errorf("failed to create code pair: %v", err)
   }
   fmt.Println(codes)
   return c.cache.Encode(codes)
}

///

func (c *client) do_playback() error {
   var token_pair amazon.TokenPair
   err := c.cache.Decode(&token_pair)
   if err != nil {
      return err
   }
   if err = token_pair.Refresh(); err != nil {
      return err
   }
   profile, err := amazon.GetPrimaryProfile(&token_pair, string(c.DeviceTypeId))
   if err != nil {
      return fmt.Errorf("failed to get primary profile: %v", err)
   }
   actor_token, err := amazon.GetActorToken(
      &token_pair, profile, string(c.DeviceTypeId),
   )
   if err != nil {
      return fmt.Errorf("failed to get actor token: %v", err)
   }
   resource, err := amazon.GetItemDetails(
      actor_token, string(c.TitleId), string(c.DeviceTypeId),
   )
   if err != nil {
      return err
   }
   metadata, err := resource.GetPlaybackExperienceMetadata()
   if err != nil {
      return err
   }
   playback := amazon.VodPlaybackParams{
      DRMType:                    "PlayReady",
      MaxVideoResolution:         "2160p",
      BitrateAdaptation:          string(c.bitrate_adaptation),
      DynamicRangeFormat:         string(c.dynamic_range),
      VideoCodec:                 string(c.video_codec),
      DeviceTypeID:               string(c.DeviceTypeId),
      TitleId:                    string(c.TitleId),
      ActorToken:                 actor_token,
      PlaybackExperienceMetadata: metadata,
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
   return c.cache.Encode(actor_token, manifest, metadata)
}

func (c *client) do_title_id() error {
   resource, err := amazon.GetItemDetails(
      nil, string(c.TitleId), string(c.DeviceTypeId),
   )
   if err != nil {
      return err
   }
   fmt.Println(resource)
   return c.cache.Encode(c)
}
