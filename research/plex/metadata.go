package plex

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type Metadata struct {
   MediaContainer MetadataContainer `json:"MediaContainer"`
}

type MetadataContainer struct {
   Metadata []MetadataItem `json:"Metadata"`
}

type MetadataItem struct {
   Media []MediaItem `json:"Media"`
}

type MediaItem struct {
   Part []MediaPart `json:"Part"`
}

type MediaPart struct {
   Id  string `json:"id"`
   Key string `json:"key"`
}

func GetMetadata(authToken string, ratingKey string) (*Metadata, error) {
   targetUrl := &url.URL{
      Scheme: "https",
      Host:   "vod.provider.plex.tv",
      Path:   "/library/metadata/" + ratingKey,
   }

   headers := map[string]string{
      "x-plex-token": authToken,
   }

   resp, err := maya.Get(targetUrl, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var metadata Metadata
   if err := json.NewDecoder(resp.Body).Decode(&metadata); err != nil {
      return nil, err
   }

   return &metadata, nil
}
