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

type WidevineFolder string

func (c *client) do_dash() error {
   var (
      manifest   maya.Manifest
      media_file itv.MediaFile
      widevine   WidevineFolder
   )
   err := c.cache.Decode(&manifest, &media_file, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.dash.Value, &manifest, &maya.Options{
      Device:  string(widevine),
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
   address, err := c.playlist.ParseUrl()
   if err != nil {
      return err
   }
   playlist, err := itv.FetchWidevine(address)
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
   cache    maya.Cache
   flag     maya.FlagSet
   address  maya.Flag
   dash     maya.Flag
   playlist maya.Flag
   widevine maya.Flag
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/itv"); err != nil {
      return err
   }
   c.flag.AddValue(&c.address, "a", "address")
   c.flag.AddValue(&c.playlist, "p", "playlist URL")
   c.flag.AddValue(&c.dash, "d", "DASH ID")
   c.flag.AddValue(&c.widevine, "w", "Widevine")
   if err := c.flag.Parse(); err != nil {
      return err
   }
   switch {
   case c.widevine.Set:
      return c.cache.Encode(WidevineFolder(c.widevine.Value))
   case c.address.Set:
      return c.do_address()
   case c.playlist.Set:
      return c.do_playlist()
   case c.dash.Set:
      return c.do_dash()
   }
   fmt.Println(c.flag)
   return nil
}
