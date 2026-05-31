package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
   "os"
)

type client struct {
   Proxy    maya.FlagString
   Widevine maya.FlagString
   address  maya.FlagString
   dash     maya.FlagString
   playlist maya.FlagString
   threads  maya.FlagInt

   cache maya.Cache
}

func (c *client) do() error {
   if err := c.cache.Setup(); err != nil {
      return err
   }
   if err := c.cache.Decode(c); err != nil {
      return c.cache.Encode(c)
   }
   flags := maya.FlagSet{
      {Name: "widevine-folder", Value: &c.Widevine},
      {Name: "proxy", Value: &c.Proxy},
      {Name: "address", Value: &c.address},
      {Name: "playlist", Value: &c.playlist},
      {Name: "dash-id", Value: &c.dash},
      {Name: "threads", Value: &c.threads, Needs: "dash-id"},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if flags.IsSet(&c.Proxy) {
      return c.cache.Encode(c)
   }
   if c.address != "" {
      return c.do_address()
   }
   if err := maya.SetProxy(string(c.Proxy)); err != nil {
      return err
   }
   if c.playlist != "" {
      return c.do_playlist()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return flags.Usage(os.Stderr, "itv")
}

func (c *client) do_dash() error {
   var (
      manifest   maya.Manifest
      media_file itv.MediaFile
   )
   err := c.cache.Decode(&manifest, &media_file)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: media_file.FetchKeyService,
      Threads: int(c.threads),
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (*client) CachePath() string {
   return "rosso/examples/itv/client"
}

func (c *client) do_address() error {
   titles, err := itv.FetchTitles(
      itv.ParseLegacyId(string(c.address)),
   )
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
   playlist, err := itv.FetchWidevine(string(c.playlist))
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
