package main

import (
   "io"
   "net/http"
   "net/url"
   "os"
   "strings"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "account.bellmedia.ca"
   req.URL.Path = "/api/login/v2.2"
   req.URL.Scheme = "https"
   req.Body = io.NopCloser(strings.NewReader(data))
   req.Header.Add("Authorization", "Basic Y3JhdmUtd2ViOmRlZmF1bHQ=")
   req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}

var data = url.Values{
   "grant_type":[]string{"refresh_token"},
   "refresh_token":[]string{"eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI2OTk2N2RhOWM5M2VlZjVkZjIwZjg3MTIiLCJ1c2VyX25hbWUiOiJjcmF2ZUB3b21lbi1hdC13b3JrLm9yZyIsImF0aSI6IjA1Y2NiNzVhLTkyNDYtNDliYS04MWRjLTAyYTE4NGUxYTc0OCIsInNjb3BlIjoiYWNjb3VudDp3cml0ZSBkZWZhdWx0IHBhc3N3b3JkX3Rva2VuIiwiY29udGV4dCI6eyJwcm9maWxlX2lkIjpudWxsLCJicmFuZF9pZHMiOlsiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTExIiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE0IiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE1Il19LCJleHAiOjE4MDU3NDQ0MjEsImlhdCI6MTc3NDE4NzQ5NSwidmVyc2lvbiI6IlYyIiwianRpIjoiMDIyNzdhOTEtYmMxYi00MTBmLTk2ZjItMmJiMzAyNTc5YWMxIiwiYXV0aG9yaXRpZXMiOlsiUkVHVUxBUl9VU0VSIl0sImNsaWVudF9pZCI6ImNyYXZlLXdlYiJ9.l_JQ4BfRLCw4EJ2Bj8gMDGDk02ro4Dp755bRB8uVKGjojr4BlwS3bNcOGeb6rxbiTejJFfnygwsLCn9aEI8RfX52MqwWAO3BgZtDKLkAQ7BXChkzdBEoFVJ0S0un5GPQrjs1x_kPhdO7WH9dOzuvpZA23qep07sFmdUK9dhtPcbSD44l6VzwK3zvYkV9_4rOVHwCi5OHhbLUM2rQdlNHJDQbuagPKOY8kzWIO9vufexhoP_yE5SghzdKKGKIVp2e4tdE6OnDFjQZa9pfNTjWbB5Dd2n5PWlokZkXYrWY_y1CSFequgWeebpPHzWOhsnMhma5R629KbT_cAIqgvly4g"},
}.Encode()

