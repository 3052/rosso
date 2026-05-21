package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
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

func (c *client) do_address() error {
   path, err := plex.ParsePath(c.Address.Value)
   if err != nil {
      return err
   }
   user, err := plex.CreateUser()
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

type WidevineFolder maya.Flag[string]

type client struct {
   cache          maya.Cache
   WidevineFolder WidevineFolder
   Address        maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/plex"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.WidevineFolder.Set:
      return c.cache.Encode(c.WidevineFolder)
   case c.Address.Set:
      return c.do_address()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "plex", c)
}

func (c *client) do_dash_id() error {
   var (
      manifest maya.Manifest
      media    plex.Media
      user     plex.User
      widevine WidevineFolder
   )
   err := c.cache.Decode(&manifest, &media, &user, &widevine)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(&media, &user, body)
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: license,
   })
}
