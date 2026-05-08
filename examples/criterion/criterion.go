package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/criterion"
   "log"
   "path"
)

func (c *client) do_address() error {
   var token criterion.Token
   err := token.Refresh()
   if err != nil {
      return err
   }
   files_href, err := criterion.FetchFilesHref(
      token.AccessToken, path.Base(c.address),
   )
   if err != nil {
      return err
   }
   files, err := criterion.FetchFiles(token.AccessToken, files_href)
   if err != nil {
      return err
   }
   file, err := criterion.GetDash(files)
   if err != nil {
      return err
   }
   manifest, err := file.GetManifest()
   if err != nil {
      return err
   }
   dash, err := maya.ListDash(manifest)
   if err != nil {
      return err
   }
   return c.cache.Encode(dash, file, token)
}

func (c *client) do_dash() error {
   if c.err != nil {
      return c.err
   }
   var (
      dash maya.Dash
      file criterion.File
   )
   err := c.cache.Decode(&dash, &file)
   if err != nil {
      return err
   }
   return dash.Download(&c.job, file.FetchWidevine)
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
   email    string
   err      error
   job      maya.Job
   password string
}

func (c *client) do_email_password() error {
   token, err := criterion.FetchToken(c.email, c.password)
   if err != nil {
      return err
   }
   return c.cache.Encode(token)
}

func (c *client) do() error {
   if err := c.cache.Setup("rosso/criterion"); err != nil {
      return err
   }
   address := maya.StringFlag(&c.address, "a", "address")
   email := maya.StringFlag(&c.email, "e", "email")
   password := maya.StringFlag(&c.password, "p", "password")
   c.err = c.cache.Decode(&c.job)
   dash := maya.StringFlag(&c.job.Dash, "d", "DASH ID")
   widevine := maya.StringFlag(&c.job.Widevine, "w", "Widevine")
   if err := maya.ParseFlags(); err != nil {
      return err
   }
   if widevine.IsSet {
      return c.cache.Encode(c.job)
   }
   if email.IsSet {
      if password.IsSet {
         return c.do_email_password()
      }
   }
   if address.IsSet {
      return c.do_address()
   }
   if dash.IsSet {
      return c.do_dash()
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {email, password},

      {address},
      {dash},
   })
}
