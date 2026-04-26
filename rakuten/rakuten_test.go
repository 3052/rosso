package rakuten

import "testing"

var classification_tests = map[int][]string{
   34: {"https://rakuten.tv/pt/movies/bound"},
   40: {"https://rakuten.tv/ie/movies/miss-sloane"},
   45: {"https://rakuten.tv/es/movies/una-obra-maestra"},
   60: {"https://rakuten.tv/uk?content_type=tv_shows&tv_show_id=clink"},
   61: {"https://rakuten.tv/pl?content_type=movies&content_id=ad-astra"},
   68: {
      "https://rakuten.tv/fr?content_type=movies&content_id=michael-clayton",
      "https://rakuten.tv/nl?content_type=movies&content_id=made-in-america",
   },
   70: {"https://rakuten.tv/se?content_type=movies&content_id=i-heart-huckabees"},
   83: {"https://rakuten.tv/cz?content_type=movies&content_id=transvulcania-the-people-s-run"},
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
