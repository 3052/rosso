package rakuten

import (
   "encoding/json"
   "net/url"

   "41.neocities.org/maya"
)

type UserSessionClassification struct {
   NumericalId int `json:"numerical_id"`
}

type UserSessionProfile struct {
   Classification    UserSessionClassification `json:"classification"`
   AudioLanguage     UserSessionLanguage       `json:"audio_language"`
   SubtitlesLanguage UserSessionLanguage       `json:"subtitles_language"`
}
