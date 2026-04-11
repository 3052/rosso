package crave

import (
   "fmt"
   "net/url"
   "strconv"
   "strings"
)

func ParseMedia(rawUrl string) (*Media, error) {
   parsedUrl, err := url.Parse(rawUrl)
   if err != nil {
      return nil, err
   }
   pathParts := strings.Split(parsedUrl.Path, "/")
   if len(pathParts) < 3 {
      return nil, fmt.Errorf("invalid url path structure")
   }
   urlType := pathParts[1]
   lastSegment := pathParts[len(pathParts)-1]
   dashIndex := strings.LastIndex(lastSegment, "-")
   if dashIndex == -1 {
      return nil, fmt.Errorf("id not found in url")
   }
   idString := lastSegment[dashIndex+1:]
   parsedId, err := strconv.Atoi(idString)
   if err != nil {
      return nil, fmt.Errorf("failed to parse id: %v", err)
   }
   m := &Media{}
   switch urlType {
   case "movie":
      m.Id = parsedId
   case "play":
      m.FirstContent.Id = parsedId
   default:
      return nil, fmt.Errorf("unknown media type: %s", urlType)
   }
   return m, nil
}

type Media struct {
   FirstContent struct {
      Id int `json:"id,string"`
   }
   Id int `json:"id,string"`
}
