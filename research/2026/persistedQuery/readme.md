# persistedQuery

with GraphQL, if the client is using `persistedQuery`, you can either modify
the request `sha256Hash`, or modify the response message to
`PersistedQueryNotFound` - if you do either, the client will download a
JavaScript that includes the full queries instead of just the hashes, then
client will try the request again with the full query which you can then edit
or whatever

https://crawlee.dev/blog/graphql-persisted-query
