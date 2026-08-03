package peacock

import (
   "bytes"
   "crypto/md5"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "net/http"
   "net/url"
)

// PlayoutVodParams holds the parameters for a VOD playout request.
type PlayoutVodParams struct {
   UserToken                    string
   ContentID                    string
   ProviderVariantID            string
   ParentalControlPin           string
   PersonaParentalControlRating string
}

// PlayoutVodResponse is the response from POST /video/playouts/vod.
type PlayoutVodResponse struct {
   Asset struct {
      Endpoints []struct {
         Cdn string `json:"cdn"`
         Url string `json:"url"`
      } `json:"endpoints"`
   } `json:"asset"`
   Protection struct {
      LicenceAcquisitionUrl string `json:"licenceAcquisitionUrl"`
   } `json:"protection"`
   ErrorCode   string `json:"errorCode"`
   Description string `json:"description"`
}

// PlayoutVod requests a VOD playout URL from POST /video/playouts/vod using the
// embedded mTLS certificate.
func PlayoutVod(params *PlayoutVodParams) (*PlayoutVodResponse, error) {
   if params.UserToken == "" {
      return nil, fmt.Errorf("playout vod: empty userToken")
   }
   if params.ContentID == "" {
      return nil, fmt.Errorf("playout vod: empty contentID")
   }
   if params.ProviderVariantID == "" {
      return nil, fmt.Errorf("playout vod: empty providerVariantID")
   }
   client, err := mtlsClient()
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
   body := playoutRequest{
      Device: playoutDevice{
         Capabilities: []playoutCapability{
            {
               Acodec:     "AAC",
               Container:  "ISOBMFF",
               Protection: "WIDEVINE",
               Transport:  "DASH",
               Vcodec:     "H264",
            },
         },
         MaxVideoFormat: "HD",
      },
      ProviderVariantID:            params.ProviderVariantID,
      PersonaParentalControlRating: "9",
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
   req.Header.Set("x-skyott-usertoken", params.UserToken)
   req.Header.Set("Content-Type", "application/vnd.playvod.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   resp, err := doRequest(client, req)
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
   defer resp.Body.Close()

   var out PlayoutVodResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("playout vod: decode: %w", err)
   }
   if out.ErrorCode != "" {
      return nil, &out
   }
   return &out, nil
}

func (*PlayoutVodResponse) CachePath() string {
   return "rosso/peacock/PlayoutVodResponse"
}

// Error implements the error interface.
func (r *PlayoutVodResponse) Error() string {
   return r.ErrorCode + ": " + r.Description
}

// Fastly returns the parsed URL of the FASTLY CDN endpoint from the playout response.
func (r *PlayoutVodResponse) Fastly() (*url.URL, error) {
   for _, endpoint := range r.Asset.Endpoints {
      if endpoint.Cdn == "FASTLY" {
         parsed, err := url.Parse(endpoint.Url)
         if err != nil {
            return nil, fmt.Errorf("fastly: parse url: %w", err)
         }
         return parsed, nil
      }
   }
   return nil, fmt.Errorf("fastly cdn endpoint not found")
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

// playout.go
