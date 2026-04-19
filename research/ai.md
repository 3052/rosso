# ai

1. package kanopy
2. Go language 1.26
3. one file per request
4. after naming a variable, review the name. if the name is multiple bytes like
   "resp", carry on. if the name is a single byte like "a" and a receiver,
   carry on. if the name is a single byte and not a receiver, rename to one or
   two words.
5. do not ignore errors
6. if function input comes from a previous response field, it should be either
   use the same name and type of the field, or pass the struct itself
7. user will provide license payload
8. do not hard code authorization
9. decode HTTP responses as needed
10. do not use net/http
11. use 41.neocities.org/maya for HTTP
```
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
```
