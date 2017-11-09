package wrapper

import (
	"bytes"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"time"

	"github.com/micro/go-micro/server"
	"golang.org/x/net/context"
)

// GracefulRestartWrapper 优雅重启
type GracefulRestartWrapper struct {
	wg sync.WaitGroup
}

// HandlerWrapper HandlerWrapper
func (grw *GracefulRestartWrapper) HandlerWrapper() server.HandlerWrapper {
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			grw.wg.Add(1)
			defer grw.wg.Done()

			ch := make(chan error)
			go func() {
				defer func() {
					if e := recover(); e != nil {
						log.Println(e, "\t", string(bytes.Replace(debug.Stack(), []byte("\n"), []byte("\t"), -1)))
						ch <- fmt.Errorf("Panic: %+v", e)
					}
				}()
				ch <- fn(ctx, req, rsp)
			}()

			timeout := time.NewTimer(5 * time.Second)
			select {
			case err := <-ch:
				timeout.Stop()
				return err
			case <-timeout.C:
				return nil
			}
		}
	}
}

// Wait Wait
func (grw *GracefulRestartWrapper) Wait() {
	grw.wg.Wait()
}
