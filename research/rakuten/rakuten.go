package rakuten

import (
   "errors"
   "net/url"
   "strings"
)

type ParsedUrl struct {
   MarketCode  string
   ContentType string
   ContentId   string
}

func ParseUrl(targetUrl string) (*ParsedUrl, error) {
   parsed, err := url.Parse(targetUrl)
   if err != nil {
      return nil, err
   }

   if !strings.HasSuffix(parsed.Host, "rakuten.tv") {
      return nil, errors.New("invalid host")
   }

   segments := strings.Split(strings.Trim(parsed.Path, "/"), "/")
   if len(segments) == 0 || segments[0] == "" {
      return nil, errors.New("missing market code in path")
   }

   info := &ParsedUrl{
      MarketCode: segments[0],
   }

   if len(segments) == 3 {
      info.ContentType = segments[1]
      info.ContentId = segments[2]
   } else if len(segments) >= 5 && segments[1] == "player" {
      info.ContentType = segments[2]
      info.ContentId = segments[4]
   }

   query := parsed.Query()
   if info.ContentType == "" {
      info.ContentType = query.Get("content_type")
   }

   if info.ContentId == "" {
      if id := query.Get("content_id"); id != "" {
         info.ContentId = id
      } else if id := query.Get("tv_show_id"); id != "" {
         info.ContentId = id
         if info.ContentType == "" {
            info.ContentType = "tv_shows"
         }
      } else if id := query.Get("movie_id"); id != "" {
         info.ContentId = id
         if info.ContentType == "" {
            info.ContentType = "movies"
         }
      }
   }

   if info.MarketCode == "" || info.ContentType == "" || info.ContentId == "" {
      return nil, errors.New("could not extract all required components from url")
   }

   return info, nil
}

func (s *StreamInfo) GetManifest() (*url.URL, error) {
   return url.Parse(s.Url)
}
