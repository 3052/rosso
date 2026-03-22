package crave

import (
   "io"
   "net/http"
   "net/url"
   "strings"
)

func Four() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "account.bellmedia.ca"
   req.URL.Path = "/api/login/v2.2"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(four_data))
   req.Header.Add("Authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
   return http.DefaultClient.Do(&req)
}

var four_data = url.Values{
   "grant_type":[]string{"refresh_token"},
   "refresh_token":[]string{"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI2OTk2N2RhOWM5M2VlZjVkZjIwZjg3MTIiLCJ1c2VyX25hbWUiOiJjcmF2ZUB3b21lbi1hdC13b3JrLm9yZyIsImF0aSI6IjQwNTAwMWZhLWFiMWItNDQ4Zi1hMWQ3LWU4OTE1NWNkNzIxNCIsInNjb3BlIjoiYWNjb3VudDp3cml0ZSBkZWZhdWx0IG1hdHVyaXR5OmFkdWx0IiwiY29udGV4dCI6eyJwcm9maWxlX2lkIjoiNjk5NzBlZmE3MzE3NmQyYjJlNTNhNWEzIiwiYnJhbmRfaWRzIjpbIjFkNzJkOTkwY2I3NjVkZTdlNDIxMTExMSIsIjFkNzJkOTkwY2I3NjVkZTdlNDIxMTExNCIsIjFkNzJkOTkwY2I3NjVkZTdlNDIxMTExNSJdfSwiZXhwIjoxODA1NzQ0NDIxLCJpYXQiOjE3NzQxODc1MDAsInZlcnNpb24iOiJWMiIsImp0aSI6IjM0NGNmYmI2LWEyY2EtNGZlNC1iODUzLWZkMjA5YWNlMjVkNyIsImF1dGhvcml0aWVzIjpbIlJFR1VMQVJfVVNFUiJdLCJjbGllbnRfaWQiOiJjcmF2ZS13ZWIifQ.vSBGMMA4fzpTzs4DxcOAB9iTGVGxFAX7mZhIguwpKCcfwNWEPbI52EIXD8h8zJEL93P2ecGFyXuy_95GnryIGSC8aC3PEsLvfocRInUV2neHU9TdnKICkmQfd9Mxj4Mf4VsB2S-7t8IJ9rWXkycSyQ8tQA1uM68PzKiHZx422b_VHZ6j2o3X79U1ujagtHNYNm59P9WxchwtJR8e786kVbdenZyODRmqajzPvSWWgve6oDxBTCmzLwbllss-ty-FNmC_HSmcMBtMmSPzEO4zwhw-iFtk1YSbvfl32S_fVNcvOzHljzlq7axCjhmIKAZDxFPcmDITIQ79QmymTWP61g"},
}.Encode()

