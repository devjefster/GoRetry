package retry

import (
	"math"
	"math/rand"
	"time"
)

func CalculateBackoff(cfg Config, attempt int) time.Duration {
	switch cfg.BackoffStrategy {
	case FixedBackoff:
		return cfg.InitialBackoff
	case LinearBackoff:
		delay := cfg.InitialBackoff + time.Duration(attempt)*cfg.InitialBackoff
		if delay > cfg.MaxBackoff {
			return cfg.MaxBackoff
		}
		return delay
	case ExponentialBackoff:
		delay := float64(cfg.InitialBackoff) * math.Pow(cfg.BackoffFactor, float64(attempt))
		if time.Duration(delay) > cfg.MaxBackoff {
			delay = float64(cfg.MaxBackoff)
		}
		jitter := delay * cfg.JitterFactor * (rand.Float64()*2 - 1)
		return time.Duration(delay + jitter)
	case CustomBackoff:
		if cfg.CustomBackoffFunc != nil {
			return cfg.CustomBackoffFunc(attempt)
		}
	}
	return cfg.InitialBackoff
}
