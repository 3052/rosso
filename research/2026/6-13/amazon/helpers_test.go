package amazon

import (
   "os"
   "path/filepath"
)

type authState struct {
   PublicCode  string `json:"public_code"`
   PrivateCode string `json:"private_code"`
}

type tokenState struct {
   AccessToken  string `json:"access_token"`
   RefreshToken string `json:"refresh_token"`
}

type actorTokenState struct {
   ActorId     string `json:"actor_id"`
   AccessToken string `json:"access_token"`
}

type envelopeState struct {
   PlaybackEnvelope string `json:"playback_envelope"`
}

func getStateFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_auth_state.json")
}

func getTokensFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_tokens.json")
}

func getActorTokensFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_actor_tokens.json")
}

func getEnvelopeFilePath() string {
   return filepath.Join(os.TempDir(), "amazon_envelope.json")
}
