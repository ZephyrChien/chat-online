# A simple CLI tool for online chatting
## Build
>to build server:
```
go build -o server -tags=server
```
>to build client:
```
go build -o client
```
## Usage
>Start the server with:
```
server -s <ip_address> -p <port> -k <16byte_key> --ssl --cert <crt> --key <key>
```
>Or you can simply run it with:
```
server -p <port> //default address is 0.0.0.0:8000
```
>Tips: the default key is 00...0 (x16)
You'd better set another one by the [-k] option.  
  
>As for the client, it would be a little bit troublesome. 
>In order not to mix up Text and Input(Stdin), 
>You need to run two clients at the same time. 
>One to send & receive message, the other used as keyboard.  
>This might be solved later(not sure, I need to rewrite it one day=_=  
  

>Example:
```
client -n <name> -h <remote_address:port> -l <local_address:port> //receive &send message
client -l <local_address:port> --write //input area
```
>The two clients are locally connected through the LoopBack Address. 
>You can alse use other tools as a writer, such as netcat:
```
nc -nv <local_address> <port>
```
## TODO:
1. more cmds for client and server(e.g. /kick)
2. more elements attached to a single message(such as level, privilege..)
3. beautify client(may be one with GUI)
4. enhance log

