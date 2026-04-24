# ai

1. Use the requested package name. You may use Go 1.26 features, but NEVER include version build tags in the files. Because Go 1.22+ fixes loop variable scoping, you MUST take the address of the range value variable directly instead of indexing the slice when returning a pointer from a loop.
2. Generate EXACTLY ONE file PER HTTP REQUEST found in the source material. If the provided HAR contains multiple HTTP requests, you MUST output a separate Go file for each individual HTTP request. For EVERY file, use the delimiters:
```
// --- START OF FILE path/to/filename.go ---
// [file contents go here]
// --- END OF FILE path/to/filename.go ---
```
3. NEVER use the standard library HTTP package for making requests. You MUST use the specified custom library for HTTP requests. You MUST explicitly qualify calls to the custom HTTP library with its package name; do not assume the generated code resides in the same package.
4. Unmarshal JSON responses into domain-specific struct pointers. When parsing a JSON response body, use the standard library JSON decoder directly on the response body stream. NEVER read the entire response body into a byte slice prior to unmarshaling. If you must read the body into memory, reuse the existing body variable; NEVER declare a new response body variable.
5. Use URL struct literals for static URLs. NEVER use parsing functions on a URL that is known at compile time. For dynamic URLs, NEVER combine parsing with path escaping. NEVER construct raw queries via string concatenation; ALWAYS use the standard library's values encoding method to generate query parameters safely. When assigning to the raw query field, instantiate the standard library's values map as a separate variable on a preceding line rather than nesting it inline.
6. NEVER explicitly add standard or automatically generated headers (e.g., User-Agent, Content-Length, Accept-Encoding). ONLY set header keys if their values are non-standard. If no custom headers are required, pass a nil or empty value to the request function instead of an initialized empty map.
7. NEVER parameterize static, structural, dummy, or enum-like values in query parameters, headers, or JSON request bodies. Hardcode these constants directly into the request construction instead of exposing them as function arguments.
8. NEVER use anonymous structs. Either define an explicit named type or use a map.
9. When constructing JSON payloads, NEVER mix structs and maps. Choose ONE approach: either use a fully defined hierarchy of named structs, OR use maps entirely. NEVER embed a struct inside a map.
10. NEVER use double capitals (consecutive uppercase letters) in identifiers, including acronyms (e.g., use `Id`, not `ID`; `Url`, not `URL`). For struct fields: match the tag exactly if possible, but you MUST uppercase the first letter to export it, sanitize it if the tag is not a valid identifier, and lowercase consecutive capital letters to comply with the double-capital rule.
11. If a type is not fully known based on the provided attachment, OMIT the field from the structs entirely.
12. NEVER alias standard library imports.
13. Identifier naming rules are strictly separated by category. NEVER apply rules meant for one type of identifier to another:
    * Variables, Parameters, and Loop Variables: Use simple, direct, idiomatic Go names. If a variable, parameter, or loop variable name would identically match its type name (ignoring case and pointer prefixes), you MUST append a full-word suffix to the name to prevent repetition. If the type name contains secondary descriptive words or suffixes, use only the single primary base word for the variable or parameter name. NEVER use abbreviations.
    * Functions: MUST begin with a verb followed by the descriptive name of the entity or operation; NEVER invent alternative action verbs. NEVER use overly brief function names consisting only of a bare verb. NEVER use abbreviations.
    * Types (Structs): The root response struct type MUST closely match the entity name used in the related function name. If this causes a collision with a nested struct field, either rename both the function and the root struct to align on a new concept, or append a standard suffix to the root struct type. NEVER use abbreviations. NEVER append generic suffixes unless resolving a collision.
    * Struct Fields: Exempt from general word-choice rules. Struct field names MUST match the original JSON keys exactly when possible. When a struct field uses a custom type, the custom type name MUST match the field name if possible. Exception: If the field is a slice or collection, the custom type representing a single element MUST use the singular form of the specific logical entity it represents, and MUST NOT be a generic term derived from the JSON key.
14. ONLY use pointers for struct fields, slice elements, or map values if there is a specific reason to do so. Default to using value types for nested structures.
15. Unwrapped Widevine responses MUST ALWAYS be returned as a byte slice, NEVER as a string.
16. If input comes from the user, use standard built-in types. If input comes from a previous response, you MUST pass the parent response struct directly or define a new type for the field. When passing structs as function arguments, use a pointer if the struct has two or more fields; otherwise, pass it by value.
17. When naming the variable for a URL struct literal, use EXACTLY ONE word. Use two words IF AND ONLY IF one word is genuinely ambiguous. NEVER apply this rule to anything else unless it is the exact situation.
18. When a HAR file's response content includes an encoding flag indicating base64, this indicates the capturing tool base64-encoded raw binary data to store it in JSON. The actual HTTP response body over the wire is raw binary bytes. NEVER implement base64 decoding for the response body in the generated code.
19. ALWAYS align variable and parameter names with standard library conventions and function signatures. Use standard short names for HTTP responses. When serializing a payload to pass as the body parameter of a request function, name the resulting byte slice variable identically to the function signature parameter. If constructing a struct before serialization, name the struct variable something else so the serialized byte slice can utilize the parameter identifier. When declaring a variable, parameter, or loop variable whose type is exactly the single-word base entity (ignoring pointers), you MUST append a suffix to the name to avoid a case-insensitive match with the type. NEVER carry over secondary descriptive words from the type name into the variable or parameter name, and NEVER use stuttering or repetitively suffixed names.

~~~go
package maya // import "41.neocities.org/maya"
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done

kanopy
