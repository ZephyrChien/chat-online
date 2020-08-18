# A simple CLI tool for online chatting
## Build
>>to build server:
```
go build -o server -tags=server
```
>>to build client:
```
go build -o client
```
## Usage
```
server -s <ip_address> -p <port>
```
```
client -n <name> -h <remote_address:port> -l <local_address:port>
client -l <local_address:port> --write
```
## TODO:
1. add more cmd for client and server(e.g. /help)
2. encrypt data during transport
3. beautify client

