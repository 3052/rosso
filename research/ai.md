# ai

note I would prefer simpler rules, but the AI is too fucking stupid to understand
anything simpler than this:

1. package kanopy
2. Go language 1.26
3. one file per request
4. Variable Naming Rules:
   - **For names of 2 or more letters:** DO NOT change, lengthen, or expand
      standard idiomatic Go variable names
   - **For 1-letter method receivers:** DO NOT change them
   - **For 1-letter names that are NOT method receivers:** You MUST rename them
     to exactly one or two words
5. When an API request payload requires a mix of static/hardcoded structural
   fields and dynamic user data, functions must only accept the parameters or the
   specific inner struct containing the dynamic user data (as a pointer if it
   contains two or more fields). Do not force the caller to instantiate top-level
   wrapper structs just to satisfy the full JSON payload shape
6. user will provide license payload
7. do not hard code authorization
8. decode HTTP responses as needed
9. do not use net/http
10. use 41.neocities.org/maya for HTTP
```
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
```
