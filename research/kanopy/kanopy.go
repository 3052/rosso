package kanopy

const (
   BaseURL   = "https://www.kanopy.com"
   UserAgent = "!"
   XVersion  = "!/!/!/!"
)

// Session represents an authenticated user context.
// It acts as the receiver for all authenticated Kanopy API methods.
type Session struct {
   JWT       string `json:"jwt"`
   VisitorID string `json:"visitorId"`
   UserID    int    `json:"userId"`
   UserRole  string `json:"userRole"`
}
