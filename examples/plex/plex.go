package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/plex"
   "log"
)

func (c *client) do_dash() error {
   var (
      dash  maya.Dash
      media plex.Media
      user  plex.User
   )
   err := c.cache.Decode(&c.job, &dash, &media, &user)
   if err != nil {
      return err
   }
   return dash.Download(c.dash, &c.job, func(body []byte) ([]byte, error) {
      return plex.AcquireWidevineLicense(&media, &user, body)
   })
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   address string
   cache   maya.Cache
   dash    string
   job     maya.Job
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/plex"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
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
   case dash.IsSet:
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},

      {address},
      {dash},
   })
}

func (c *client) do_address() error {
   user, err := plex.CreateUser()
   if err != nil {
      return err
   }
   path, err := plex.ParsePath(c.address)
   if err != nil {
      return err
   }
   match, err := plex.GetMetadataMatches(path, user)
   if err != nil {
      return err
   }
   vod_metadata, err := plex.GetVodMetadata(&match.Metadata[0], user)
   if err != nil {
      return err
   }
   media, err := vod_metadata.GetDash()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(media.GetManifest(user))
   if err != nil {
      return err
   }
   return c.cache.Encode(user, dash, media)
}
