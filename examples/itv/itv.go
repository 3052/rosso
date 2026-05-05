package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "fmt"
   "log"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/itv.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
   //----------------------------------------------------------
   playlist := maya.StringFlag(&c.playlist, "p", "playlist URL")
   //----------------------------------------------------------
   dash := maya.StringFlag(&c.Job.Dash, "d", "DASH ID")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   switch {
   case widevine.IsSet:
      return cache.Write(c)
   case address.IsSet:
      return c.do_address()
   case playlist.IsSet:
      return c.do_playlist()
   case dash.IsSet:
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{{
      widevine,
      address,
      playlist,
      dash,
   }})
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
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
   c.Dash, err = maya.ListDash(c.MediaFile.GetManifest)
   if err != nil {
      return err
   }
   return cache.Write(c)
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, c.MediaFile.FetchKeyService)
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   // cache
   Dash      *maya.Dash
   Job       maya.Job
   MediaFile *itv.MediaFile
   // flags
   address  string
   playlist string
   // state
   cache_err error
}
