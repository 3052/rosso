package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/molotov"
   "fmt"
   "log"
   "net/url"
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
   asset_id maya.FlagString
   dash_id  maya.FlagString
   username maya.FlagString
   password maya.FlagString
   search   maya.FlagString
   refresh  maya.FlagBool
   threads  maya.FlagInt

   cache maya.Cache
}

func (*client) CachePath() string {
   return "rosso/examples/molotov/client"
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
      {Name: "username", Value: &c.username, Needs: "password"},
      {Name: "password", Value: &c.password, Needs: "username"},
      {Name: "refresh", Value: &c.refresh},
      {Name: "search", Value: &c.search},
      {Name: "asset-id", Value: &c.asset_id},
      {Name: "dash-id", Value: &c.dash_id},
      {Name: "threads", Value: &c.threads, Needs: "dash-id"},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   if flags.IsSet(&c.Widevine) {
      return c.cache.Encode(c)
   }
   if c.username != "" {
      if c.password != "" {
         return c.do_username_password()
      }
   }
   if c.refresh {
      return c.do_refresh()
   }
   if c.search != "" {
      return c.do_search()
   }
   if c.asset_id != "" {
      return c.do_asset_id()
   }
   if c.dash_id != "" {
      return c.do_dash_id()
   }
   return flags.Usage(os.Stderr, "molotov")
}

func (c *client) do_asset_id() error {
   var signin molotov.SigninResponse
   err := c.cache.Decode(&signin)
   if err != nil {
      return err
   }
   asset, err := molotov.GetAsset(string(c.asset_id), &signin)
   if err != nil {
      return err
   }
   address, err := url.Parse(asset.Stream.URL)
   if err != nil {
      return err
   }
   manifest, err := maya.ListDash(address)
   if err != nil {
      return err
   }
   return c.cache.Encode(asset, manifest)
}

func (c *client) do_dash_id() error {
   var (
      asset    molotov.AssetResponse
      manifest maya.Manifest
   )
   err := c.cache.Decode(&asset, &manifest)
   if err != nil {
      return err
   }
   return maya.DownloadDash(string(c.dash_id), &manifest, &maya.Options{
      Device:  string(c.Widevine),
      Drm:     maya.DrmWidevine,
      License: asset.GetLicense,
      Threads: int(c.threads),
   })
}

func (c *client) do_refresh() error {
   var signin molotov.SigninResponse
   err := c.cache.Decode(&signin)
   if err != nil {
      return err
   }
   err = signin.Refresh()
   if err != nil {
      return err
   }
   return c.cache.Encode(&signin)
}

func (c *client) do_search() error {
   var signin molotov.SigninResponse
   err := c.cache.Decode(&signin)
   if err != nil {
      return err
   }
   results, err := molotov.Search(string(c.search), &signin)
   if err != nil {
      return err
   }
   for i, result := range results {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(&result)
   }
   return nil
}

func (c *client) do_username_password() error {
   signin, err := molotov.Signin(string(c.username), string(c.password))
   if err != nil {
      return err
   }
   return c.cache.Encode(signin)
}
