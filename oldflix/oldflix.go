package oldflix

import (
   "errors"
   "net/url"
)

func (w *Watch) GetManifest() (*url.URL, error) {
   return url.Parse(w.Playlist[0].File)
}

type Watch struct {
   Message  string
   Playlist []struct {
      File string
   }
}

const azure = "oldflix-api.azurewebsites.net"

func (b *Browse) GetOriginal() (*Track, error) {
   for _, track_data := range b.Movie.Tracks {
      if track_data.Lang == "Original" {
         return &track_data, nil
      }
   }
   return nil, errors.New("track with language 'Original' not found")
}

type Browse struct {
   Id    string
   Movie struct {
      Id     string
      Tracks []Track
   }
}

type Login struct {
   Status int
   Token  string
}

type Track struct {
   Id   string
   Lang string
   Lnk  string
}
