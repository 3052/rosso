package paramount

import (
   "fmt"
   "strings"
)

func CbsAppIds() string {
   var data strings.Builder
   for i, app := range CbsApps {
      if i >= 1 {
         data.WriteByte(' ')
      }
      data.WriteString(app.Id)
   }
   return data.String()
}

func GetCbsApp(id string) (*CbsApp, error) {
   for _, app := range CbsApps {
      if app.Id == id {
         return &app, nil
      }
   }
   return nil, fmt.Errorf("CBS app not found %q", id)
}

var CbsApps = []CbsApp{
   {
      Id:      "com.cbs.app",
      Host:    "www.paramountplus.com",
      Secret:  "7081400bd4143bf3",
      Version: "Paramount+ 16.8.0",
   },
   {
      Id:      "com.cbs.ca",
      Host:    "www.paramountplus.com",
      Secret:  "1c5d27627d71b420",
      Version: "Paramount+ 16.8.0",
   },
   {
      Id:      "com.cbs.tve",
      Host:    "www.cbs.com",
      Secret:  "cef32931dc01412e",
      Version: "CBS 15.6.0",
   },
}

type CbsApp struct {
   Id      string
   Host    string
   Secret  string
   Version string
}
