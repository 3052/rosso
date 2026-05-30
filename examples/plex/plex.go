package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
   "os"
)

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "address", Value: &c.address},
      {Name: "dash-id", Value: &c.dash},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.Widevine):
      return c.cache.Encode(c)
   case c.address != "":
      return c.do_address()
   case c.dash != "":
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "plex")
}

func (c *client) do_address() error {
   path, err := plex.ParsePath(string(c.address))
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

func (c *client) do_dash() error {
   var (
      manifest maya.Manifest
      media    plex.Media
      user     plex.User
   )
   err := c.cache.Decode(&manifest, &media, &user)
   if err != nil {
      return err
   }
   license := func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(&media, &user, body)
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
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
   Widevine maya.FlagString

   address maya.FlagString
   dash    maya.FlagString

   cache maya.Cache
}
