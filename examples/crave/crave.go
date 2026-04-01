package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/crave"
   "fmt"
   "log"
)

func (c *client) do_username_password() error {
   var err error
   c.Account, err = crave.Login(c.username, c.password)
   if err != nil {
      return err
   }
   profiles, err := c.Account.FetchProfiles()
   if err != nil {
      return err
   }
   for i, profile := range profiles {
      if i >= 1 {
         fmt.Println()
      }
      fmt.Println(profile)
   }
   return nil
}

type client struct {
   Account *crave.Account
   //--------------------
   Job maya.Job
   //--------------------
   username string
   password string
}

func (c *client) do() error {
   err := cache.Setup("rosso/crave.xml")
   if err != nil {
      return err
   }
   // with_cache := cache.Read(c)
   widevine := maya.StringFlag(&c.Job.Widevine, "w", "Widevine")
   //-----------------------------------------------------------
   username := maya.StringFlag(&c.username, "u", "username")
   password := maya.StringFlag(&c.password, "p", "password")
   err = maya.ParseFlags()
   if err != nil {
      return err
   }
   if widevine.IsSet {
      return cache.Write(c)
   }
   if username.IsSet {
      if password.IsSet {
         return c.do_username_password()
      }
   }
   return maya.PrintFlags([][]*maya.Flag{
      {widevine},
      {username, password},
   })
}

var cache maya.Cache

func main() {
   log.SetFlags(log.Ltime)
   maya.SetProxy("", "")
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}
