#### To generate singer.pb.go and descriptors.pb file from singer.proto using `protoc`
```shell
cd spanner_snippets/spanner/testdata
protoc --proto_path=./protos/ --include_imports --descriptor_set_out=./protos/descriptors.pb --go_out=./protos/ protos/singer.proto
```
