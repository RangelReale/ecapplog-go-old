package ecapplog

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"net"
	"time"
)

type Client struct {
	address      string
	isOpen       bool
	cmdChan      chan interface{}
	cmdCtx       context.Context
	cmdCtxCancel context.CancelFunc
}

func NewClient(options ...Option) *Client {
	ret := &Client{
		address: "127.0.0.1:13991",
		isOpen:  false,
	}
	for _, opt := range options {
		opt(ret)
	}
	return ret
}

func (c *Client) Open() {
	if !c.isOpen {
		c.cmdChan = make(chan interface{})
		c.cmdCtx, c.cmdCtxCancel = context.WithCancel(context.Background())
		c.isOpen = true

		go c.handleConnection(c.cmdCtx, c.cmdChan)
	}
}

func (c *Client) Close() {
	if c.isOpen {
		c.cmdCtxCancel()
		close(c.cmdChan)

		c.cmdChan = nil
		c.cmdCtx = nil
		c.cmdCtxCancel = nil
		c.isOpen = false
	}
}

func (c *Client) handleConnection(c_cmdCtx context.Context, c_cmdChan chan interface{}) {
	rfor:
	for {
		err := func() error {
			var d net.Dialer
			cctx, ccancel := context.WithTimeout(c_cmdCtx, time.Second * 10)
			defer ccancel()
			conn, err := d.DialContext(cctx, "tcp", c.address)
			if err != nil {
				return err
			}
			defer conn.Close()

			if conntcp, ok := conn.(*net.TCPConn); ok {
				err = conntcp.SetNoDelay(false)
				if err != nil {
					return err
				}
			}

			// write banner
			_, err = conn.Write([]byte("ECAPPLOG"))
			if err != nil {
				return err
			}

			//logs := make([]interface{}, 0)
			//out := make(chan interface{})
			//outCh := func() chan interface{} {
			//	if len(logs) == 0 {
			//		return nil
			//	}
			//	return out
			//}
			//curVal := func() interface{} {
			//	if len(logs) == 0 {
			//		return nil
			//	}
			//	return logs[0]
			//}

			for {
				var err error
				select {
				case <-c_cmdCtx.Done():
					return nil
				case cmd := <-c_cmdChan:
					switch xcmd := cmd.(type) {
					case *cmdLog:
						err = c.handleCmdLog(conn, xcmd)
					}
				}
				if err != nil {
					if errors.Is(err, net.ErrClosed) {
						return nil
					}
					c.handleError(err)
				}
			}
		}()
		if err != nil {
			c.handleError(err)
		}

		select {
		case <-c_cmdCtx.Done():
			break rfor
		case <-time.After(time.Second * 5):
			break
		}
	}
}

func (c *Client) handleCmdLog(conn net.Conn, cmd *cmdLog) error {
	jcmd, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	// write command
	err = binary.Write(conn, binary.BigEndian, command_Log)
	if err != nil {
		return err
	}

	// write size
	size := int32(len(jcmd))
	err = binary.Write(conn, binary.BigEndian, size)
	if err != nil {
		return err
	}

	// write encoded json
	err = binary.Write(conn, binary.BigEndian, jcmd)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) handleError(err error) {

}

func (c *Client) Log(time time.Time, priority Priority, source string, text string) {
	if c.isOpen {
		c.cmdChan <- &cmdLog{
			Time:     cmdTime{time},
			Priority: priority,
			Source:   source,
			Text:     text,
		}
	}
}
