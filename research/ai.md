# ai

1. Use the requested package name and Go 1.26.
2. One file per request.
3. Do not use `net/http`.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use `41.neocities.org/maya` for HTTP requests.
6. Use `url.URL` struct literals for static URLs; for dynamic ones, do not combine `url.Parse` with `url.PathEscape` (use one or the other).
7. Never explicitly add `accept-encoding` headers. Do not parameterize headers that contain static or non-standard values; hardcode them directly in the request headers instead of passing them as function arguments.
8. If a struct has two or more fields, it must be passed as a pointer when used as a function parameter.
9. Do not create parameter structs to hold function arguments. Pass variables directly as standard function arguments.
10. Do not parameterize static, structural, or enum-like fields in JSON request bodies. Hardcode these constants directly into the payload generation.
11. If multiple input parameters naturally originate from the same previously defined struct (such as a response struct), pass that struct directly as the argument instead of extracting its individual fields. However, if only a single field from a struct is needed, pass that specific field directly using its exact name and type instead of passing the entire struct.

OUTPUT ONLY THE TYPE DEFINITIONS AND FUNCTION SIGNATURES

~~~
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~
