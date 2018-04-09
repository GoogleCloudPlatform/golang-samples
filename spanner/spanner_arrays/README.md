This is a demonstration of the use of Cloud Spanner arrays using the Go
client library.  

The key takeaway here is that even if your string fields are not nullable, if
they are returned as an arrray, you still _must_ use the `spanner.NullString`
type to _receive_ them.  _Populating_ them, however, can be done with either
`spanner.NullString` _or_ plain `string`.

