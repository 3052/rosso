package rakuten

import "testing"

var tests = []struct {
   height   int
   language []string
   url      string
   why      []string
}{
   {
      height:   1080,
      language: []string{"CAT", "ENG", "SPA"},
      url:      "https://rakuten.tv/es/movies/una-obra-maestra",
      why:      []string{"1080", "language", "/movies/"},
   },
   {
      height:   2160,
      language: []string{"ENG"},
      url:      "https://rakuten.tv/ie?content_id=blair-witch&content_type=movies",
      why:      []string{"2160"},
   },
   {
      height:   1080,
      language: []string{"DEU", "ENG", "ITA", "POL"},
      url:      "https://rakuten.tv/nl/player/movies/stream/made-in-america",
      why:      []string{"/player/movies/stream/"},
   },
   {
      height:   2160,
      language: []string{"CAT", "ENG", "FRA", "ITA", "POR", "SPA"},
      url:      "https://rakuten.tv/nl/tv_shows/matchday",
      why:      []string{"height", "language", "/tv_shows/"},
   },
   {
      height:   1080,
      language: []string{"ENG"},
      url:      "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
      why:      []string{"tv_show_id="},
   },
}

func TestLog(t *testing.T) {
   t.Log(tests)
}
