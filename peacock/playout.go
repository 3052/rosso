package peacock

import (
   "bytes"
   "crypto/md5"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "net/url"
)

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
// embedded mTLS certificate. token supplies the userToken; providerVariantID
// identifies the asset to play. vcodec selects the requested video codec
// (e.g. "H264" or "H265"). protection selects the DRM system ("WIDEVINE" or
// "PLAYREADY").
func PlayoutVod(token *TokenResponse, providerVariantID, vcodec, protection string) (*PlayoutVodResponse, error) {
   if token == nil {
      return nil, fmt.Errorf("playout vod: nil token")
   }
   if token.UserToken == "" {
      return nil, fmt.Errorf("playout vod: empty userToken")
   }
   if providerVariantID == "" {
      return nil, fmt.Errorf("playout vod: empty providerVariantID")
   }
   if vcodec == "" {
      return nil, fmt.Errorf("playout vod: empty vcodec")
   }
   if protection == "" {
      return nil, fmt.Errorf("playout vod: empty protection")
   }
   body := playoutRequest{
      Device: playoutDevice{
         Capabilities: []playoutCapability{
            {
               Acodec:     "AAC",
               Container:  "ISOBMFF",
               Protection: protection,
               Transport:  "DASH",
               Vcodec:     vcodec,
            },
         },
         MaxVideoFormat: "UHD",
      },
      ProviderVariantID:            providerVariantID,
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
   req.Header.Set("x-skyott-usertoken", token.UserToken)
   req.Header.Set("Content-Type", "application/vnd.playvod.v1+json")
   req.Header.Set("Content-MD5", contentMD5)
   client, err := mtlsClient()
   if err != nil {
      return nil, fmt.Errorf("playout vod: %w", err)
   }
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

// AcquireLicense sends a Widevine license acquisition request to the licence acquisition URL
// and returns the raw license bytes. The licenceAcquisitionUrl is the full URL returned in
// the PlayoutVod response.
func (p *PlayoutVodResponse) AcquireLicense(challenge []byte) ([]byte, error) {
   if p == nil {
      return nil, fmt.Errorf("acquire license: nil playout")
   }
   if p.Protection.LicenceAcquisitionUrl == "" {
      return nil, fmt.Errorf("acquire license: empty licenceAcquisitionUrl")
   }
   if len(challenge) == 0 {
      return nil, fmt.Errorf("acquire license: empty challenge")
   }
   req, err := http.NewRequest(http.MethodPost, p.Protection.LicenceAcquisitionUrl, bytes.NewReader(challenge))
   if err != nil {
      return nil, fmt.Errorf("acquire license: create request: %w", err)
   }
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("acquire license: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("acquire license: bad status %d", resp.StatusCode)
   }

   license, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("acquire license: read body: %w", err)
   }

   return license, nil
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
   ProviderVariantID            string        `json:"providerVariantId"`
   PersonaParentalControlRating string        `json:"personaParentalControlRating"`
}

// playout.go
