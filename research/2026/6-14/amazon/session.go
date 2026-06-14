package amazon

import (
   "bytes"
   "crypto/rand"
   "encoding/json"
   "fmt"
   "io"
   "net/http"
   "time"
)

// generateUUID creates a pseudo-random UUID for the session.
func generateUUID() string {
   b := make([]byte, 16)
   _, _ = rand.Read(b)
   // Set version (4) and variant (RFC4122)
   b[6] = (b[6] & 0x0f) | 0x40
   b[8] = (b[8] & 0x3f) | 0x80
   return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

// StartSession initiates the playback session using the actor access token and playback envelope.
func StartSession(actorAccessToken, playbackEnvelope string) error {
   url := "https://ab8mt4dd97et.na.api.amazonvideo.com/cdp/playback/pes/StartSession"

   req, err := http.NewRequest("POST", url, nil)
   if err != nil {
      return err
   }

   q := req.URL.Query()
   q.Add("deviceTypeID", "A2SNKIF736WF4T")
   q.Add("deviceID", "uuidcbb2f9705f13437e9e515622dce02106")
   q.Add("firmware", "1")
   q.Add("version", "1")
   req.URL.RawQuery = q.Encode()

   payload := map[string]interface{}{
      "playbackEnvelope": playbackEnvelope,
      "streamInfo": map[string]interface{}{
         "eventType":        "START",
         "streamIntent":     "AUTOPLAY",
         "streamUpdateTime": time.Now().UTC().Format("2006-01-02T15:04:05.000Z"),
         "vodProgressInfo": map[string]string{
            "currentProgressTime": "PT0S",
            "timeFormat":          "ISO8601DURATION",
         },
      },
      "userWatchSessionId": generateUUID(),
   }

   body, err := json.Marshal(payload)
   if err != nil {
      return err
   }

   req.Body = io.NopCloser(bytes.NewBuffer(body))
   req.ContentLength = int64(len(body))

   req.Header.Set("User-Agent", "Android/google/sdk_gphone_x86/generic_x86_arm:11/RSR1.240422.006/12134477:userdebug/dev-keys, Ignition X/15.5.2026042820-android, Google")
   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Accept", "application/json")
   req.Header.Set("Authorization", "Bearer "+actorAccessToken)
   req.Header.Set("x-client-app", "avlrc")

   client := &http.Client{}
   resp, err := client.Do(req)
   if err != nil {
      return err
   }
   defer resp.Body.Close()

   if resp.StatusCode != http.StatusOK {
      return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
   }

   return nil
}
