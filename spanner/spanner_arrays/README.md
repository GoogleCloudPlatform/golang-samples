This is a demonstration of the use of Cloud Spanner arrays using the golang
client library.  The key takeaway here is that even if your fields are not
nullable, you still _must_ use the `spanner.NullString` type to receive them.
Populating them can be done with either `spanner.NullString` _or_ plain
`string`.

