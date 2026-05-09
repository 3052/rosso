package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
)

func (c *client) do_playlist() error {
   playlist, err := itv.FetchWidevine(c.playlist)
   if err != nil {
      return err
   }
   media_file, err := playlist.Get1080()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(&media_file.Href.Url)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, media_file)
}

func (c *client) do_dash() error {
   var (
      dash       maya.Dash
      media_file itv.MediaFile
   )
   err := c.cache.Decode(&c.job, &dash, &media_file)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, media_file.FetchKeyService)
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
   job      maya.Job
   playlist string
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/itv"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   playlist := maya.StringFlag(&c.playlist, "p", "playlist URL")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   dash := maya.StringFlag(&c.dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return c.cache.Encode(c.job)
   case address.IsSet:
      return c.do_address()
   case playlist.IsSet:
      return c.do_playlist()
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      address,
      playlist,
      dash,
   }})
}
