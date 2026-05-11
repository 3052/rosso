package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
)

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      media    plex.Media
      user     plex.User
      widevine device
   )
   err := c.cache.Decode(&manifest, &media, &user, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(&media, &user, body)
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(widevine),
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

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/plex"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   dash := c.flag.String(&c.dash, "d", "DASH ID")
   widevine := c.flag.String(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(device(c.widevine))
   case address.IsSet:
      return c.do_address()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{{
      widevine,
      address,
      dash,
   }})
}

func (c *client) do_address() error {
   user, err := plex.CreateUser()
   if err != nil {
      return err
   }
   path, err := plex.ParsePath(c.address)
   if err != nil {
      return err
   }
   match, err := plex.GetMetadataMatches(path, user)
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
