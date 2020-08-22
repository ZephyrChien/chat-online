package cmd

//Client :specified by name
type Client struct {
	IP     string
	Name   string
	Extra  string
	Tunnel chan<- string
}

//Status :the combination of channel
type Status struct {
	Entering chan Client
	Leaving  chan Client
	Message  chan string
}

//Command :empty unless the message starts with "/"
type Command struct {
	Is  bool   `json:"apply"`
	Key string `json:"option"`
	Val string `json:"value"`
}

//Data :client send data in this format
type Data struct {
	Name    string  `json:"name"`
	Message string  `json:"message"`
	CMD     Command `json:"cmd"`
	Extra   string  `json:"extra"`
}
