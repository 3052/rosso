package amazon

import (
   "encoding/json"
   "errors"
   "io"
   "net/http"
   "net/url"
   "time"
)

func (v *vod) playback_envelope() (string, error) {
   for _, item := range v.Resource.Header.TitleActionsV2.ListActionsView.PlaybackGroup.Items {
      if envelope := item.ItemReference.PlaybackExperienceMetadata.PlaybackEnvelope; envelope != "" {
         return envelope, nil
      }
   }
   return "", errors.New("envelope not found")
}

type vod struct {
   Resource struct {
      Header struct {
         TitleActionsV2 struct {
            ListActionsView struct {
               PlaybackGroup struct {
                  Items []struct {
                     ItemReference struct {
                        PlaybackExperienceMetadata struct {
                           PlaybackEnvelope string `json:"playbackEnvelope,omitempty"`
                        } `json:"playbackExperienceMetadata,omitempty"`
                     } `json:"itemReference,omitempty"`
                  } `json:"items,omitempty"`
               } `json:"playbackGroup,omitempty"`
            } `json:"listActionsView,omitempty"`
         } `json:"titleActionsV2,omitempty"`
      } `json:"header,omitempty"`
   } `json:"resource,omitempty"`
}

func create_vod() (*vod, error) {
   time.Sleep(time.Second)
   client := &http.Client{}
   reqURL := &url.URL{
      Scheme: "https",
      Host:   "abzq7aq4866p.na.api.amazonvideo.com",
      Path:   "/cdp/switchblade/android/getDataByJvmTransform/v1/dv-android/detail/vod/v1.kt",
   }
   q := url.Values{}
   q.Add("itemId", "amzn1.dv.gti.28b85d90-1338-720b-4be7-3247683a7624")
   q.Add("deviceId", "7f93fd3fd18646658ccd555423c8e5b8")
   q.Add("deviceTypeID", "A43PXU4ZN2AL1")
   q.Add("featureScheme", "mobile-android-features-v11.1")
   q.Add("isPlaybackEnvelopeSupported", "true")
   q.Add("swiftPriorityLevel", "critical")
   reqURL.RawQuery = q.Encode()
   req, err := http.NewRequest("GET", reqURL.String(), nil)
   if err != nil {
      return nil, err
   }
   req.Header.Add("Authorization", "Bearer Atna|EwMDIB9J3MGbf9E3tgBbHoUygRtpo-XCA6ykvhcHPAGs-ewGQZosfIejyc24UeQXMyuvXu40fd32xlYpS4hp5vlN5eEwG9dc_3Q-xDIWA-2APscp8pdvRJ09KnGvDqNtU4hI3Bm6f0lUI2mtd7YP86s9zCJb4SlWGRSDYxs_VZKBNXPpdoTsSwWRwju80wsWx1rF1Dr0KFCMsedsGUufl8kR-ah68mKvzdmlrsGdePXDviHZ9vmOOxCev-iUrAxDwj_T-etHp6KkpOyMjhgwVi4_Kx4zO8lCY_czqf0vkJ17hxCh2ZJLeaRoUQ1Fpkwtl8Xun7772kd7DJblCcTxwK_6VOnsgVo0ANo0d3PkPQHUqUYL6A")
   resp, err := client.Do(req)
   if err != nil {
      return nil, err
   }
   defer resp.Body.Close()
   data, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, err
   }
   if resp.StatusCode != http.StatusOK {
      return nil, errors.New(string(data))
   }
   var result vod
   err = json.Unmarshal(data, &result)
   if err != nil {
      return nil, err
   }
   return &result, nil
}
