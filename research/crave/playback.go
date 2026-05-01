package crave

import (
   "encoding/json"
   "net/url"
   "strconv"

   "41.neocities.org/maya"
)

type Playback struct {
   Playback  string `json:"playback"`
   Trickplay string `json:"trickplay"`
}

func GetPlayback(activeSession *Session, activeMedia *Media, available *AvailableContentPackage) (*Playback, error) {
   endpoint := url.URL{
      Scheme: "https",
      Host:   "stream.video.9c9media.com",
      Path:   "/meta/content/" + activeMedia.FirstContent.Id + "/contentpackage/" + strconv.Itoa(available.Id) + "/destination/" + strconv.Itoa(available.DestinationId) + "/platform/48",
   }

   values := url.Values{}
   values.Set("filter", "ff")
   values.Set("format", "mpd")
   values.Set("hd", "true")
   values.Set("mcv", "true")
   values.Set("uhd", "true")
   endpoint.RawQuery = values.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + activeSession.AccessToken,
   }

   resp, err := maya.Get(&endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   var targetPlayback Playback
   if err := json.NewDecoder(resp.Body).Decode(&targetPlayback); err != nil {
      return nil, err
   }

   return &targetPlayback, nil
}
