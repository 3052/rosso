package main

import (
   "41.neocities.org/maya"
   "41.neocities.org/rosso/amazon"
   "fmt"
   "log"
   "os"
)

func (c *client) do_actor_token() error {
   return nil
}

func (*client) CachePath() string {
   return "rosso/examples/amazon/client"
}

func main() {
   log.SetFlags(log.Ltime)
   err := new(client).do()
   if err != nil {
      log.Fatal(err)
   }
}

type client struct {
   Widevine       maya.FlagString
   actor_token    maya.FlagBool
   complete_login maya.FlagBool
   initiate_login maya.FlagBool

   cache maya.Cache
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
      {Name: "initiate-login", Value: &c.initiate_login},
      {Name: "complete-login", Value: &c.complete_login},
      {Name: "actor-token", Value: &c.actor_token},
   }
   if err := flags.Parse(os.Args[1:]); err != nil {
      return err
   }
   switch {
   case flags.IsSet(&c.Widevine):
      return c.cache.Encode(c)
   case bool(c.initiate_login):
      return c.do_initiate_login()
   case bool(c.complete_login):
      return c.do_complete_login()
   case bool(c.actor_token):
      return c.do_actor_token()
   }
   return flags.Usage(os.Stderr, "amazon")
}

func (c *client) do_complete_login() error {
   var code_pair amazon.CodePair
   err := c.cache.Decode(&code_pair)
   if err != nil {
      return err
   }
   // Call the updated function which now returns a *TokenPair
   tokenPair, err := amazon.PollRegister(
      code_pair.PublicCode, code_pair.PrivateCode,
   )
   if err != nil {
      return fmt.Errorf("Login incomplete or failed: %v", err)
   }
   // Map the properties of the returned struct into your local test struct
   return c.cache.Encode(tokenPair)
}

func (c *client) do_initiate_login() error {
   // Call the updated function which now returns a *CodePair
   codes, err := amazon.CreateCodePair()
   if err != nil {
      return fmt.Errorf("Failed to create code pair: %v", err)
   }
   // Access the properties using dot notation
   err = amazon.InitiateMDSO(codes.PublicCode)
   if err != nil {
      return fmt.Errorf("Failed to initiate MDSO: %v", err)
   }
   return c.cache.Encode(codes)
}
