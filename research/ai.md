# ai

1. Use the requested package name and Go 1.26.
2. One file per request.
3. Do not use `net/http`.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use `41.neocities.org/maya` for HTTP requests.
6. Avoid single-letter variables except for receivers. Receivers MUST use short, 1-2 letter names. Standard Go variable names (such as `resp`, `err`, `req`, `body`) are acceptable and preferred for all other variables.
7. Use `url.URL` struct literals for static URLs; for dynamic ones, do not combine `url.Parse` with `url.PathEscape` (use one or the other).
8. Never explicitly add `accept-encoding` headers. Never manually define or parameterize `user-agent` headers UNLESS the value is explicitly non-standard.
9. Avoid inline nested anonymous structs for JSON payloads to prevent duplicating type definitions. Use either exclusively `map[string]any` (including for nested fields) or exclusively explicitly named struct types; do not mix maps and structs within the same payload definition.
10. When constructing nested JSON request payloads, either use `map[string]any` exclusively, or if using explicit struct types, pass the component structs directly as function inputs.
11. If a struct has two or more fields, it must be passed as a pointer when used as a function parameter.
12. Function parameters that represent header values or specific fields must match the name of the target header or source field exactly. If a parameter represents a dynamically sourced value used to construct a header, name the parameter exactly after its source field rather than the resulting header name.
13. Do not parameterize static or hardcoded request payload fields or headers. Hardcode these values inside the function body so only dynamic data is passed as function input.
14. Pass individual extracted fields (such as primitive IDs) directly as function parameters instead of passing their parent domain structs ONLY if the function requires a single field from that struct. If multiple function inputs originate from the same parent domain struct, do not split them into individual parameters; pass the parent struct instead.
15. If a function requires a custom domain struct (such as a previous response struct) to supply its inputs, implement the function as a method on that struct rather than a standalone function. If inputs originate from multiple different domain structs, implement the function as a method on the most contextually relevant struct and pass the others as parameters.

~~~
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~
