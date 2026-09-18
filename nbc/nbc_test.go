package nbc

import "testing"

var video_tests = []struct {
   url     string
   program string
   lock    bool
}{
   {
      lock:    true,
      program: "movie",
      url:     "https://nbc.com/las-sobrinas-del-diablo/video/las-sobrinas-del-diablo/8000011982",
   },
   {
      lock:    false,
      program: "episode",
      url:     "https://nbc.com/saturday-night-live/video/november-8-nikki-glaser/9000454169",
   },
   {
      lock:    true,
      program: "episode",
      url:     "https://nbc.com/saturday-night-live/video/november-1-miles-teller/9000454168",
   },
}

func TestVideo(t *testing.T) {
   t.Log(video_tests)
}
