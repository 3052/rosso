package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
   "net/url"
)

func (c *client) do() error {
   if err := cache.Setup("rosso/plex.xml"); err != nil {
      return err
   }
   c.cache_err = cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //----------------------------------------------------------
   address := maya.StringFlag(&c.address, "a", "address")
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
   case dash.IsSet:
      return c.run(c.do_dash)
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {address},
      {dash},
   })
}

func (c *client) run(action func() error) error {
   if c.cache_err != nil {
      return c.cache_err
   }
   return action()
}

func (c *client) do_dash() error {
   return c.Dash.Download(&c.Job, func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(c.VodMedia, c.AnonymousUser, body)
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

func (c *client) do_address() error {
   var err error
   c.AnonymousUser, err = plex.CreateAnonymousUser()
   if err != nil {
      return err
   }
   path, err := plex.ParsePath(c.address)
   if err != nil {
      return err
   }
   match, err := plex.GetMetadataMatches(path, c.AnonymousUser)
   if err != nil {
      return err
   }
   vod_metadata, err := plex.GetVodMetadata(&match.Metadata[0], c.AnonymousUser)
   if err != nil {
      return err
   }
   c.VodMedia, err = vod_metadata.GetDashMedia()
   if err != nil {
      return err
   }
   c.Dash, err = maya.ListDash(func() (*url.URL, error) {
      return c.VodMedia.GetMpdUrl(c.AnonymousUser)
   })
   if err != nil {
      return err
   }
   return cache.Write(c)
}

var cache maya.Cache

type client struct {
   // cache
   AnonymousUser *plex.AnonymousUser
   Dash          *maya.Dash
   Job           maya.Job
   VodMedia      *plex.VodMedia
   // flags
   address string
   // state
   cache_err error
}
