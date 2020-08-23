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
Start the server with
```
server -s <ip_address> -p <port> -key <16byte_key> --ssl <cert> <key> [optional]
```
Or you can simply run it with
```
server -p <port> //default address is 0.0.0.0:8000
```
Tips: the default key is 000...0 (x16)
You'd better set another one by [-key] option

As for the client, it would be a little bit troublesome
In order not to mix up Text and Input(Stdin)
You need to run two clients at the same time
One to send & receive message, the other one used as keyboard
This problem might be solved later(not sure, I need to rewrite it one day=_=
Example:
```
client -n <name> -h <remote_address:port> -l <local_address:port> //receive &send message
client -l <local_address:port> --write //input area
//The two clients are locally connected with the LoopBack Address
//You can alse use other tools as a writer, such as netcat:
nc -nv <local_address> <port>
```
## TODO:
1. more cmds for client and server(e.g. /kick)
2. more elements attached to a single client/message(such as level,privilege..)
3. beautify client(may be one with GUI)
4. enhance log

