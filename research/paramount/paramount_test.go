package paramount

import (
   "slices"
   "testing"
)

func TestCookie(t *testing.T) {
   i := slices.IndexFunc(Apps, func(a *app) bool {
      return a.id == "com.cbs.app"
   })
   at, err := GetAt(Apps[i].secret)
   if err != nil {
      t.Fatal(err)
   }
   cookie, err := FetchCbsCom(at, username, password)
   if err != nil {
      t.Fatal(err)
   }
   t.Log(cookie)
}
