package peacock

import (
   _ "embed"
   "encoding/json"
   "fmt"
   "log"
   "net/http"
   "net/url"
   "strings"
)

const idBase = "https://rango.id.peacocktv.com"

// doRequest logs the request method and URL, then sends the request
// using the provided http.Client.
func doRequest(client *http.Client, req *http.Request) (*http.Response, error) {
   log.Println(req.Method, req.URL)
   return client.Do(req)
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
   req.Header.Set("user-agent", "user-agent (Android; Build/user-agent)")
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
