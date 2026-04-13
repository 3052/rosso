package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/itv.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   proxy := maya.StringFlag(&c.Proxy, "x", "proxy")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   playlist := maya.StringFlag(&c.playlist, "p", "playlist URL")
   //----------------------------------------------------------
   dash_id := maya.StringFlag(&c.dash_id, "d", "DASH ID")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   err = maya.SetProxy(c.Proxy, "*.dash")
   if err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case proxy.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case playlist.IsSet:
      return c.do_playlist()
   case dash_id.IsSet:
      return with_cache(c.do_dash_id)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      proxy,
      address,
      playlist,
      dash_id,
   }})
}

var cache maya.Cache

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
func (c *client) do_playlist() error {
   playlist, err := itv.FetchWidevine(c.playlist)
   if err != nil {
      return err
   }
   c.MediaFile, err = playlist.Get1080()
   if err != nil {
      return err
   }
   c.Dash, err = c.MediaFile.FetchDash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.MediaFile.FetchKeyService,
   )
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
   //----------------------
   Job maya.Job
   //----------------------
   Proxy string
   //----------------------
   address string
   //----------------------
   playlist string
   //----------------------
   dash_id string
}
