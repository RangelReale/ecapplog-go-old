package ecapplog

type Option func(*Client)

func WithAddress(address string) Option {
	return func (c *Client) {
		c.address = address
	}
}

func WithAppName(appname string) Option {
	return func (c *Client) {
		c.appname = appname
	}
}

func WithBufferSize(bufferSize int) Option {
	return func (c *Client) {
		c.bufferSize = bufferSize
	}
}
