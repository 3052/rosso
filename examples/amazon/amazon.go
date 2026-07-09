package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
)

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
      Device:       string(c.PlayReady),
      Drm:          maya.DrmPlayReady,
      License:      license,
      MinBandwidth: int(c.min_bandwidth),
   })
}
