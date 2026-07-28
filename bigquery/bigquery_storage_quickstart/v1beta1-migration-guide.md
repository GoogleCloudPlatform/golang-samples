# Migrating BigQuery Storage API from v1beta1 to v1: Go

This guide shows how to migrate Go code using the BigQuery Storage API from
version `v1beta1` to `v1`.

## Key Changes

*   **Package Imports**:
    *   Client library: `cloud.google.com/go/bigquery/storage/apiv1beta1` ->
        `cloud.google.com/go/bigquery/storage/apiv1`
    *   Proto types:
        `google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1beta1`
        -> `cloud.google.com/go/bigquery/storage/apiv1/storagepb`
*   **Service Client**: `NewBigQueryStorageClient` is replaced by
    `NewBigQueryReadClient`.
*   **Table Reference**: `TableReference` struct is replaced by a simple string
    representation of the table path in `ReadSession.Table`.
*   **Session Configuration**: Configuration fields (table, format, read
    options) have moved into `ReadSession` struct, which is passed in
    `CreateReadSessionRequest`.
*   **Parallelism**: `RequestedStreams` is replaced by `MaxStreamCount`.
*   **Sharding Strategy**: `ShardingStrategy` field is removed. The server now
    automatically balances the streams.
*   **Read Rows Request**: `ReadPosition` is flattened. You now pass the stream
    name directly as `ReadStream` and the `Offset` as a top-level field in
    `ReadRowsRequest`.
*   **Stream Type**: `Stream` type is renamed to `ReadStream`.

## Code Comparison

### 1. Imports and Client Initialization

**v1beta1:**

```go
import (
    "context"

    bqStorage "cloud.google.com/go/bigquery/storage/apiv1beta1"
    bqStoragepb "google.golang.org/genproto/googleapis/cloud/bigquery/storage/v1beta1"
)

ctx := context.Background()
client, err := bqStorage.NewBigQueryStorageClient(ctx)
if err != nil {
    // handle error
}
defer client.Close()
```

**v1:**

```go
import (
    "context"

    bqStorage "cloud.google.com/go/bigquery/storage/apiv1"
    bqStoragepb "cloud.google.com/go/bigquery/storage/apiv1/storagepb"
)

ctx := context.Background()
client, err := bqStorage.NewBigQueryReadClient(ctx)
if err != nil {
    // handle error
}
defer client.Close()
```

### 2. Creating a Read Session

**v1beta1:**

```go
tableRef := &bqStoragepb.TableReference{
    ProjectId: "bigquery-public-data",
    DatasetId: "usa_names",
    TableId:   "usa_1910_current",
}

readOptions := &bqStoragepb.ReadSession_TableReadOptions{
    SelectedFields: []string{"name"},
    RowRestriction: "state = 'WA'",
}

req := &bqStoragepb.CreateReadSessionRequest{
    Parent:         "projects/read-session-project",
    TableReference: tableRef,
    ReadOptions:    readOptions,
    RequestedStreams: 1,
    Format:         bqStoragepb.DataFormat_AVRO,
    ShardingStrategy: bqStoragepb.ShardingStrategy_LIQUID,
}

session, err := client.CreateReadSession(ctx, req)
```

**v1:**

```go
// Table path is now a string: projects/{project}/datasets/{dataset}/tables/{table}
tablePath := "projects/bigquery-public-data/datasets/usa_names/tables/usa_1910_current"

readOptions := &bqStoragepb.ReadSession_TableReadOptions{
    SelectedFields: []string{"name"},
    RowRestriction: "state = 'WA'",
}

// ReadSession holds the session configuration
readSession := &bqStoragepb.ReadSession{
    Table:      tablePath,
    DataFormat: bqStoragepb.DataFormat_AVRO, // Format renamed to DataFormat
    ReadOptions: readOptions,
}

req := &bqStoragepb.CreateReadSessionRequest{
    Parent:      "projects/read-session-project",
    ReadSession: readSession,
    MaxStreamCount: 1, // RequestedStreams renamed to MaxStreamCount
}

session, err := client.CreateReadSession(ctx, req)
```

### 3. Reading Rows

**v1beta1:**

```go
// Use the first stream
stream := session.GetStreams()[0]

req := &bqStoragepb.ReadRowsRequest{
    ReadPosition: &bqStoragepb.StreamPosition{
        Stream: stream, // Stream object
        Offset: 0,
    },
}

rowStream, err := client.ReadRows(ctx, req)
// Iterate over rowStream.Recv()
```

**v1:**

```go
// Use the first stream
stream := session.GetStreams()[0] // returns *bqStoragepb.ReadStream

req := &bqStoragepb.ReadRowsRequest{
    ReadStream: stream.GetName(), // Pass stream name string directly
    Offset:     0,                // Offset is top-level
}

rowStream, err := client.ReadRows(ctx, req)
// Iterate over rowStream.Recv()
```
