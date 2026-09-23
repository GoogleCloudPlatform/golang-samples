# PR Title

```text
feat(bigtable): migrate table admin samples to cloud.google.com/go/bigtable/admin/apiv2
```

---

# PR Description

## Description

Migrates all Cloud Bigtable table administration samples and integration tests in `bigtable/` from the legacy handwritten `cloud.google.com/go/bigtable.AdminClient` to the modernized GAPIC-generated `cloud.google.com/go/bigtable/admin/apiv2.BigtableTableAdminClient` (and protobuf request/schema types in `cloud.google.com/go/bigtable/admin/apiv2/adminpb`).

### Key API Surface Changes
- **Client Initialization**: Replaced project/instance-bound `bigtable.NewAdminClient(ctx, project, instance)` (`*bigtable.AdminClient`) with stateless `admin.NewBigtableTableAdminClient(ctx)` (`*admin.BigtableTableAdminClient`).
- **Request-Scoped Resource Names**: Updated table admin calls (`CreateTable`, `DeleteTable`, `GetTable`, `ListTables`, `ModifyColumnFamilies`, `DropRowRange`) to pass structured `adminpb` request protos containing full resource paths (`projects/{project}/instances/{instance}` and `projects/{project}/instances/{instance}/tables/{table}`).
- **Garbage Collection Rules (`adminpb.GcRule`)**: Replaced legacy `bigtable.GCPolicy` helpers (`MaxAgePolicy`, `MaxVersionsPolicy`, `UnionPolicy`, `IntersectionPolicy`) and `SetGCPolicy` calls with `adminpb.GcRule` oneof configurations (`GcRule_MaxAge`, `GcRule_MaxNumVersions`, `GcRule_Union_`, `GcRule_Intersection_`) applied via `ModifyColumnFamilies`.
- **Aggregate Column Families**: Updated `view_count` column family setup in `writes/writes_test.go` to configure `adminpb.Type_AggregateType` (`Int64Type` with `BigEndianBytes` encoding and `Sum` aggregator) directly on `adminpb.ColumnFamily`.
- **Dependencies**: Updated `cloud.google.com/go/bigtable` to `v1.57.0` in `bigtable/go.mod`.

### Summary of Modified Files
- **`deletes/`**:
  - `deletes/drop_row_range.go`: Migrated `dropRowRange` to `admin.NewBigtableTableAdminClient` and `adminpb.DropRowRangeRequest`.
  - `deletes/deletes_test.go`: Migrated table and column family test setup/teardown to `admin.BigtableTableAdminClient`.
- **`garbagecollection/`**:
  - `garbagecollection/create_family_gc_max_age.go`: Migrated to `admin.BigtableTableAdminClient` and `adminpb.GcRule_MaxAge`.
  - `garbagecollection/create_family_gc_max_versions.go`: Migrated to `admin.BigtableTableAdminClient` and `adminpb.GcRule_MaxNumVersions`.
  - `garbagecollection/create_family_gc_union.go`: Migrated to `admin.BigtableTableAdminClient` and `adminpb.GcRule_Union_`.
  - `garbagecollection/create_family_gc_intersect.go`: Migrated to `admin.BigtableTableAdminClient` and `adminpb.GcRule_Intersection_`.
  - `garbagecollection/create_family_gc_nested.go`: Migrated to `admin.BigtableTableAdminClient` with nested `adminpb.GcRule_Union_` and `adminpb.GcRule_Intersection_`.
  - `garbagecollection/update_gc_rule.go`: Migrated to `admin.BigtableTableAdminClient` using `ModifyColumnFamiliesRequest_Modification_Update` and `fieldmaskpb.FieldMask`.
  - `garbagecollection/garbagecollection_test.go`: Migrated test table creation/deletion and output assertions.
- **`helloworld/`**:
  - `helloworld/main.go`: Migrated `bigtable_hw_imports`, `bigtable_hw_connect`, `bigtable_hw_create_table`, and `bigtable_hw_delete_table` regions to `admin.BigtableTableAdminClient`.
- **`search/` & `usercounter/`**:
  - `search/search.go`: Migrated `main`, `handleReset`, `copyTable`, and `handleCopy` to `*admin.BigtableTableAdminClient`.
  - `usercounter/main.go`: Migrated table and column family initialization to `admin.NewBigtableTableAdminClient`.
- **`filters/`, `reads/`, & `writes/`**:
  - `filters/filters_test.go`, `reads/reads_test.go`, `writes/writes_test.go`: Migrated test setup/teardown to `admin.NewBigtableTableAdminClient`.
  - `filters/filters.go`, `writes/write_batch.go`, `writes/write_conditionally.go`, `writes/write_increment.go`: Corrected error wrapping strings on `bigtable.NewClient` from `"bigtable.NewAdminClient: %w"` to `"bigtable.NewClient: %w"`.

## Checklist
- [x] I have followed [Contributing Guidelines from CONTRIBUTING.MD](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/CONTRIBUTING.md)
- [x] **Tests** pass: `go test -v ./...` (see [Testing](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/CONTRIBUTING.md#testing))
- [x] **Code formatted**: `gofmt` (see [Formatting](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/CONTRIBUTING.md#formatting))
- [x] **Vetting** pass: `go vet` (see [Formatting](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/CONTRIBUTING.md#formatting))
- [ ] These samples need a new **API enabled** in testing projects to pass (let us know which ones)
- [ ] These samples need a new/updated **env vars** in testing projects set to pass (let us know which ones)
- [ ] This sample adds a new sample directory, and I updated the [CODEOWNERS file](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/.github/CODEOWNERS) with the codeowners for this sample
- [ ] This sample adds a new **Product API**, and I updated the [Blunderbuss issue/PR auto-assigner](https://github.com/GoogleCloudPlatform/golang-samples/blob/main/.github/blunderbuss.yml) with the codeowners for this sample
- [ ] Please **merge** this PR for me once it is approved
