Instructions to create and work with the workspace

Create the workspace folder
```bash
mkdie workspacenae
```
cd inside the workspace folder
```bash
cd workspace
```
Initiate module as:
```bash
go work init ./modulename
```
For example:
```bash
go work init ./walistner
```
The module folder will be created as empty folder, and the `go.work` file will looks something like:
```bash
go 1.18

use ./walistner
```
cd inside the module folder
```bash
cd modulename
```
Run the mod init command
```bash
go mod init
```
Run the mod tidy command
```bash
go mod tidy
```

To run any module inside the workspace use:
```bash
go run ./modulename
```
As:
```bash
go run ./walistner
```

To create gRPC and gRPC gateway connection:

Install `vscode-proto3` extenstion at your VS code
Install `protoc` binary from https://github.com/protocolbuffers/protobuf/releases
Install the `protoc` library:
https://developers.google.com/protocol-buffers/docs/reference/go-generated
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
```
At the proto file define the Go package it will belongs to, as:
```bash
option go_package = "example.com/project/protos/fizz";
```
In our case, we have in top of the file this line:
```bash
option go_package = "/walistner";
```
Be careful in the code above, The import path must contain at least one period ('.') or forward slash ('/') character.

SO, in our case, our compilation line will become:
1. Generate the messages
2. Generate the services
3. Generate the Gateway

All of these are generated in single line as below:
```bash
protoc -I ./proto \
  --go_out ./proto --go_opt paths=source_relative \
  --go-grpc_out ./proto --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./proto --grpc-gateway_opt paths=source_relative \
  ./proto/data.proto
  ```
So, to get it done, we need to do 2 things:
1. Creat a `proto` folder inside the required module `walistner`
```bash
$ mkdir ./walistner/proto 
```
2. Generate the gRPC messages, services and cateway
```bash
$ protoc -I ./proto --go_out ./walistner/proto --go_opt paths=source_relative --go-grpc_out ./walistner/proto --go-grpc_opt paths=source_relative --grpc-gateway_out ./walistner/proto --grpc-gateway_opt paths=source_relative ./proto/data.proto
```
