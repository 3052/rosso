package canal

import (
   "encoding/json"
   "net/http"
   "os/exec"
   "testing"
)

func Test(t *testing.T) {
   app := &App{
      Client:       &http.Client{},
      DeviceSerial: "w76d15b90-3215-11f1-87ca-01f0af932fb7", // Replace with a dynamic UUID if needed
   }
   app.TVApiBaseURL = "https://tvapi-hlm2.solocoo.tv"
   if err := app.LoginInit(); err != nil {
      t.Fatalf("LoginInit failed: %v", err)
   }
   data, err := exec.Command("credential", "-j=canalplus.cz").Output()
   if err != nil {
      t.Fatal(err)
   }
   var credential []struct {
      Username string
      Password string
   }
   err = json.Unmarshal(data, &credential)
   if err != nil {
      t.Fatal(err)
   }
   err = app.LoginSubmit(credential[0].Username, credential[0].Password)
   if err != nil {
      t.Fatalf("LoginSubmit failed: %v", err)
   }
}
