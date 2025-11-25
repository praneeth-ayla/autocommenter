package providerutil

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/googleapis/gax-go/v2/apierror"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
)

const (
	DefaultRetryDelay = 5 * time.Second
	MaxRetryAttempts  = 3
	PerRequestTimeout = 60 * time.Second
)

// DoWithRetry runs fn with retry, timeout and rate-limit handling.
// fn itself is synchronous; this helper runs it in a goroutine and enforces timeout.
func DoWithRetry[T any](maxAttempts int, timeout time.Duration, fn func() (T, error)) (T, error) {
	var zero T
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		done := make(chan struct{})
		var result T

		go func() {
			result, lastErr = fn()
			close(done)
		}()

		select {
		case <-ctx.Done():
			lastErr = fmt.Errorf("request timed out after %s", timeout)
		case <-done:
			// fn completed, lastErr is already set
		}
		cancel()

		if lastErr == nil {
			return result, nil
		}

		if delay, isRateLimit := CheckRateLimitError(lastErr); isRateLimit {
			if attempt < maxAttempts {
				SleepWithJitter(delay)
				continue
			}
			return zero, fmt.Errorf("rate limit after %d attempts: %w", maxAttempts, lastErr)
		}

		// Non-rate-limit error: fail fast
		return zero, fmt.Errorf("operation failed: %w", lastErr)
	}

	return zero, fmt.Errorf("retries exhausted after %d attempts: %w", maxAttempts, lastErr)
}

func SleepWithJitter(base time.Duration) {
	if base <= 0 {
		base = DefaultRetryDelay
	}
	j := time.Duration(rand.Int63n(int64(base / 2))) // jitter up to half of base
	time.Sleep(base + j)
}

func CheckRateLimitError(err error) (time.Duration, bool) {
	if err == nil {
		return 0, false
	}

	errStr := err.Error()
	if strings.Contains(errStr, "RESOURCE_EXHAUSTED") ||
		strings.Contains(errStr, "429") ||
		strings.Contains(errStr, "Quota exceeded") {
		return extractRetryDelay(errStr), true
	}

	var apiErr *apierror.APIError
	if errors.As(err, &apiErr) {
		if status := apiErr.GRPCStatus(); status != nil && status.Code() == codes.ResourceExhausted {
			for _, detail := range status.Details() {
				if retryInfo, ok := detail.(*errdetails.RetryInfo); ok {
					if retryInfo.RetryDelay != nil {
						if d := retryInfo.RetryDelay.AsDuration(); d > 0 {
							return d, true
						}
					}
				}
			}
			return DefaultRetryDelay, true
		}
	}

	return 0, false
}

func extractRetryDelay(errStr string) time.Duration {
	re := regexp.MustCompile(`retry in ([0-9.]+)s`)
	if matches := re.FindStringSubmatch(errStr); len(matches) > 1 {
		if seconds, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return time.Duration(seconds * float64(time.Second))
		}
	}
	return DefaultRetryDelay
}
