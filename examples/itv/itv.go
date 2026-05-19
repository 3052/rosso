package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
   "os"
)

type client struct {
   cache          maya.Cache
   WidevineFolder maya.Flag[string]
   Address        maya.Flag[string]
   Playlist       maya.Flag[string]
   DashId         maya.Flag[string]
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/itv"); err != nil {
      return err
   }
   if err := maya.ParseFlags(os.Args[1:], c); err != nil {
      return err
   }
   switch {
   case c.WidevineFolder.Set:
      return c.cache.Encode(WidevineFolder(c.WidevineFolder.Value))
   case c.Address.Set:
      return c.do_address()
   case c.Playlist.Set:
      return c.do_playlist()
   case c.DashId.Set:
      return c.do_dash_id()
   }
   return maya.FormatFlags(os.Stderr, "itv", c)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type WidevineFolder string

func (c *client) do_dash_id() error {
   var (
      manifest   maya.Manifest
      media_file itv.MediaFile
      widevine   WidevineFolder
   )
   err := c.cache.Decode(&manifest, &media_file, &widevine)
   if err != nil {
      return err
   }
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  string(widevine),
      Drm:     maya.DrmWidevine,
      License: media_file.FetchKeyService,
   })
}

func (c *client) do_address() error {
   titles, err := itv.FetchTitles(itv.ParseLegacyId(c.Address.Value))
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

func (c *client) do_playlist() error {
   playlist, err := itv.FetchWidevine(c.Playlist.Value)
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
