package wrapper

import (
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/client"

	microErrors "github.com/micro/go-micro/errors"
	"golang.org/x/net/context"
)

func init() {
	hystrix.DefaultTimeout = int(client.DefaultRequestTimeout)
	hystrix.DefaultMaxConcurrent = 5000
}

type clientWrapper struct {
	client.Client
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	errChan := make(chan error, 1)
	errors := hystrix.Go(req.Service()+"."+req.Method(), func() error {
		err := c.Client.Call(ctx, req, rsp, opts...)
		if err == nil {
			errChan <- nil
			return nil
		}
		e := microErrors.Parse(err.Error())
		if e.Code < 1000 {
			return err
		}
		errChan <- err
		return nil
	}, nil)

	select {
	case err := <-errChan:
		return err
	case err := <-errors:
		return err
	}
}

func (c *clientWrapper) CallRemote(ctx context.Context, addr string, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	errChan := make(chan error, 1)

	errors := hystrix.Go(addr+"."+req.Service()+"."+req.Method(), func() error {
		err := c.Client.CallRemote(ctx, addr, req, rsp, opts...)
		if err == nil {
			errChan <- nil
			return nil
		}
		e := microErrors.Parse(err.Error())
		if e.Code < 1000 {
			return err
		}
		errChan <- err
		return nil
	}, nil)

	select {
	case err := <-errChan:
		return err
	case err := <-errors:
		return err
	}
}

// NewHystrixClientWrapper returns a hystrix client Wrapper.
func NewHystrixClientWrapper() client.Wrapper {
	return func(c client.Client) client.Client {
		return &clientWrapper{c}
	}
}
