package ecapplog

type Option func(*Client)

func WithAddress(address string) Option {
	return func (c *Client) {
		c.address = address
	}
}
