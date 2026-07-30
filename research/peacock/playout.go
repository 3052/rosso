// playout.go
package peacock

import (
   "bytes"
   "crypto/md5"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
)

// PlayoutVodParams holds the parameters for a VOD playout request.
type PlayoutVodParams struct {
   UserToken                    string
   ContentID                    string
   ProviderVariantID            string
   ParentalControlPin           string
   PersonaParentalControlRating string
}

type playoutCapability struct {
   Acodec     string `json:"acodec"`
   Container  string `json:"container"`
   Protection string `json:"protection"`
   Transport  string `json:"transport"`
   Vcodec     string `json:"vcodec"`
}

type playoutClient struct {
   ThirdParties   []string `json:"thirdParties"`
   VariantCapable bool     `json:"variantCapable"`
}

type playoutDevice struct {
   Capabilities          []playoutCapability `json:"capabilities"`
   MaxVideoFormat        string              `json:"maxVideoFormat"`
   SupportedColourSpaces []string            `json:"supportedColourSpaces"`
   Model                 string              `json:"model"`
   HdcpEnabled           bool                `json:"hdcpEnabled"`
}

type playoutRequest struct {
   Device                       playoutDevice `json:"device"`
   Client                       playoutClient `json:"client"`
   ContentID                    string        `json:"contentId"`
   ProviderVariantID            string        `json:"providerVariantId"`
   ParentalControlPin           string        `json:"parentalControlPin"`
   PersonaParentalControlRating string        `json:"personaParentalControlRating"`
}

// PlayoutVod requests a VOD playout URL from POST /video/playouts/vod using the
// embedded mTLS certificate. The response is returned as raw JSON since the
// structure may vary.
func (c *Client) PlayoutVod(params PlayoutVodParams) (json.RawMessage, error) {
   if params.UserToken == "" {
      return nil, fmt.Errorf("playout vod: empty userToken")
   }
   if params.ContentID == "" {
      return nil, fmt.Errorf("playout vod: empty contentID")
   }
   if params.ProviderVariantID == "" {
      return nil, fmt.Errorf("playout vod: empty providerVariantID")
   }

   if params.ParentalControlPin == "" && params.PersonaParentalControlRating == "" {
      params.PersonaParentalControlRating = "9"
   }

   client, err := mtlsClient(c.HTTP.Timeout)
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }

   body := playoutRequest{
      Device: playoutDevice{
         Capabilities: []playoutCapability{
            {Acodec: "AAC", Container: "ISOBMFF", Protection: "WIDEVINE", Transport: "DASH", Vcodec: "H264"},
            {Acodec: "AAC", Container: "ISOBMFF", Protection: "NONE", Transport: "DASH", Vcodec: "H264"},
         },
         MaxVideoFormat:        "HD",
         SupportedColourSpaces: []string{"SDR"},
         Model:                 "ANDROIDTV",
         HdcpEnabled:           false,
      },
      Client: playoutClient{
         ThirdParties:   []string{"FREEWHEEL", "MEDIATAILOR", "CONVIVA"},
         VariantCapable: true,
      },
      ContentID:                    params.ContentID,
      ProviderVariantID:            params.ProviderVariantID,
      ParentalControlPin:           params.ParentalControlPin,
      PersonaParentalControlRating: params.PersonaParentalControlRating,
   }

   raw, err := json.Marshal(body)
   if err != nil {
      return nil, fmt.Errorf("playout vod: marshal: %w", err)
   }

   hash := md5.Sum(raw)
   contentMD5 := hex.EncodeToString(hash[:])

   req, err := http.NewRequest(http.MethodPost, playBase+"/video/playouts/vod", bytes.NewReader(raw))
   if err != nil {
      return nil, fmt.Errorf("playout vod: create request: %w", err)
   }

   req.Header.Set("Accept", "application/vnd.playvod.v1+json")
   req.Header.Set("Content-Type", "application/vnd.playvod.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   req.Header.Set("Origin", "https://tv.clients.peacocktv.com")
   req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 12; sdk_gphone64_x86_64 Build/SE1A.220826.008; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/91.0.4472.114 Mobile Safari/537.36")
   req.Header.Set("x-skyott-ab-atom", "BrandIcons:variant_brand_icons;VisionDeCampoTileImagery:variation__new_tile_image_;pzEntitlementOrder:pzeoFree1")
   req.Header.Set("x-skyott-ab-suggest", "olympicsSportsBoost:control")
   req.Header.Set("x-skyott-activeterritory", "US")
   req.Header.Set("x-skyott-broadcastregions", "INPATTERN_US_CENTRAL")
   req.Header.Set("x-skyott-coppa", "false")
   req.Header.Set("x-skyott-device", "TV")
   req.Header.Set("x-skyott-language", "en-US")
   req.Header.Set("x-skyott-pinoverride", "false")
   req.Header.Set("x-skyott-platform", "ANDROIDTV")
   req.Header.Set("x-skyott-proposition", "NBCUOTT")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-usertoken", params.UserToken)

   resp, err := client.Do(req)
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("playout vod: bad status %d", resp.StatusCode)
   }

   var out json.RawMessage
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("playout vod: decode: %w", err)
   }
   return out, nil
}
