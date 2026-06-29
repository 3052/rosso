package rtbf

import (
   "testing"
   "time"
)

var tests = []struct {
   category string
   path     string
   url      string
}{
   {
      category: "series",
      path:     "/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
      url:      "https://auvio.rtbf.be/media/the-durrells-une-famille-anglaise-a-corfou-the-durrells-une-famille-anglaise-a-corfou-s01-3351856",
   },
   {
      category: "films",
      path:     "/media/l-affaire-thomas-crown-avec-steve-mcqueen-et-faye-dunawa-3381405",
      url:      "https://auvio.rtbf.be/media/l-affaire-thomas-crown-avec-steve-mcqueen-et-faye-dunawa-3381405",
   },
}

func TestRtbf(t *testing.T) {
   for _, test := range tests {
      _, err := FetchAssetId(test.path)
      if err != nil {
         t.Fatal(err)
      }
      time.Sleep(time.Second)
   }
}
