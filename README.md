# Links gRPC Service and Client

This is a gRPC service for managing saved links for personal use and to learn gRPC.

## Application

├── client
│   └── main.go
├── links
│   ├── links.pb.go
│   └── links.proto
├── mock_links
│   ├── server_mock.go
│   └── server_mock_test.go
└── server
    ├── main.go
    ├── main_test.go
    ├── model.go
    ├── saved.sqlite
    └── saved_test.sqlite

### client

Initial testing of the server was done using this but also to test implementing a client.

### links

The proto file and generated go code from running `compile.sh` which just runs the `protoc` command.

Install proto3 from source

```bash
git clone https://github.com/google/protobuf
./autogen.sh ; ./configure ; make ; make install
```

Update protoc Go bindings via

```bash
go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
```
### server

This is the service that is the main application and what you would actually put in a Docker container. The tests here setup a `Server` in order to be able to test streams.

### mock_links

**Under Development** is an attempt to utilize `mochgen` for client testing.

---

## Development

Be sure to add to and run tests for development.  This will create a test database `saved_test.sqlite` if one does not exist.  The file is not deleted after tests run so if changes are made to an existing table's structure, the database file should be removed before running the tests.

Once the tests all pass you can build your own docker image with the included docker file like so:

```shell
docker build -t contd/links-rpc .
```
If you use `go install` the binary will expect the saved.sqlite file to be in the same directory.  You can override this by passing an environment variable like so:

```shell
SQLITE_PATH=/some/other/path/saved.sqlite links
```

This assumes your `PATH` includes `$GOPATH/bin` and you must include the file name of the sqlite database file you want to use.  The tables are not created by the application and it expects them to be there so you can use the one created from running the tests and just rename it.

## Docker

To run this in docker and persist the sqlite database, use the following once you've created an image from the `Dockerfile`:

```shell
docker run --name golinksrpc -d -p 5555:5555 -v $GOPATH/src/githib.com/contd/links:/data contd/links-rpc
```
