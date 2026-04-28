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
   target, err := url.Parse(targetUrl)
   if err != nil {
      return nil, err
   }

   if !strings.HasSuffix(target.Host, "rakuten.tv") {
      return nil, errors.New("invalid host")
   }

   segments := strings.Split(strings.Trim(target.Path, "/"), "/")
   if len(segments) == 0 || segments[0] == "" {
      return nil, errors.New("missing market code in path")
   }

   parsed := &ParsedUrl{
      MarketCode: segments[0],
   }

   if len(segments) == 3 {
      parsed.ContentType = segments[1]
      parsed.ContentId = segments[2]
   } else if len(segments) >= 5 && segments[1] == "player" {
      parsed.ContentType = segments[2]
      parsed.ContentId = segments[4]
   }

   query := target.Query()
   if parsed.ContentType == "" {
      parsed.ContentType = query.Get("content_type")
   }

   if parsed.ContentId == "" {
      if id := query.Get("content_id"); id != "" {
         parsed.ContentId = id
      } else if id := query.Get("tv_show_id"); id != "" {
         parsed.ContentId = id
         if parsed.ContentType == "" {
            parsed.ContentType = "tv_shows"
         }
      } else if id := query.Get("movie_id"); id != "" {
         parsed.ContentId = id
         if parsed.ContentType == "" {
            parsed.ContentType = "movies"
         }
      }
   }

   if parsed.MarketCode == "" || parsed.ContentType == "" || parsed.ContentId == "" {
      return nil, errors.New("could not extract all required components from url")
   }

   return parsed, nil
}

func (parsed *ParsedUrl) IsMovie() bool {
   return parsed.ContentType == "movies"
}

func (parsed *ParsedUrl) IsTvShow() bool {
   return parsed.ContentType == "tv_shows"
}
