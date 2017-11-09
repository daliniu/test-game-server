package trace

import (
	"math/rand"
	"time"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/metadata"
	"github.com/micro/go-micro/registry"
	"github.com/micro/go-micro/server"

	"github.com/micro/go-os/config"
	"golang.org/x/net/context"
)

// Config 配置
type RateConfig struct {
	Numerator   int `json:"numerator"`
	Denominator int `json:"denominator"`
}

type clientWrapper struct {
	client.Client
	t      Trace
	s      *registry.Service
	cfg    config.Config
	isRoot bool
}

func (c *clientWrapper) CallRemote(ctx context.Context, address string, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if !c.cfg.Get("isTraceOn").Bool(false) {
		return c.Client.CallRemote(ctx, address, req, rsp, opts...)
	}
	var span *Span
	var ok, mk bool
	var err error
	var md metadata.Metadata

	if c.isRoot {
		conf := RateConfig{1, 1000}
		_ = c.cfg.Get("traceRate").Scan(&conf)
		isTrace := isHit(conf.Numerator, conf.Denominator)
		// do not trace
		if !isTrace {
			return c.Client.CallRemote(ctx, address, req, rsp, opts...)
		}

		md, mk = metadata.FromContext(ctx)
		if !mk {
			md = make(metadata.Metadata)
		}
		span = c.t.NewSpan(nil)

	} else {

		md, mk = metadata.FromContext(ctx)
		if !mk {
			md = make(metadata.Metadata)
		}
		// try pull span from context
		span, ok = SpanFromContext(ctx)
		if !ok && mk {
			span, _ = SpanFromHeader(md)
		}

		if span == nil {
			return c.Client.CallRemote(ctx, address, req, rsp, opts...)
		}

		// setup the span with parent
		span = c.t.NewSpan(&Span{
			// same trace id
			TraceId: span.TraceId,
			// set parent id to parent span id
			ParentId: span.Id,
			// use previous debug
			Debug: span.Debug,
		})
	}

	newSpan, newCtx := c.newSpanCtx(ctx, req, rsp, span, md, opts...)
	defer c.collect(newSpan, err)

	// now just make a regular call down the stack
	err = c.Client.CallRemote(newCtx, address, req, rsp, opts...)
	return err
}

func (c *clientWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	if !c.cfg.Get("isTraceOn").Bool(false) {
		return c.Client.Call(ctx, req, rsp, opts...)
	}

	var span *Span
	var ok, mk bool
	var err error
	var md metadata.Metadata

	if c.isRoot {
		conf := RateConfig{1, 1000}
		_ = c.cfg.Get("traceRate").Scan(&conf)
		isTrace := isHit(conf.Numerator, conf.Denominator)
		// do not trace
		if !isTrace {
			return c.Client.Call(ctx, req, rsp, opts...)
		}

		md, mk = metadata.FromContext(ctx)
		if !mk {
			md = make(metadata.Metadata)
		}
		span = c.t.NewSpan(nil)

	} else {
		md, mk = metadata.FromContext(ctx)
		if !mk {
			md = make(metadata.Metadata)
		}
		// try pull span from context
		span, ok = SpanFromContext(ctx)
		if !ok && mk {
			span, _ = SpanFromHeader(md)
		}

		if span == nil {
			return c.Client.Call(ctx, req, rsp, opts...)
		}

		// setup the span with parent
		span = c.t.NewSpan(&Span{
			// same trace id
			TraceId: span.TraceId,
			// set parent id to parent span id
			ParentId: span.Id,
			// use previous debug
			Debug: span.Debug,
		})
	}

	newSpan, newCtx := c.newSpanCtx(ctx, req, rsp, span, md, opts...)
	defer c.collect(newSpan, err)
	// now just make a regular call down the stack
	err = c.Client.Call(newCtx, req, rsp, opts...)
	return err
}

func (c *clientWrapper) newSpanCtx(ctx context.Context, req client.Request, rsp interface{}, span *Span, md metadata.Metadata, opts ...client.CallOption) (*Span, context.Context) {
	// start the span
	span.Annotations = append(span.Annotations, &Annotation{
		Timestamp: time.Now(),
		Type:      AnnStart,
		Service:   c.s,
	})

	// and mark as debug? might want to do this based on a setting
	span.Debug = true
	// set unique span name
	span.Name = req.Service() + "." + req.Method()
	// set source/dest
	span.Source = c.s
	span.Destination = &registry.Service{Name: req.Service()}

	// set context key
	newCtx := ContextWithSpan(ctx, span)
	// set metadata
	newCtx = metadata.NewContext(newCtx, HeaderWithSpan(md, span))

	// mark client request
	span.Annotations = append(span.Annotations, &Annotation{
		Timestamp: time.Now(),
		Type:      AnnClientRequest,
		Service:   c.s,
	})

	return span, newCtx
}

func (c *clientWrapper) collect(span *Span, err error) {
	// mark client response
	span.Annotations = append(span.Annotations, &Annotation{
		Timestamp: time.Now(),
		Type:      AnnClientResponse,
		Service:   c.s,
	})

	// if we were the creator
	var debug map[string]string
	if err != nil {
		debug = map[string]string{"error": err.Error()}
	}
	// mark end of span
	span.Annotations = append(span.Annotations, &Annotation{
		Timestamp: time.Now(),
		Type:      AnnEnd,
		Service:   c.s,
		Debug:     debug,
	})

	span.Duration = time.Now().Sub(span.Timestamp)

	// flush the span to the collector on return
	c.t.Collect(span)
}

func handlerWrapper(fn server.HandlerFunc, t Trace, s *registry.Service, cfg config.Config) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		if !cfg.Get("isTraceOn").Bool(false) {
			return fn(ctx, req, rsp)
		}
		// embed trace instance
		newCtx := NewContext(ctx, t)

		var span *Span
		var err error
		var ok, okk bool
		var md metadata.Metadata

		// get trace info from metadata
		md, ok = metadata.FromContext(ctx)
		if ok {
			span, okk = SpanFromHeader(md)
		}

		if !ok || !okk {
			return fn(ctx, req, rsp)
		}

		// mark client request
		span.Annotations = append(span.Annotations, &Annotation{
			Timestamp: time.Now(),
			Type:      AnnServerRequest,
			Service:   s,
		})

		// and mark as debug? might want to do this based on a setting
		span.Debug = true
		// set unique span name
		span.Name = req.Service() + "." + req.Method()
		// set source/dest
		span.Source = s
		span.Destination = &registry.Service{Name: req.Service()}

		// embed the span in the context
		newCtx = ContextWithSpan(newCtx, span)

		// defer the completion of the span
		defer func() {
			var debug map[string]string
			if err != nil {
				debug = map[string]string{"error": err.Error()}
			}
			// mark server response
			span.Annotations = append(span.Annotations, &Annotation{
				Timestamp: time.Now(),
				Type:      AnnServerResponse,
				Service:   s,
				Debug:     debug,
			})

			span.Duration = time.Now().Sub(span.Timestamp)

			// flush the span to the collector on return
			t.Collect(span)
		}()
		err = fn(newCtx, req, rsp)
		return err
	}
}

// use numerator/denominator
func isHit(numerator, denominator int) bool {
	if denominator <= 0 {
		denominator = 1000
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	random := rand.Intn(denominator)
	if random < numerator {
		return true
	}
	return false
}
