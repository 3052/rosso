package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
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
   Widevine maya.FlagString

   address  maya.FlagString
   playlist maya.FlagString
   dash     maya.FlagString

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/itv"); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      if !os.IsNotExist(err) {
         return err
      }
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "address", Value: &c.address},
      {Name: "playlist", Value: &c.playlist},
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
   case c.playlist != "":
      return c.do_playlist()
   case c.dash != "":
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "itv")
}

///

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
   return maya.DownloadDash(c.DashId.Value, &manifest, &maya.Options{
      Device:  widevine.Value,
      Drm:     maya.DrmWidevine,
      License: media_file.FetchKeyService,
   })
}
