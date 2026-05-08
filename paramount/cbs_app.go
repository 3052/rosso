package paramount

import (
   "fmt"
   "maps"
   "slices"
)

func GetCbsApp(id string) (*CbsApp, error) {
   app, found := CbsApps[id]
   if !found {
      return nil, fmt.Errorf("CBS app not found %q", id)
   }
   return &app, nil
}

var CbsApps = map[string]CbsApp{
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

type CbsApp struct {
   Host    string
   Version string
   Secret  string
}

func CbsAppIds() []string {
   return slices.Sorted(maps.Keys(CbsApps))
}
