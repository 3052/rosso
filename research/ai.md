# ai

1. Use the requested package name and Go 1.26.
2. One file per request.
3. Do not use `net/http`.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use `41.neocities.org/maya` for HTTP requests.
6. Use `url.URL` struct literals for static URLs. Do not use `url.Parse` on a URL that is known at compile time. For dynamic URLs, do not combine `url.Parse` with `url.PathEscape` (use one or the other). Never construct `RawQuery` via string concatenation; always use `url.Values` and its `Encode()` method to generate query parameters safely.
7. Do not use single-letter variables. Use a single word instead for variables. If and only if a single word is not clear, use two words. This rule applies ONLY to variables, not function names.
8. Never explicitly add standard or automatically generated headers like `accept-encoding` or default `user-agent` strings. Only set the `user-agent` key if its value is non-standard.
9. Do not parameterize static, structural, dummy, or enum-like values in query parameters, headers, or JSON request bodies. Hardcode these constants directly into the request construction instead of exposing them as function arguments.
10. Never use anonymous structs. Either define an explicit named type or use a map.
11. When constructing JSON payloads, do not mix structs and maps. Choose one approach or the other: either use a fully defined hierarchy of named structs, or use maps entirely. Do not embed a struct inside a map.
12. Do not use any double capitals (consecutive uppercase letters) in identifiers, including acronyms.
13. Function input variables must match the corresponding JSON/request field names exactly. If that is not possible for whatever reason (e.g., the field name is ambiguous because multiple different parent structs contain a field with the exact same name), pass the parent struct instead. Never invent combined variable names unless they match the exact field string.
14. If a type is not fully known based on the provided attachment (e.g., empty JSON objects like `{}` or arrays `[]` where the inner type is ambiguous), omit the field from the structs entirely.
15. If a request's complete URL is provided within the JSON response of a previous request, the function must accept that complete URL as a single string argument and process it using `url.Parse` (adhering strictly to rule 13 for naming the argument, or passing its parent struct if the field name is ambiguous). Do not hardcode the base URL or attempt to reconstruct it by extracting and passing individual query parameters.
16. If a struct passed as a function argument has two or more fields, it must be passed as a pointer rather than by value.

~~~
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done

kanopy
