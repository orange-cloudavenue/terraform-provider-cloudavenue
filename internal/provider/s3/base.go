package s3

import (
	"context"
	"errors"
	"time"

	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const (
	categoryName = "s3"

	// ErrCodeObjectLockConfigurationNotFoundError object lock configuration not found.
	ErrCodeObjectLockConfigurationNotFoundError = "ObjectLockConfigurationNotFoundError"
)

// DefaultWaitRetryInterval is used to set the retry interval to 0 during acceptance tests.
var DefaultWaitRetryInterval *time.Duration

// NotFound returns true if the error represents a "resource not found" condition.
// Specifically, NotFound returns true if the error or a wrapped error is of type
// retry.NotFoundError.
func NotFound(err error) bool {
	var e *retry.NotFoundError // nosemgrep:ci.is-not-found-error
	return errors.As(err, &e)
}

type RetryWhenConfig[T any] struct {
	Timeout  time.Duration
	Interval time.Duration
	Function func() (T, error)
}

var ErrRetryWhenTimeout = errors.New("timeout reached")

// retryWhen executes the function passed in the configuration object until the timeout is reached or the context is cancelled.
// It will retry if the shouldRetry function returns true. It will stop if the shouldRetry function returns false.
func retryWhen[T any](ctx context.Context, config *RetryWhenConfig[T], shouldRetry func(error) bool) (T, error) { //nolint: ireturn,unused
	retryInterval := config.Interval
	if DefaultWaitRetryInterval != nil {
		retryInterval = *DefaultWaitRetryInterval
	}

	timer := time.NewTimer(config.Timeout)

	for {
		result, err := config.Function()
		if shouldRetry(err) {
			select {
			case <-timer.C:
				return result, ErrRetryWhenTimeout
			case <-ctx.Done():
				return result, ctx.Err()
			default:
				time.Sleep(retryInterval) // lintignore:R018
				continue
			}
		}

		return result, err
	}
}

// retryWhenAWSErrCodeEquals retries a function when it returns a specific AWS error.
func retryWhenAWSErrCodeEquals[T any](ctx context.Context, codes []string, config *RetryWhenConfig[T]) (T, error) { //nolint: ireturn,unused
	return retryWhen(ctx, config, func(err error) bool {
		return tfawserr.ErrCodeEquals(err, codes...)
	})
}

// retryWhenAWSErrCodeNotEquals retries a function until it returns a specific AWS error.
func retryWhenAWSErrCodeNotEquals[T any](ctx context.Context, codes []string, config *RetryWhenConfig[T]) (T, error) { //nolint: ireturn,unused
	return retryWhen(ctx, config, func(err error) bool {
		if err == nil {
			return true
		}

		return !tfawserr.ErrCodeEquals(err, codes...)
	})
}
