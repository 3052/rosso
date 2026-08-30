package peacock

import (
   "bytes"
   "crypto/tls"
   _ "embed"
   "encoding/json"
   "fmt"
   "io"
   "log"
   "net/http"
   "net/url"
   "strings"
)

const idBase = "https://rango.id.peacocktv.com"

const playBase = "https://play.clients.peacocktv.com"

const user_agent = "user-agent (Android; Build/user-agent)"

//go:embed cert.pem
var certPEM []byte

//go:embed key.pem
var keyPEM []byte

// auth.go
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

// CategoryError is a Siren category-level error.
type CategoryError struct {
   Target  string `json:"target"`
   Message string `json:"message"`
   Code    string `json:"code"`
}

// Errors holds the field and category errors from a Siren response.
// It implements the error interface when it contains at least one
// error with a non-empty message.
type Errors struct {
   FieldErrors    []FieldError    `json:"fieldErrors"`
   CategoryErrors []CategoryError `json:"categoryErrors"`
}

// Error implements the error interface. It returns all field and
// category errors with non-empty messages, joined by "; ".
func (e *Errors) Error() string {
   var errs []string
   for _, fe := range e.FieldErrors {
      if fe.Message != "" {
         errs = append(errs, fe.Code+": "+fe.Message)
      }
   }
   for _, ce := range e.CategoryErrors {
      if ce.Message != "" {
         errs = append(errs, ce.Code+": "+ce.Message)
      }
   }
   return strings.Join(errs, "; ")
}

// HasErrors returns true if the Errors value contains at least one
// field or category error with a non-empty message.
func (e *Errors) HasErrors() bool {
   for _, fe := range e.FieldErrors {
      if fe.Message != "" {
         return true
      }
   }
   for _, ce := range e.CategoryErrors {
      if ce.Message != "" {
         return true
      }
   }
   return false
}

// FieldError is a Siren field-level validation error.
type FieldError struct {
   Target  string `json:"target"`
   Message string `json:"message"`
   Code    string `json:"code"`
}

// OAuthAuthorizeResponse is the Siren JSON response from GET /oauth/authorize/service/international.
type OAuthAuthorizeResponse struct {
   Class      []string `json:"class"`
   Properties struct {
      EventType   string `json:"eventType"`
      Href        string `json:"href"`
      Type        string `json:"type"`
      AccessToken string `json:"access_token"`
   } `json:"properties"`
}

// OAuthAuthorize follows the OAuth authorize redirect to obtain the access token.
// It requires the skyCEsidmesso01 cookie set by SignIn.
func OAuthAuthorize(cookie *http.Cookie) (*OAuthAuthorizeResponse, error) {
   if cookie == nil {
      return nil, fmt.Errorf("oauth authorize: nil cookie")
   }
   if cookie.Value == "" {
      return nil, fmt.Errorf("oauth authorize: empty cookie value")
   }
   params := url.Values{}
   params.Set("client_id", "nbcu_tvclient")
   params.Set("response_type", "token")
   oauthUrl := idBase + "/oauth/authorize/service/international?" + params.Encode()
   req, err := http.NewRequest(http.MethodGet, oauthUrl, nil)
   if err != nil {
      return nil, fmt.Errorf("oauth authorize: create request: %w", err)
   }
   req.AddCookie(cookie)
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-provider", "NBCU")
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("oauth authorize: %w", err)
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return nil, fmt.Errorf("oauth authorize: bad status %d", resp.StatusCode)
   }

   var out OAuthAuthorizeResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("oauth authorize: decode: %w", err)
   }

   if out.Properties.AccessToken == "" {
      return nil, fmt.Errorf("oauth authorize: access_token not found in response")
   }

   return &out, nil
}

func (*OAuthAuthorizeResponse) CachePath() string {
   return "rosso/peacock/oauth_authorize"
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

   license, err := io.ReadAll(resp.Body)
   if err != nil {
      return nil, fmt.Errorf("acquire license: read body: %w", err)
   }

   // The endpoint returns a JSON error body instead of a non-200 status
   // when the request is rejected (e.g. "Unsupported browser/client").
   var lr licenseErrorResponse
   if err := json.Unmarshal(license, &lr); err == nil && lr.ErrorCode != "" {
      return nil, fmt.Errorf("acquire license: %s: %s", lr.ErrorCode, lr.Description)
   }

   return license, nil
}

// SignInParams holds the credentials for the sign-in request.
type SignInParams struct {
   UserIdentifier string
   Password       string
   RememberMe     bool
}

// SignInResponse is the Siren JSON response from POST /signin/service/international.
// Cookie is the skyCEsidmesso01 cookie from the Set-Cookie header; it is required
// as input to OAuthAuthorize.
type SignInResponse struct {
   Cookie     *http.Cookie `json:"-"`
   Class      []string     `json:"class"`
   Properties struct {
      TargetApi string `json:"targetApi"`
      Type      string `json:"type"`
      EventType string `json:"eventType"`
      Href      string `json:"href"`
      Data      struct {
         DeviceID string `json:"deviceid"`
      } `json:"data"`
      Errors Errors `json:"errors"`
   } `json:"properties"`
}

// SignIn authenticates with email and password via POST /signin/service/international.
// The skyCEsidmesso01 cookie from the response is returned in SignInResponse.Cookie
// for use by OAuthAuthorize. The device id is available in Properties.Data.DeviceID for
// use by ExchangeToken.
func SignIn(params *SignInParams) (*SignInResponse, error) {
   if params.UserIdentifier == "" {
      return nil, fmt.Errorf("sign in: empty userIdentifier")
   }
   if params.Password == "" {
      return nil, fmt.Errorf("sign in: empty password")
   }
   signInUrl := idBase + "/signin/service/international"
   form := url.Values{}
   form.Set("password", params.Password)
   form.Set("userIdentifier", params.UserIdentifier)
   req, err := http.NewRequest(http.MethodPost, signInUrl, strings.NewReader(form.Encode()))
   if err != nil {
      return nil, fmt.Errorf("sign in: create request: %w", err)
   }
   req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
   req.Header.Set("x-skyott-territory", "US")
   req.Header.Set("x-skyott-provider", "NBCU")
   req.Header.Set("user-agent", user_agent)
   resp, err := doRequest(http.DefaultClient, req)
   if err != nil {
      return nil, fmt.Errorf("sign in: %w", err)
   }
   defer resp.Body.Close()

   var out SignInResponse
   if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
      return nil, fmt.Errorf("sign in: decode: %w", err)
   }

   if out.Properties.Errors.HasErrors() {
      return nil, fmt.Errorf("sign in: %w", &out.Properties.Errors)
   }

   for _, cookie := range resp.Cookies() {
      if cookie.Name == "skyCEsidmesso01" {
         out.Cookie = cookie
         break
      }
   }
   if out.Cookie == nil {
      return nil, fmt.Errorf("sign in: skyCEsidmesso01 cookie not found in response")
   }

   return &out, nil
}

func (*SignInResponse) CachePath() string {
   return "rosso/peacock/SignInResponse"
}

// licenseErrorResponse is the error body returned by the licence
// acquisition endpoint when the request is rejected.
type licenseErrorResponse struct {
   ErrorCode   string `json:"errorCode"`
   Description string `json:"description"`
}

// license.go
