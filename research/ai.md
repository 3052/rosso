# ai

1. Use the requested package name. You may use Go 1.26 features, but do not include version build tags in the files.
2. One file per request.
3. Do not use the standard library HTTP package. Use the specified custom library for HTTP requests. You must explicitly qualify calls to the custom HTTP library with its package name; do not assume the generated code resides in the same package.
4. Unmarshal JSON responses into domain-specific struct pointers.
5. Use URL struct literals for static URLs. Do not use parsing functions on a URL that is known at compile time. For dynamic URLs, do not combine parsing with path escaping. Never construct raw queries via string concatenation; always use the standard library's values encoding method to generate query parameters safely. When assigning to the `RawQuery` field, instantiate the standard library's values map as a separate variable on a preceding line rather than nesting it inline.
6. Never explicitly add standard or automatically generated headers. Only set header keys if their values are non-standard. If no custom headers are required, pass `nil` to the request function instead of an initialized empty map.
7. Do not parameterize static, structural, dummy, or enum-like values in query parameters, headers, or JSON request bodies. Hardcode these constants directly into the request construction instead of exposing them as function arguments.
8. Never use anonymous structs. Either define an explicit named type or use a map.
9. When constructing JSON payloads, do not mix structs and maps. Choose one approach or the other: either use a fully defined hierarchy of named structs, or use maps entirely. Do not embed a struct inside a map.
10. Do not use any double capitals (consecutive uppercase letters) in identifiers, including acronyms. For struct fields: match the tag exactly if possible, but you must uppercase the first letter to export it, sanitize it if the tag is not a valid identifier, and lowercase consecutive capital letters to comply with the double capital rule.
11. If a type is not fully known based on the provided attachment, omit the field from the structs entirely.
12. Do not alias standard library imports.
13. Identifier naming rules are strictly separated by category. Do not apply rules meant for one type of identifier to another:
    * Variables and Parameters: Use simple, direct, idiomatic Go names. If a variable or parameter name would identically match its type name (ignoring case and pointer prefixes), you must append a full-word suffix (such as `Data`) to the variable or parameter to prevent repetition (e.g., `entityData *Entity`). If the type name contains secondary descriptive words or suffixes (such as `Response` or `Session`), use only the single primary base word for the variable or parameter name (e.g., `entity *EntityResponse`, `entity *EntitySession`). Must not use abbreviations (e.g., never use `entityResp`).
    * Functions: Must begin with a verb followed by the descriptive name of the entity or operation; do not invent alternative action verbs. Do not use overly brief function names consisting only of a bare verb. Must not use abbreviations.
    * Types (Structs): The root response struct type must closely match the entity name used in the related function name. If this causes a collision with a nested struct field, either rename both the function and the root struct to align on a new concept, or append a standard suffix such as `Response` to the root struct type. Must not use abbreviations. Do not append generic suffixes unless resolving a collision.
    * Struct Fields: Exempt from general word-choice rules. Struct field names must match the original JSON keys exactly when possible. When a struct field uses a custom type, the custom type name must match the field name if possible. Exception: If the field is a slice or collection, the custom type representing a single element must use the singular form of the specific logical entity it represents, and must not be a generic term derived from the JSON key.
14. Only use pointers for struct fields, slice elements, or map values if there is a specific reason to do so. Default to using value types for nested structures.
15. Unwrapped Widevine responses must always be returned as a byte slice, never as a string.
16. If input comes from the user, use standard built-in types. If input comes from a previous response, you must pass the parent response struct directly or define a new type for the field. When passing structs as function arguments, use a pointer if the struct has two or more fields; otherwise, pass it by value.
17. When naming the variable for a URL struct literal, use a single word. Use two words if and only if one word is ambiguous. Do not apply this rule to anything else unless it is the exact situation.
18. When a HAR file's response content includes an encoding flag indicating base64, this indicates the capturing tool base64-encoded raw binary data to store it in JSON. The actual HTTP response body over the wire is raw binary bytes. Do not implement base64 decoding for the response body in the generated code.
19. Always align variable and parameter names with standard library conventions and function signatures. Use `resp` instead of `res` for HTTP responses. When serializing a payload to pass as the body parameter of a request function, name the resulting byte slice variable `body` instead of `data` or `buf` to match the function signature. If constructing a struct before serialization, name the struct variable something else so the serialized byte slice can utilize the `body` identifier. When declaring a variable or parameter whose type is exactly the single-word base entity (ignoring pointers), you must append `Data` to the name (e.g., `entityData *Entity`) to avoid a case-insensitive match with the type. Do not carry over secondary descriptive words from the type name into the variable or parameter name (e.g., do not use `entitySessionData *EntitySession`), and do not use stuttering or repetitively suffixed names.

~~~go
package maya // import "41.neocities.org/maya"
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done
