package kanopy

const BaseURL = "https://www.kanopy.com"

// Client holds the authentication token and standard headers.
type Client struct {
   Token     string
   UserAgent string
   XVersion  string
}

// NewClient creates a new Kanopy config.
func NewClient() *Client {
   return &Client{
      UserAgent: "!",
      XVersion:  "!/!/!/!",
   }
}

// SetToken sets the Bearer JWT token for subsequent authenticated requests.
func (c *Client) SetToken(token string) {
   c.Token = token
}
