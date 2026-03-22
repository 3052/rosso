package main

import (
   "net/http"
   "net/url"
   "os"
)

func main() {
   var req http.Request
   req.Header = http.Header{}
   req.Method = "POST"
   req.URL = &url.URL{}
   req.URL.Host = "account.bellmedia.ca"
   req.URL.Path = "/api/magic-link/v2.1/generate"
   req.URL.Scheme = "https"
   req.Header.Add("Authorization", "Bearer " + bearer)
   resp, err := http.DefaultClient.Do(&req)
   if err != nil {
      panic(err)
   }
   err = resp.Write(os.Stdout)
   if err != nil {
      panic(err)
   }
}

const bearer = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJsb2NhbGl6YXRpb24iOiJlbi1DQSIsInVzZXJfbmFtZSI6ImNyYXZlQHdvbWVuLWF0LXdvcmsub3JnIiwiYnJhbmRfcG9saWNpZXMiOlsiY2FzdGluZzphaXJwbGF5IiwiY2FzdGluZzpjaHJvbWVjYXN0IiwiZGV2aWNlOjUiLCJvZmZsaW5lX2Rvd25sb2FkIiwicGxhdGZvcm06YW5kcm9pZCIsInBsYXRmb3JtOmFuZHJvaWRfdHYiLCJwbGF0Zm9ybTpmaXJlX3R2IiwicGxhdGZvcm06aGlzZW5zZSIsInBsYXRmb3JtOmlvcyIsInBsYXRmb3JtOmxnX3R2IiwicGxhdGZvcm06cHM1IiwicGxhdGZvcm06cm9rdSIsInBsYXRmb3JtOnNhbXN1bmdfdHYiLCJwbGF0Zm9ybTpzb255X3BzNCIsInBsYXRmb3JtOnN0YiIsInBsYXRmb3JtOnR2b3MiLCJwbGF0Zm9ybTp3ZWIiLCJwbGF0Zm9ybTp4MSIsInBsYXRmb3JtOnhib3hfb25lIiwicGxheWJhY2tfcXVhbGl0eTo0ayIsInBsYXliYWNrX3F1YWxpdHk6aGQiLCJwbGF5YmFja19xdWFsaXR5OnNkIiwic3RyZWFtX2NvbmN1cnJlbmN5OjQiLCJzdWJzY3JpcHRpb246Y3JhdmVfdG90YWwiLCJzdWJzY3JpcHRpb246Y3JhdmVwIiwic3Vic2NyaXB0aW9uOmNyYXZldHYiLCJzdWJzY3JpcHRpb246ZnJlZSIsInN1YnNjcmlwdGlvbjpzZSJdLCJjcmVhdGlvbl9kYXRlIjoxNzc0MTg3NDkxMDE3LCJhaXNfaWQiOm51bGwsImF1dGhvcml0aWVzIjpbIlJFR1VMQVJfVVNFUiJdLCJjbGllbnRfaWQiOiJjcmF2ZS13ZWIiLCJicmFuZF9pZCI6IjFkNzJkOTkwY2I3NjVkZTdlNDIxMTExMSIsImFjY291bnRfaWQiOiI2OTk2N2RhOWM5M2VlZjVkZjIwZjg3MTIiLCJwcm9maWxlX2lkIjpudWxsLCJzY29wZSI6WyJhY2NvdW50OndyaXRlIiwiZGVmYXVsdCIsIm9mZmxpbmVfZG93bmxvYWQ6MTAiLCJwYXNzd29yZF90b2tlbiIsInN1YnNjcmlwdGlvbjpjcmF2ZV90b3RhbCxjcmF2ZXAsY3JhdmV0dixmcmVlLHNlIl0sImV4cCI6MTc3NDIwMTg5MSwiaWF0IjoxNzc0MTg3NDkxLCJqdGkiOiJkZTk2MTdiMy02M2QwLTQxMmMtYTRkOC00YmNkNWZiOGZmODYiLCJhY2NvdW50Ijp7InN0YXR1cyI6IkFDVElWRSJ9fQ.EtC7c92yKrDP-RRAZITt8LrybWU2o_iAh6TOj3U8j68e9nU77VLlOQfic-20QsZwP-n9sY_Rk5ESyz-mB80HzYcf9d3vDlVv4emBedOo-ulVH8ZFOo7OOs6zpQZiD-TEEcS51sQ3tRgFEaNEio42_AwTSlMeCObUWrOpga_rSgb7T3yMdpp9V4ud8QBwhKuwzkQcPMnKLHtLh8lmT5XJJxKjnoNHsRw2UN06eePQ2sC8DMF-4bvsQu5t65erI_8zKgqqev_pgEL-NNJBdRk2Fhauz9XTirOv-MrCCcik9WVP3tM_KDd6dBWzcTwodfOPdzHpegHDvMcu3sQyvF1Q4w"

