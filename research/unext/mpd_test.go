package unext

import (
	"fmt"
	"os"
   "io"
   "net/http"
)

const mpdURLFile = "mpd_url.txt"

// SaveMPDURL writes the DASH MPD URL (with play_token appended) to a file.
func SaveMPDURL(playlist *PlaylistResponse) error {
	mpdURL, err := playlist.GetDASHPlaylistURL()
	if err != nil {
		return fmt.Errorf("getting DASH playlist URL: %w", err)
	}

	if err := os.WriteFile(mpdURLFile, []byte(mpdURL.String()), 0600); err != nil {
		return fmt.Errorf("writing MPD URL file: %w", err)
	}

	return nil
}

// LoadAndPrintMPD reads the MPD URL from file, fetches the content, and prints it.
func LoadAndPrintMPD(client *http.Client) error {
	urlBytes, err := os.ReadFile(mpdURLFile)
	if err != nil {
		return fmt.Errorf("reading MPD URL file: %w", err)
	}

	mpdURL := string(urlBytes)
	fmt.Printf("MPD URL: %s\n", mpdURL)

	req, err := http.NewRequest("GET", mpdURL, nil)
	if err != nil {
		return fmt.Errorf("creating MPD request: %w", err)
	}

	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:140.0) Gecko/20100101 Firefox/140.0")
	req.Header.Set("accept", "*/*")
	req.Header.Set("origin", "https://video.unext.jp")
	req.Header.Set("referer", "https://video.unext.jp/")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("fetching MPD content: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading MPD content: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected MPD status code %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("MPD Content:\n%s\n", string(body))
	return nil
}
