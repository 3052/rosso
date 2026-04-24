# ai

1. Use the requested package name. With Go 1.22+, you MUST take the address of the range value directly instead of indexing the slice when returning a pointer from a loop.
2. Generate EXACTLY ONE file PER HTTP REQUEST. If the HAR contains multiple requests, output a separate Go file for each. The generated code MUST be inside markdown code blocks. For every file, print the marker on the very first line exactly like this:
// FILE: path/to/filename.go
3. NEVER use the standard library HTTP package. You MUST use the specified custom library for HTTP requests. Explicitly qualify calls to the custom HTTP library with its package name.
4. Unmarshal JSON responses into domain-specific struct pointers. Use the standard library JSON decoder directly on the response body stream. NEVER read the entire response body into a byte slice prior to unmarshaling. If reading into memory is necessary, reuse the existing body variable; NEVER declare a new response body variable.
5. Use URL struct literals for static URLs. NEVER use parsing functions on compile-time known URLs. For dynamic URLs, NEVER combine parsing with path escaping. NEVER construct raw queries via string concatenation; ALWAYS use the standard library values encoding method to safely generate parameters. Instantiate the values map as a separate variable on a preceding line rather than inline.
6. NEVER explicitly add standard or auto-generated headers (e.g., User-Agent, Content-Length, Accept-Encoding). ONLY set header keys for non-standard values. If no custom headers are required, pass a nil or empty value to the request function.
7. NEVER parameterize static, structural, dummy, enum-like values, or device IDs in queries, headers, or JSON request bodies. Hardcode these constants directly into the request construction instead of exposing them as arguments.
8. NEVER use anonymous structs. Either define an explicit named type or use a map.
9. When constructing JSON payloads, NEVER mix structs and maps. Choose ONE approach: use entirely a fully defined hierarchy of named structs, OR use maps entirely. NEVER embed a struct inside a map.
10. NEVER use double capitals (consecutive uppercase letters) in identifiers, including acronyms (e.g., use Id, not ID; Url, not URL). For struct fields: match the tag exactly if possible, but you MUST uppercase the first letter to export it, sanitize invalid identifiers, and lowercase consecutive capitals to comply with the double-capital rule.
11. If a type is not fully known based on the provided attachment, OMIT the field from the structs entirely.
12. NEVER alias standard library imports.
13. Identifier naming rules are strictly categorized. NEVER apply rules meant for one type to another:
    * Variables/Parameters/Loop Variables: Use simple, idiomatic Go names. When naming variables or function arguments that share an entity root with their type, you MUST use one of the following exact name/type structures: `alfaBravo Alfa`, `alfaBravo Bravo`, `alfa AlfaBravo`, or `alfa BravoAlfa`. NEVER invent arbitrary disjointed words. When declaring a variable for a parsed struct, NEVER use generic names (e.g., "data"). Derive the variable name from the entity and append a standard suffix, OR suffix the struct type so the variable can use the bare entity name. NEVER abbreviate.
    * Functions: MUST begin with a verb followed by the descriptive entity or operation name. NEVER invent action verbs. NEVER use overly brief bare verbs. NEVER abbreviate.
    * Types (Structs): The root response struct MUST closely match the entity name in the related function. If this causes a collision with a nested struct field, rename both the function and root struct to align on a new concept, or append a standard suffix to the root struct. NEVER abbreviate. NEVER append generic suffixes unless resolving a collision or pairing with a bare variable name.
    * Struct Fields: Exempt from general word-choice rules. Names MUST match original JSON keys exactly when possible. Custom types MUST match the field name if possible. Exception: For slices/collections, the custom element type MUST use the singular form of the specific logical entity it represents, NOT a generic term derived from the JSON key.
14. ONLY use pointers for struct fields, slice elements, or map values if there is a specific reason to do so. Default to using value types for nested structures.
15. Unwrapped Widevine responses MUST ALWAYS be returned as a byte slice, NEVER as a string.
16. If input comes from the user, use standard built-in types. If input comes from a previous response, you MUST pass the parent response struct directly or define a new type. When passing structs as function arguments, use a pointer if the struct has two or more fields; otherwise, pass it by value.
17. When naming the variable for a URL struct literal, use EXACTLY ONE word. Use two words IF AND ONLY IF one word is genuinely ambiguous. NEVER apply this rule to anything else unless it is the exact situation.
18. A base64 encoding flag in a HAR response means the capturing tool base64-encoded raw binary data to store it in JSON. The actual HTTP response body over the wire is raw binary bytes. NEVER implement base64 decoding for the response body in the generated code.
19. ALWAYS align variable and parameter names with standard library conventions. You MUST use standard idiomatic short names for common variables (e.g., `resp` for HTTP responses, `req` for requests, `err` for errors). When serializing a payload to pass as the body parameter of a request function, name the resulting byte slice variable identically to the function signature parameter. If constructing a struct before serialization, name the struct variable something else so the serialized byte slice can utilize the parameter identifier. When declaring a variable, parameter, or loop variable whose type shares its base entity, you MUST adhere strictly to the naming structures defined in Rule 13 to avoid case-insensitive collisions. NEVER use generic standalone variable names for root entities. NEVER carry over secondary descriptive words from the type name into the variable or parameter name unless dictated by the aforementioned forms, and NEVER use stuttering or repetitively suffixed names.

~~~go
package maya // import "41.neocities.org/maya"
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done

kanopy
