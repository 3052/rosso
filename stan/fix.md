# fix

~~~
1 func (a *AppSession) FetchMedia(id int, quality, drm string) (*Media, error)

2 func (w *WebToken) FetchSession() (*AppSession, error)
3 func (a *ActivationCode) Token() (*WebToken, error)
4 func FetchActivationCode() (*ActivationCode, error)
~~~

## done

~~~
func (m *Media) BaseUrl(host string) (*url.URL, error)
func (*ActivationCode) CachePath() string
func (*AppSession) CachePath() string
func (*Media) CachePath() string
func (*WebToken) CachePath() string
func (a *ActivationCode) String() string
func (e *DTError) Error() string
func (e APIErrors) Error() string
func (m *Media) LicensePlayReady(data []byte) ([]byte, error)
func (m *Media) LicenseWidevine(data []byte) ([]byte, error)
func do(req *http.Request) (*http.Response, error)
~~~
