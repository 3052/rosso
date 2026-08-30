package peacock

import (
   "bytes"
   "crypto/md5"
   "crypto/tls"
   _ "embed"
   "encoding/hex"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "time"
)

const playBase = "https://play.clients.peacocktv.com"

//go:embed cert.pem
var certPEM []byte

//go:embed key.pem
var keyPEM []byte

// doRequest logs the request method and URL, then sends the request
// using the provided http.Client.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return client.Do(req)
}

// mtlsClient returns an *http.Client configured with the embedded mTLS
// certificate, ProxyFromEnvironment, and the given timeout.
func mtlsClient() (*http.Client, error) {
   cert, err := tls.X509KeyPair(certPEM, keyPEM)
   if err != nil {
      return nil, err
   }

   return &http.Client{
      Transport: &http.Transport{
         Proxy: http.ProxyFromEnvironment,
         TLSClientConfig: &tls.Config{
            Certificates: []tls.Certificate{cert},
         },
      },
   }, nil
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
// embedded mTLS certificate. token supplies the userToken; providerVariantID
// identifies the asset to play. vcodec selects the requested video codec
// (e.g. "H264" or "H265"). protection selects the DRM system ("WIDEVINE" or
// "PLAYREADY"). colourSpace selects the requested transfer range ("SDR",
// "HDR10" or "DolbyVision").
func PlayoutVod(
   token *TokenResponse, providerVariantID, vcodec, protection, colourSpace string,
) (*PlayoutVodResponse, error) {
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
         MaxVideoFormat:        "UHD",
         SupportedColourSpaces: []string{colourSpace},
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

// Error implements the error interface.
func (r *PlayoutVodResponse) Error() string {
   return r.ErrorCode + ": " + r.Description
}

// TokenResponse is the response from POST /auth/throttled/tokens.
// On a non-2xx response the server populates ErrorCode/Description
// instead of the token fields, in which case TokenResponse doubles
// as the returned error.
type TokenResponse struct {
   UserToken                     string    `json:"userToken"`
   TokenExpiryTime               time.Time `json:"tokenExpiryTime"`
   RecommendedTokenReacquireTime time.Time `json:"recommendedTokenReacquireTime"`
   ErrorCode                     string    `json:"errorCode"`
   Description                   string    `json:"description"`
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
