package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
)

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type widevine_device string

func (c *client) do_dash() error {
   var (
      device     widevine_device
      manifest   maya.Manifest
      media_file itv.MediaFile
   )
   err := c.cache.Decode(&device, &manifest, &media_file)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(device),
      Drm:     maya.DrmWidevine,
      License: media_file.FetchKeyService,
   })
}

func (c *client) do_address() error {
   titles, err := itv.FetchTitles(itv.ParseLegacyId(c.address.Value))
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

///

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

type client struct {
   address  string
   cache    maya.Cache
   dash     string
   flag     maya.FlagSet
   playlist string
   widevine string
}

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
      return c.cache.Encode(widevine_device(c.widevine))
   case address.IsSet:
      return c.do_address()
   case playlist.IsSet:
      return c.do_playlist()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([]maya.FlagSet{{
      widevine,
      address,
      playlist,
      dash,
   }})
}
