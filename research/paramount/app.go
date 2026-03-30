package paramount

import "fmt"

type App struct {
   Host    string
   Version string
   Secret  string
}

var apps = map[string]App{
   "com.cbs.app": {
      Host:    "www.paramountplus.com",
      Version: "Paramount+ 16.8.0",
      Secret:  "7081400bd4143bf3",
   },
   "com.cbs.ca": {
      Host:    "www.paramountplus.com",
      Version: "Paramount+ 16.8.0",
      Secret:  "1c5d27627d71b420",
   },
   "com.cbs.tve": {
      Host:    "www.cbs.com",
      Version: "CBS 15.6.0",
      Secret:  "cef32931dc01412e",
   },
}

func GetApp(id string) (*App, error) {
   app, exists := apps[id]
   if !exists {
      return nil, fmt.Errorf("app not found: %s", id)
   }
   return &app, nil
}
