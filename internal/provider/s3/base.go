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

	// defaultReadTimeout is the default timeout for read operations.
	defaultReadTimeout = 5 * time.Minute
	// defaultCreateTimeout is the default timeout for create operations.
	defaultCreateTimeout = 5 * time.Minute
	// defaultUpdateTimeout is the default timeout for update operations.
	defaultUpdateTimeout = 5 * time.Minute
	// defaultDeleteTimeout is the default timeout for delete operations.
	defaultDeleteTimeout = 5 * time.Minute

	// General timeout for S3 bucket changes to propagate.
	// See https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html#ConsistencyModel.
	s3BucketPropagationTimeout = 2 * time.Minute // nosemgrep:ci.s3-in-const-name, ci.s3-in-var-name
)

// Error code constants missing from AWS Go SDK:
// https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#pkg-constants

const (
	ErrCodeAccessDenied                         = "AccessDenied"
	ErrCodeBucketNotEmpty                       = "BucketNotEmpty"
	ErrCodeInvalidBucketState                   = "InvalidBucketState"
	ErrCodeInvalidRequest                       = "InvalidRequest"
	ErrCodeMalformedPolicy                      = "MalformedPolicy"
	ErrCodeMethodNotAllowed                     = "MethodNotAllowed"
	ErrCodeNoSuchBucket                         = "NoSuchBucket"
	ErrCodeNoSuchBucketPolicy                   = "NoSuchBucketPolicy"
	ErrCodeNoSuchConfiguration                  = "NoSuchConfiguration"
	ErrCodeNoSuchCORSConfiguration              = "NoSuchCORSConfiguration"
	ErrCodeNoSuchLifecycleConfiguration         = "NoSuchLifecycleConfiguration"
	ErrCodeNoSuchKey                            = "NoSuchKey"
	ErrCodeNoSuchPublicAccessBlockConfiguration = "NoSuchPublicAccessBlockConfiguration"
	ErrCodeNoSuchTagSet                         = "NoSuchTagSet"
	ErrCodeNoSuchTagSetError                    = "NoSuchTagSetError"
	ErrCodeNoSuchWebsiteConfiguration           = "NoSuchWebsiteConfiguration"
	ErrCodeNotImplemented                       = "NotImplemented"
	// ErrCodeObjectLockConfigurationNotFound should be used with tfawsErr.ErrCodeContains, not tfawsErr.ErrCodeEquals.
	// Reference: https://github.com/hashicorp/tErraform-provider-aws/pull/26317
	ErrCodeObjectLockConfigurationNotFound           = "ObjectLockConfigurationNotFound"
	ErrCodeObjectLockConfigurationNotFoundError      = "ObjectLockConfigurationNotFoundError"
	ErrCodeOperationAborted                          = "OperationAborted"
	ErrCodeOwnershipControlsNotFoundError            = "OwnershipControlsNotFoundError"
	ErrCodeReplicationConfigurationNotFound          = "ReplicationConfigurationNotFoundError"
	ErrCodeServerSideEncryptionConfigurationNotFound = "ServerSideEncryptionConfigurationNotFoundError"
	ErrCodeUnsupportedArgument                       = "UnsupportedArgument"
	// ErrCodeXNotImplemented is returned from Third Party S3 implementations
	// and so far has been noticed with calls to GetBucketWebsite.
	// Reference: https://github.com/hashicorp/tErraform-provider-aws/issues/14645
	ErrCodeXNotImplemented = "XNotImplemented"
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

const (
	lifeCycleStatusEnabled  = "Enabled"
	lifeCycleStatusDisabled = "Disabled"
)

var ErrRetryWhenTimeout = errors.New("timeout reached")

// retryWhen executes the function passed in the configuration object until the timeout is reached or the context is cancelled.
// It will retry if the shouldRetry function returns true. It will stop if the shouldRetry function returns false.
func retryWhen[T any](ctx context.Context, config *RetryWhenConfig[T], shouldRetry func(error) bool) (T, error) { //nolint: ireturn
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
func retryWhenAWSErrCodeEquals[T any](ctx context.Context, codes []string, config *RetryWhenConfig[T]) (T, error) { //nolint: ireturn
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

// RetryWhenNotFound retries the specified function when it returns a retry.NotFoundError.
func retryWhenNotFound(ctx context.Context, timeout time.Duration, f func() (interface{}, error)) (interface{}, error) { //nolint: ireturn
	return retryWhen(ctx, &RetryWhenConfig[interface{}]{Timeout: timeout, Function: f}, NotFound)
}
