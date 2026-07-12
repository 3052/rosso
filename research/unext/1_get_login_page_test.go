// 1_get_login_page_test.go
package unext

import (
   "strings"
   "testing"
)

func TestGetLoginPage(t *testing.T) {
   client, err := newProxyClient()
   if err != nil {
      t.Fatalf("failed to create proxy client: %v", err)
   }

   token, err := GetLoginPage(client)
   if err != nil {
      t.Fatalf("GetLoginPage failed: %v", err)
   }

   if token == "" {
      t.Fatal("expected non-empty CSRF token, got empty string")
   }

   if strings.Contains(token, " ") {
      t.Errorf("CSRF token should not contain spaces: '%s'", token)
   }

   t.Logf("received CSRF token: %s", token)
}
