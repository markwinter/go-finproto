package soupbintcp

type ClientOption func(client *Client)

// WithAuth sets the username and password to use when connecting to the Server
func WithAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.username = username
		c.password = password
	}
}

// WithCallback sets the callback function for every sequenced packet received
func WithCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.packetCallback = callback
	}
}

// WithUnsequencedCallback sets the callback function for every unsequenced packet received
func WithUnsequencedCallback(callback func([]byte)) ClientOption {
	return func(c *Client) {
		c.unsequencedCallback = callback
	}
}

// WithDebugCallback sets the callback function for every debug packet received. Not normally used
func WithDebugCallback(callback func(string)) ClientOption {
	return func(c *Client) {
		c.debugCallBack = callback
	}
}

// WithServer sets the ip and port to use to connect to the Server
func WithServer(ip, port string) ClientOption {
	return func(c *Client) {
		c.serverIp = ip
		c.serverPort = port
	}
}

// WithSession sets the initial session and sequence number when connecting to the Server
func WithSession(id string, sequence uint64) ClientOption {
	return func(c *Client) {
		c.session = id
		c.sequenceNumber = sequence
	}
}
