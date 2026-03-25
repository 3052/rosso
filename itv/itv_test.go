package itv

import "testing"

var watch_tests = []struct {
   category string
   watch    []string
}{
   {
      category: "DRAMA_AND_SOAPS",
      watch: []string{
         "https://itv.com/watch/joan/10a3918",
         "https://itv.com/watch/joan/10a3918/10a3918a0001",
      },
   },
   {
      category: "FILM",
      watch:    []string{"https://itv.com/watch/love-actually/27304"},
   },
}

func TestWatch(t *testing.T) {
   t.Log(watch_tests)
}
