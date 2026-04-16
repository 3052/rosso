package kanopy

const (
   BaseUrl   = "https://www.kanopy.com"
   UserAgent = "!"
   Xversion  = "!/!/!/!"
)

// Session represents an authenticated user context.
type Session struct {
   Jwt       string `json:"jwt"`
   VisitorId string `json:"visitorId"`
   UserId    int    `json:"userId"`
   UserRole  string `json:"userRole"`
}
