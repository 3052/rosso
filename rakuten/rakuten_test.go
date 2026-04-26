package rakuten

import "testing"

var classification_tests = []struct {
   height int
   url    string
}{
   {
      url:    "https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run",
      height: 2160,
   },
   {
      url:    "https://rakuten.tv/es/movies/una-obra-maestra",
      height: 1080,
   },
   {
      url:    "https://rakuten.tv/fr/tv_shows/une-femme-d-honneur",
      height: 1080,
   },
   {
      url:    "https://rakuten.tv/ie?content_id=blair-witch&content_type=movies",
      height: 2160,
   },
   {
      url:    "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
      height: 1080,
   },
   {
      url:    "https://rakuten.tv/pt/movies/bound",
      height: 1080,
   },
   {
      url:    "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
      height: 1080,
   },
}

func TestLog(t *testing.T) {
   t.Log(address_tests, classification_tests)
}

var address_tests = []struct {
   format string
   url    string
}{
   {
      format: "/movies/",
      url:    "https://rakuten.tv/nl/movies/made-in-america",
   },
   {
      format: "/player/movies/stream/",
      url:    "https://rakuten.tv/nl/player/movies/stream/made-in-america",
   },
   {
      format: "/tv_shows/",
      url:    "https://rakuten.tv/fr/tv_shows/une-femme-d-honneur",
   },
   {
      format: "?content_id=",
      url:    "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   },
   {
      format: "?tv_show_id=",
      url:    "https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink",
   },
}
