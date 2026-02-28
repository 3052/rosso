package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "flag"
   "fmt"
   "log"
   "net/http"
   "path"
)

func (c *command) do_playlist() error {
   var title itv.Title
   title.LatestAvailableVersion.PlaylistUrl = c.playlist
   playlist, err := title.Playlist()
   if err != nil {
      return err
   }
   media_file, err := playlist.FullHd()
   if err != nil {
      return err
   }
   err = c.cache.Set("MediaFile", media_file)
   if err != nil {
      return err
   }
   dash, err := media_file.Dash()
   if err != nil {
      return err
   }
   err = c.cache.Set("Dash", dash)
   if err != nil {
      return err
   }
   return maya.ListDash(dash.Body, dash.Url)
}

type command struct {
   cache maya.Cache
   // 1
   address string
   // 2
   playlist string
   // 3
   dash string
   job  maya.WidevineJob
}

func (c *command) run() error {
   c.cache.Init("L3")
   c.job.ClientId = c.cache.Join("client_id.bin")
   c.job.PrivateKey = c.cache.Join("private_key.pem")
   c.cache.Init("itv")
   // 1
   flag.StringVar(&c.address, "a", "", "address")
   // 2
   flag.StringVar(&c.playlist, "p", "", "playlist URL")
   // 3
   flag.StringVar(&c.dash, "d", "", "DASH ID")
   flag.StringVar(&c.job.ClientId, "C", c.job.ClientId, "client ID")
   flag.StringVar(&c.job.PrivateKey, "P", c.job.PrivateKey, "private key")
   flag.Parse()
   if c.address != "" {
      return c.do_address()
   }
   if c.playlist != "" {
      return c.do_playlist()
   }
   if c.dash != "" {
      return c.do_dash()
   }
   return maya.Usage([][]string{
      {"a"},
      {"p"},
      {"d", "C", "P"},
   })
}

func main() {
   // ALL REQUEST ARE GEO BLOCKED
   maya.SetProxy(func(req *http.Request) (string, bool) {
      return "", path.Ext(req.URL.Path) != ".dash"
   })
   err := new(command).run()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *command) do_address() error {
   titles, err := itv.Titles(itv.LegacyId(c.address))
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

func (c *command) do_dash() error {
   var media itv.MediaFile
   err := c.cache.Get("MediaFile", &media)
   if err != nil {
      return err
   }
   c.job.Send = media.Widevine
   var dash itv.Dash
   err = c.cache.Get("Dash", &dash)
   if err != nil {
      return err
   }
   return c.job.DownloadDash(dash.Body, dash.Url, c.dash)
}
