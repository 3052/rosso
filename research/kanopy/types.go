// types.go
package kanopy

// Session holds connection-related identifiers used across multiple API requests.
type Session struct {
   Authorization string
   UserId        int
   DomainId      int
}

// Video holds video-specific identifiers.
type Video struct {
   VideoId int
   Alias   string
}
