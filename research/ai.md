# ai

1. Use the requested package name and Go 1.26.
2. One file per request.
3. Do not use `net/http`.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use `41.neocities.org/maya` for HTTP requests.
6. Use `url.URL` struct literals for static URLs. Do not use `url.Parse` on a URL that is known at compile time. For dynamic URLs, do not combine `url.Parse` with `url.PathEscape` (use one or the other). Never construct `RawQuery` via string concatenation; always use `url.Values` and its `Encode()` method to generate query parameters safely.
7. Do not use single-letter variables (e.g., do not use `u` for URLs). Use a single word instead. If and only if a single word is not clear, use two words. Do not replace any other variables unless it's the same situation.
8. Never explicitly add standard or automatically generated headers like `accept-encoding` or default `user-agent` strings. Only set the `user-agent` key if its value is non-standard. Do not parameterize headers that contain static or non-standard values; hardcode them directly in the request headers instead of passing them as function arguments.
9. Do not parameterize static, structural, or enum-like fields in JSON request bodies. Hardcode these constants directly into the payload generation.
10. Never use anonymous structs. Either define an explicit named type or use a map.
11. When constructing JSON payloads, do not mix structs and maps. Choose one approach or the other: either use a fully defined hierarchy of named structs, or use maps entirely. Do not embed a struct inside a map.
12. Do not use any double capitals (consecutive uppercase letters) in identifiers, including acronyms.
13. Function input variables must match the corresponding JSON/request field names exactly (e.g., use `email` and `password`). If that is not possible, pass the parent struct instead.

~~~
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done
