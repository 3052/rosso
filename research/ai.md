# ai

1. package kanopy
2. Go language 1.26
3. one file per request
4. Variable Naming Rules:
   - If the name is 2 or more letters, DO NOT change it
   - If the name is 1 letter AND is a method receiver, DO NOT change it
   - If the name is 1 letter AND is NOT a method receiver, you MUST rename it to exactly one or two words
5. do not ignore errors
6. if two function inputs come from the same struct, pass the struct instead. if
   one function input comes from a struct, use the same name and type as the
   field or pass the struct itself
7. if passing a struct with two or more fields, pass a pointer
8. user will provide license payload
9. do not hard code authorization
10. decode HTTP responses as needed
11. do not use net/http
12. use 41.neocities.org/maya for HTTP
```
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
```
