package wrapper

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	microErrors "github.com/micro/go-micro/errors"
	"github.com/micro/go-micro/server"
	"github.com/micro/go-os/metrics"
	"golang.org/x/net/context"
)

type ServerStats struct {
	sync.RWMutex
	serverID   string
	instanceID string

	connNum    int64
	requestNum uint64

	m                 metrics.Metrics
	connNumGauge      metrics.Gauge
	requestNumCounter metrics.Counter
	errorCounter      map[int32]metrics.Counter
}

func NewServerStats(m metrics.Metrics, ip, serviceName, serviceID string) *ServerStats {
	serverID := strings.Replace(serviceName, ".", "-", -1)
	instanceID := strings.Replace(ip, ".", "-", -1) + "." + serviceID

	return &ServerStats{
		m:                 m,
		serverID:          serverID,
		instanceID:        instanceID,
		connNumGauge:      m.Gauge("connection_num." + serverID),
		requestNumCounter: m.Counter("request_num." + serverID),
		errorCounter:      make(map[int32]metrics.Counter),
	}
}

func (s *ServerStats) HandleStats(fn func() error) error {
	atomic.AddInt64(&s.connNum, 1)
	defer atomic.AddInt64(&s.connNum, -1)
	atomic.AddUint64(&s.requestNum, 1)

	err := fn()
	if err != nil {
		if e, ok := err.(*microErrors.Error); ok {
			s.RLock()
			c, ok := s.errorCounter[e.Code]
			if !ok {
				s.RUnlock()
				s.Lock()
				c, ok = s.errorCounter[e.Code]
				if !ok {
					c = s.m.Counter("error_num." + s.serverID + "." + strconv.Itoa(int(e.Code)))
					s.errorCounter[e.Code] = c
				}
				s.Unlock()
			} else {
				s.RUnlock()
			}
			c.Incr(1)
		}
		return err
	}
	return nil
}

func (s *ServerStats) Run() {
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			s.connNumGauge.Set(atomic.LoadInt64(&s.connNum))
			s.requestNumCounter.Incr(atomic.LoadUint64(&s.requestNum))
			atomic.StoreUint64(&s.requestNum, 0)
		}
	}()
}

func (s *ServerStats) HandlerStatsWrapper() server.HandlerWrapper {
	s.Run()
	return func(fn server.HandlerFunc) server.HandlerFunc {
		return func(ctx context.Context, req server.Request, rsp interface{}) error {
			return s.HandleStats(func() error {
				return fn(ctx, req, rsp)
			})
		}
	}
}
