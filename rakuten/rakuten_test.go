package rakuten

import "testing"

var tests = []struct{
   height int
   language []string
   url    string
   why string
}{
   {
      height: 1080,
      language: []string{"CAT", "ENG", "SPA"},
      url:    "https://rakuten.tv/es/movies/una-obra-maestra",
      why: "1080",
   },
   {
      url:    "https://rakuten.tv/ie?content_id=blair-witch&content_type=movies",
      height: 2160,
      language: []string{"ENG"},
      why: "2160",
   },
   {
      height: 1080,
      language: []string{"DEU", "ENG", "ITA", "POL"},
      url: "https://rakuten.tv/nl/player/movies/stream/made-in-america",
      why: "/player/movies/stream/",
   },
   {
      height: 1080,
      language: []string{"ENG"},
      url:    "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
      why: "?tv_show_id=",
   },
   {
      height: 1080,
      language: []string{"FRA"},
      url:    "https://rakuten.tv/fr/tv_shows/une-femme-d-honneur",
      why: "/tv_shows/",
   },
}

func TestLog(t *testing.T) {
   t.Log(address_tests, classification_tests)
}

