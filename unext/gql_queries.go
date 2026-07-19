package unext

import (
   _ "embed"
   "strings"
)

// Pre-minified at package init; never recomputed.
var (
   minPlaylistQuery    = gqlMinify(rawPlaylistQuery)
   minAllEpisodesQuery = gqlMinify(rawAllEpisodesQuery)
   minVideoDetailQuery = gqlMinify(rawVideoDetailQuery)
)

//go:embed mad_all_episodes.graphql
var rawAllEpisodesQuery string

//go:embed mad_playlist.graphql
var rawPlaylistQuery string

//go:embed mad_video_detail.graphql
var rawVideoDetailQuery string

// gqlMinify collapses insignificant whitespace in a GraphQL operation
// string. None of the embedded queries contain string literals or
// comments, so a pure whitespace-collapser is sufficient.
func gqlMinify(s string) string {
   var b strings.Builder
   b.Grow(len(s))

   prevSpace := false
   for i := 0; i < len(s); i++ {
      c := s[i]
      if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
         prevSpace = true
         continue
      }
      if prevSpace {
         b.WriteByte(' ')
         prevSpace = false
      }
      b.WriteByte(c)
   }

   return b.String()
}
