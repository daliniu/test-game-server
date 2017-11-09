package stk

import "github.com/uber-go/zap"
import "time"

// NewLogger NewLogger
func NewLogger(name string, level zap.Level) zap.Logger {
	return zap.New(
		zap.NewTextEncoder(zap.TextNoTime()),
		level,
		zap.AddCaller(),
		zap.Hook(namedHook(name, 8*time.Hour)),
		// zap.AddStacks(zap.ErrorLevel),
	)
}

func namedHook(name string, t time.Duration) zap.Hook {
	return func(entry *zap.Entry) error {
		message := []byte{}
		message = entry.Time.Add(t).AppendFormat(message, "2006-01-02 15:04:05")
		message = append(message, ' ', '[')
		message = append(message, name...)
		message = append(message, ']', ' ')
		message = append(message, entry.Message...)
		entry.Message = string(message)
		return nil
	}
}
