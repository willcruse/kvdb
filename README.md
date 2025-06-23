# KVDB

A simple in memory key-value DB written in Go.

## Server
Run the server with `go run server/main.go`

This will run a server on port 1337 with a log file at kv.db

### CLI
Port number and log file path can be customised via command line arguments
```
  --help: display this message and exit
  --port <INT> Port to run the server on
  --log-file-path <FILE_PATH> Filepath to store the write log to
```


## Client
Run the client with `go run client/main.go`

This will run through a series of operations against a server on `localhost:1337`



## TODO
- [ ] Test harness
- [ ] CI pipeline
- [x] Document protocol
- [ ] Persistence

