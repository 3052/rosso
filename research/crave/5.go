package crave

import (
   "net/http"
   "net/url"
)

func five() (*http.Response, error) {
   var req http.Request
   req.Header = http.Header{}
   req.URL = &url.URL{}
   req.URL.Host = "stream.video.9c9media.com"
   req.URL.Path = "/meta/content/938361/contentpackage/8143402/destination/1880/platform/1"
   value := url.Values{}
   value["format"] = []string{"mpd"}
   req.Header.Add("Authorization", "Bearer " + five_bearer)
   req.URL.RawQuery = value.Encode()
   req.URL.Scheme = "https"
   return http.DefaultClient.Do(&req)
}

// {
//   "sub": "69967da9c93eef5df20f8712",
//   "scope": "account:write default maturity:adult",
//   "iss": "https://account.bellmedia.ca",
//   "context": {
//     "profile_id": "69970efa73176d2b2e53a5a3",
//     "brand_ids": [
//       "1d72d990cb765de7e4211111",
//       "1d72d990cb765de7e4211114",
//       "1d72d990cb765de7e4211115"
//     ]
//   },
//   "exp": 1774201902,
//   "iat": 1774187502,
//   "version": "V2",
//   "jti": "9af409f3-5334-49d0-b976-be26ff8fd087",
//   "authorities": [
//     "REGULAR_USER"
//   ],
//   "client_id": "crave-web"
// }
const five_bearer = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJzdWIiOiI2OTk2N2RhOWM5M2VlZjVkZjIwZjg3MTIiLCJzY29wZSI6ImFjY291bnQ6d3JpdGUgZGVmYXVsdCBtYXR1cml0eTphZHVsdCIsImlzcyI6Imh0dHBzOi8vYWNjb3VudC5iZWxsbWVkaWEuY2EiLCJjb250ZXh0Ijp7InByb2ZpbGVfaWQiOiI2OTk3MGVmYTczMTc2ZDJiMmU1M2E1YTMiLCJicmFuZF9pZHMiOlsiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTExIiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE0IiwiMWQ3MmQ5OTBjYjc2NWRlN2U0MjExMTE1Il19LCJleHAiOjE3NzQyMDE5MDIsImlhdCI6MTc3NDE4NzUwMiwidmVyc2lvbiI6IlYyIiwianRpIjoiOWFmNDA5ZjMtNTMzNC00OWQwLWI5NzYtYmUyNmZmOGZkMDg3IiwiYXV0aG9yaXRpZXMiOlsiUkVHVUxBUl9VU0VSIl0sImNsaWVudF9pZCI6ImNyYXZlLXdlYiJ9.ipnWI2we9vx6wKx3u8ZJuqDjZ46nr7c_vYqn6u28IOSvLfvEuWfWiE8C2UQxVzUzYZeSiQRw-vpbzgKE-KMR4ZfBSlU2f3AcvP6wcDkBGaXkMvfao-dIbbeUHDMjeX1seCE2LzJ0N73MZ4503NZ6heHmCphkZN3wtwsgo6ZnejWD5uT3JkN-rPGcJr_y17VMf87RfuI4OG4qJH5x4NPIDcfy8uF4xcXVU6nFi6clewEb5ivV5aYMXb78lZhyCBYlL2v_DDcJ5jTgfRyThvpYnPpB6G5898JJGE-LPJUJHb17q1_7u6ocWlDPEifvNG9-c0AyG1hgLVJl28Yr4IGGVg"

