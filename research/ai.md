# ai

1. package kanopy
2. Go language 1.26
3. one file per request
4. Variable Naming Rules:
   - **For names of 2 or more letters:** DO NOT change, lengthen, or expand
      standard idiomatic Go variable names
   - **For 1-letter method receivers:** DO NOT change them
   - **For 1-letter names that are NOT method receivers:** You MUST rename them
     to exactly one or two words
6. If a function requires two or more fields that come from the same struct, I
   must pass a pointer to that struct itself
7. If a function requires only one field from a struct, I should prefer passing
   just that single field. When I do this, the parameter's name and type must
   perfectly match the field's name and type
8. if passing a struct with two or more fields, pass a pointer
9. user will provide license payload
10. do not hard code authorization
11. decode HTTP responses as needed
12. do not use net/http
13. use 41.neocities.org/maya for HTTP
```
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
```
