package amazon

import "testing"

func TestWidevineL3(t *testing.T) {
   runDeviceCombinations(
      t,
      "Widevine L3",
      `C:\Users\Steven\AppData\Local\L3`,
      "Widevine",
   )
}

func TestPlayReadySL2000(t *testing.T) {
   runDeviceCombinations(
      t,
      "PlayReady SL2000",
      `C:\Users\Steven\AppData\Local\SL2000`,
      "PlayReady",
   )
}

func TestPlayReadySL3000(t *testing.T) {
   runDeviceCombinations(
      t,
      "PlayReady SL3000",
      `C:\Users\Steven\AppData\Local\SL3000`,
      "PlayReady",
   )
}
