package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/itv"
   "flag"
   "fmt"
   "log"
)

func (c *client) do() error {
   err := cache.Setup("rosso/itv.xml")
   if err != nil {
      return err
   }
   with_cache := cache.Read(c)
   widevine := maya.StringVar(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringVar(&c.address, "a", "address")
   //----------------------------------------------------------
   playlist := maya.StringVar(&c.playlist, "p", "playlist URL")
   //----------------------------------------------------------
   dash_id := maya.StringVar(&c.dash_id, "d", "DASH ID")
   set := maya.Parse()
   switch {
   case set[widevine]:
      return cache.Write(c)
   case set[address]:
      return c.do_address()
   case set[playlist]:
      return c.do_playlist()
   case set[dash_id]:
      return with_cache(c.do_dash_id)
   }
   return maya.Usage([][]*flag.Flag{
      {widevine},
      {address},
      {playlist},
      {dash_id},
   })
}

func (c *client) do_dash_id() error {
   return c.Job.DownloadDash(
      c.Dash.Body, c.Dash.Url, c.dash_id, c.MediaFile.Widevine,
   )
}

func main() {
   log.SetFlags(log.Ltime)
   // ALL REQUEST ARE GEO BLOCKED
   maya.SetProxy("", "*.dash")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

var cache maya.Cache

func (c *client) do_address() error {
   titles, err := itv.Titles(itv.ParseLegacyId(c.address))
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
   playlist, err := itv.FetchPlaylist(c.playlist)
   if err != nil {
      return err
   }
   c.MediaFile, err = playlist.FullHd()
   if err != nil {
      return err
   }
   c.Dash, err = c.MediaFile.Dash()
   if err != nil {
      return err
   }
   err = cache.Write(c)
   if err != nil {
      return err
   }
   return maya.ListDash(c.Dash.Body, c.Dash.Url)
}

type client struct {
   Dash      *itv.Dash
   MediaFile *itv.MediaFile
   //----------------------
   Job maya.Job
   //----------------------
   address string
   //----------------------
   playlist string
   //----------------------
   dash_id string
}
