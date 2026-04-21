# ai

1. Use the requested package name and Go 1.26.
2. One file per request.
3. Do not use `net/http`.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use `41.neocities.org/maya` for HTTP requests.
6. Use `url.URL` struct literals for static URLs. Do not use `url.Parse` on a URL that is known at compile time. For dynamic URLs, do not combine `url.Parse` with `url.PathEscape` (use one or the other). Never construct `RawQuery` via string concatenation; always use `url.Values` and its `Encode()` method to generate query parameters safely.
7. Do not use single-letter variables. You must use the exact name when possible (such as matching original JSON field names). Use a single word for variables unless a second word is strictly needed for clarity or disambiguation. If you need to add a second word to a variable or function name for descriptiveness or clarity, do it. Do not force strictly one-word names if it results in altering conventional names or creating poorly named functions. Never append types, categories, or unnecessary descriptors to create a two-word name if a shorter name conveys the meaning clearly. Do not artificially lengthen variable names by appending redundant nouns to an already sufficiently descriptive base word.
8. Never explicitly add standard or automatically generated headers like `accept-encoding` or default `user-agent` strings. Only set the `user-agent` key if its value is non-standard. Do not parameterize headers that contain static or non-standard values; hardcode them directly in the request headers instead of passing them as function arguments.
9. If a struct has two or more fields, it must be passed as a pointer when used as a function parameter.
10. If you define a struct for a JSON request payload containing novel user data, you must use that struct directly as the function's input parameter. Do not decompose it into individual primitive arguments. Alternatively, avoid defining the struct entirely and use a map internally. This rule does not apply if the input is coming from a previous request, in which case you must use the previously obtained structs, nested structs, or individual fields as arguments and construct the payload internally regardless.
11. Do not parameterize static, structural, or enum-like fields in JSON request bodies. Hardcode these constants directly into the payload generation.
12. If multiple required values naturally originate from the same previously defined struct (such as a response struct), pass that struct directly as the argument instead of extracting its individual fields. If only a single field from a struct is needed, prefer to pass that specific field directly using its exact name. However, if the exact field name is ambiguous (meaning two or more structs in the package share the same field name), you must pass the entire struct instead. Do not prefix the parameter name.
13. Never use anonymous structs. Either define an explicit named type or use a map.
14. When constructing JSON payloads, do not mix structs and maps. Choose one approach or the other: either use a fully defined hierarchy of named structs, or use maps entirely. Do not embed a struct inside a map.
15. Do not use any double capitals (consecutive uppercase letters) in identifiers, including acronyms.

~~~
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done

1. kanopy
2. tubi
