package peacock

import "testing"

func TestGetProviderVariant(t *testing.T) {
   c := NewClient("")
   resp, err := c.GetProviderVariant("1cba422b-3533-33a4-84af-d57cb97bbfa1")
   if err != nil {
      t.Fatal(err)
   }
   t.Logf("body: %s", resp.Body)
}
