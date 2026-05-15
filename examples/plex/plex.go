package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "fmt"
   "log"
)

type client struct {
   cache    maya.Cache
   flag     maya.FlagSet
   address  maya.Flag
   dash     maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/plex"); err != nil {
      return err
   }
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.widevine.Set:
      return c.cache.Encode(widevine_device(c.widevine.Value))
   case c.address.Set:
      return c.do_address()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}

func (c *client) do_dash() error {
   var (
      device   widevine_device
      manifest maya.Manifest
      media    plex.Media
      user     plex.User
   )
   err := c.cache.Decode(&device, &manifest, &media, &user)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(&media, &user, body)
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: license,
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type widevine_device string

func (c *client) do_address() error {
   user, err := plex.CreateUser()
   if err != nil {
      return err
   }
   address, err := c.address.ParseUrl()
   if err != nil {
      return err
   }
   match, err := plex.GetMetadataMatches(plex.ParsePath(address), user)
   if err != nil {
      return err
   }
   vod_metadata, err := plex.GetVodMetadata(&match.Metadata[0], user)
   if err != nil {
      return err
   }
   media, err := vod_metadata.GetDash()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(media.GetManifest(user))
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, media, user)
}
