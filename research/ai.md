# ai

## max 1200 words

1. Use the requested package name. With Go 1.22+, you MUST take the address of the range value directly when returning a pointer from a loop.
2. Generate EXACTLY ONE file PER HTTP REQUEST. Output a separate Go file for each request. You MUST encapsulate the entire code for each file strictly within standard Markdown code blocks (using triple backticks for go). Print the marker on the first line inside the code block: `// FILE: path/to/filename.go`.
3. NEVER use the standard library HTTP package. Explicitly qualify calls to the custom HTTP library (i.e., `import "41.neocities.org/maya"` and call `maya.Get`). DO NOT declare the generated files as package `maya`, as this prevents proper qualification.
4. Unmarshal JSON responses into domain-specific struct pointers. If the root JSON response consists of a single wrapper field, you MUST unwrap it inside the function using an anonymous struct and return the inner struct directly; do not define a named type for the wrapper. Use the standard library JSON decoder directly on the response body stream. NEVER read the entire response body into a byte slice unless handling a DRM response (see Rule 16).
5. Use URL struct literals for static URLs. NEVER use parsing functions on compile-time known URLs. For dynamic URLs, NEVER combine parsing with path escaping. NEVER construct raw queries via string concatenation; ALWAYS use the standard library `url.Values` encoding method to safely generate parameters. Instantiate the values map on a preceding line.
6. NEVER add standard or auto-generated headers. ONLY set keys for non-standard values. If no custom headers are required, pass `nil`.
7. NEVER parameterize static, structural, dummy, enum-like values, or device IDs in queries, headers, or JSON request bodies. Hardcode these constants directly into the request construction.
8. NEVER instantiate anonymous structs using struct literals. Declare an anonymous struct variable strictly for JSON unmarshaling to unwrap outer layers. Otherwise, define an explicit named type or use a map.
9. When constructing JSON payloads, NEVER mix structs and maps. Choose ONE approach: use entirely a fully defined hierarchy of named structs, OR use maps entirely. NEVER embed a struct inside a map.
10. NEVER use double capitals (consecutive uppercase letters) in identifiers, including acronyms. For struct fields: match the tag exactly if possible, but uppercase the first letter to export it, sanitize invalid identifiers, and lowercase consecutive capitals.
11. If a type is not fully known based on the provided attachment, OMIT the field from the structs entirely.
12. NEVER alias standard library imports.
13. Variable, Parameter, Receiver, and Loop Variable naming: YOU MUST NEVER USE THE EXACT SAME STRING (IGNORING CASE) FOR BOTH THE IDENTIFIER AND ITS TYPE.
    * For SINGLE-WORD Types (e.g., `Config`, `Item`): The variable name MUST be TWO words combining a package-agnostic descriptor with the type (e.g., `targetConfig Config`, `sourceItem Item`, `activeAccount Account`). NEVER use the exact same word (e.g., `config Config`).
    * For TWO-WORD Types (e.g., `ItemDetails`, `PlaybackResponse`): The variable name MUST be a SINGLE word extracted directly from the type name to prevent stuttering (e.g., `item ItemDetails`, `playback PlaybackResponse`). NEVER add arbitrary descriptors, NEVER use the exact type name, and NEVER use generic suffixes like `resp` or `payload` for custom types.
    * Standard Library Exceptions: ONLY use standard library conventions for their specific types (e.g., `resp *http.Response`, `req *http.Request`).
14. Function and Type naming: NEVER apply rules meant for variables to functions or types.
    * Functions: MUST begin with a verb followed by the entity or operation name. NEVER abbreviate.
    * Root Structs: MUST closely match the entity name in the related function. If this causes a collision with a nested struct field, rename both the function and root struct to align on a new concept, or append a standard suffix to the root struct. NEVER abbreviate.
    * Nested Structs & Fields: Names MUST match original JSON keys exactly when possible. The custom type name for a nested struct MUST exactly match its field name. You MUST NOT blindly prefix the parent struct's name or add arbitrary qualifiers unless resolving a package-level naming collision. For slices/collections, the custom element type MUST use the singular form of the specific logical entity.
15. ONLY use pointers for struct fields, slice elements, or map values if there is a specific reason. Default to using value types for nested structures.
16. DRM licensing responses MUST ALWAYS be read completely and returned directly as a byte slice, NEVER as a string, and NEVER unmarshaled into XML or JSON structs.
17. STRICT PARAMETER TYPE RULE: You must meticulously trace the origin of every function parameter.
    * If a parameter represents the INITIAL raw user input required to start a chain of requests, you MUST use a standard built-in primitive type (e.g., `string`, `int`, `bool`).
    * If an input parameter represents data parsed, extracted, or returned from a PREVIOUS HTTP response, you MUST NEVER use a primitive type. Instead, either (A) Pass the entire parent response struct from the preceding request (by pointer if 2+ fields, otherwise by value), OR (B) Define a distinct custom named type for that specific data, assign it to the field in the parsing struct, and require that exact custom type in the subsequent function.
18. When naming the variable for a URL struct literal, use EXACTLY ONE word.
19. A base64 encoding flag in a HAR HAR response means the capturing tool base64-encoded raw binary data to store it in JSON. The actual HTTP response body over the wire is raw binary bytes. NEVER implement base64 decoding for the response body in the generated code.
20. ALWAYS align variable and parameter names with standard library conventions, EXCEPT YOU MUST NEVER USE SINGLE-LETTER VARIABLES. Use descriptive one-word (or two-word if ambiguous) names. ALWAYS reserve the name `resp` STRICTLY for `*http.Response` variables; NEVER use `resp` for custom struct types. ALWAYS use `query` or `values` for URL query maps. When serializing a payload, name the resulting byte slice variable identically to the function signature parameter. If constructing a struct before serialization, name the struct variable `payload`. When declaring a variable, parameter, receiver, or loop variable whose type shares its base entity, you MUST adhere strictly to the naming structures defined in Rule 13 to avoid case-insensitive collisions. NEVER carry over secondary descriptive words from the type name into the variable or parameter name unless dictated by the aforementioned forms, and NEVER use stuttering or repetitively suffixed names.
21. NEVER use pointers for string parameters to represent optionality. Pass the string by value and use the zero-value (empty string `""`) to check for absence.
22. NEVER deduplicate, skip, or merge requests based solely on matching URL paths. You MUST analyze the entire trace to identify state changes. If the same endpoint is called multiple times with differing headers (e.g., unauthenticated vs. authenticated), varying query parameters, or different payloads, it represents distinct operational states. You MUST account for this by either: (A) Generating a single parameterized function that conditionally applies the varying headers/queries based on the presence of the input (e.g., if a `SessionKey` string parameter is empty, omit the authorization header), or (B) Generating separate functions and files for each distinct state exactly as they appear in the flow (e.g., `GetInitialConfig` vs `GetAuthenticatedConfig`).

## maya

~~~go
package maya // import "41.neocities.org/maya"
func Get(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Head(targetUrl *url.URL, headers map[string]string) (*http.Response, error)
func Post(targetUrl *url.URL, headers map[string]string, body []byte) (*http.Response, error)
~~~

## done

1. kanopy
2. tubi
3. plex
4. rakuten
5. roku
