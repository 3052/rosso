package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
)

func (c *client) do_dash() error {
   var (
      manifest   maya.Manifest
      media_file itv.MediaFile
      widevine   device
   )
   err := c.cache.Decode(&manifest, &media_file, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: media_file.FetchKeyService,
   })
}

func (c *client) do_address() error {
   titles, err := itv.FetchTitles(itv.ParseLegacyId(c.address))
   if err != nil {
      return err
   }
   for i, title := range titles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&title)
   }
   return nil
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
   playlist string
   widevine string
}

type device string

func (c *client) do() error {
   if err := c.cache.Setup("rosso/itv"); err != nil {
      return err
   }
   address := c.flag.String(&c.address, "a", "address")
   playlist := c.flag.String(&c.playlist, "p", "playlist URL")
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
   case playlist.IsSet:
      return c.do_playlist()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{
      {widevine},
      {address},
      {playlist},
      {dash},
   })
}

func (c *client) do_playlist() error {
   playlist, err := itv.FetchWidevine(c.playlist)
   if err != nil {
      return err
   }
   media_file, err := playlist.Get1080()
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(&media_file.Href.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(manifest, media_file)
}
