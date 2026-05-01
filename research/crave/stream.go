package crave

import (
   "encoding/json"
   "fmt"
   "net/url"

   "41.neocities.org/maya"
)

type Stream struct {
   Playback  string `json:"playback"`
   Trickplay string `json:"trickplay"`
}

func GetStream(token *ProfileToken, activePlayback *Playback) (*Stream, error) {
   endpoint := &url.URL{
      Scheme: "https",
      Host:   "stream.video.9c9media.com",
      Path:   fmt.Sprintf("/meta/content/%s/contentpackage/%d/destination/%d/platform/48", activePlayback.ContentId, activePlayback.ContentPackage.Id, activePlayback.DestinationId),
   }

   values := url.Values{}
   values.Set("filter", "ff")
   values.Set("format", "mpd")
   values.Set("hd", "true")
   values.Set("mcv", "true")
   values.Set("uhd", "true")
   endpoint.RawQuery = values.Encode()

   headers := map[string]string{
      "authorization": "Bearer " + token.AccessToken,
   }

   resp, err := maya.Get(endpoint, headers)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()

   activeStream := &Stream{}
   if err := json.NewDecoder(resp.Body).Decode(activeStream); err != nil {
      return nil, err
   }

   return activeStream, nil
}
