// SPDX-FileCopyrightText: Copyright (c) 2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

// SPDX-FileCopyrightText: Copyright (c) 2023-2026 Orange
// SPDX-License-Identifier: MPL-2.0
//
// This software is distributed under the MPL-2.0 license.
// the text of which is available at https://www.mozilla.org/en-US/MPL/2.0/
// or see the "LICENSE" file for more details.

package s3

import (
	"context"
	"errors"
	"time"

	awsv2 "github.com/aws/aws-sdk-go-v2/aws"
	awss3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"

	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	v1 "github.com/orange-cloudavenue/cloudavenue-sdk-go/v1"
	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

const (
	categoryName = "s3"

	s3BucketEndpointFormat  = "https://%s.s3-region01.cloudavenue.orange-business.com"
	s3WebsiteEndpointFormat = "%s.website-region01.cloudavenue.orange-business.com"

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
	// ErrCodeObjectLockConfigurationNotFound is returned as part of a longer error code.
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

type Client struct {
	*awss3.Client
	cavSDKV1Client v1.S3Client
}

func (c Client) GetUser(username string) (*v1.S3User, *v1.OSEError) {
	// Still delegated to cloudavenue-sdk-go/v1 until the upstream S3 domain is available in v2.
	return c.cavSDKV1Client.GetUser(username)
}

func (c Client) SyncBucket(bucketName string) error {
	// Still delegated to cloudavenue-sdk-go/v1 until the upstream S3 domain is available in v2.
	return c.cavSDKV1Client.SyncBucket(bucketName)
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
		return s3ErrCodeEquals(err, codes...)
	})
}

func s3ErrCodeEquals(err error, codes ...string) bool {
	if apiErr, ok := errors.AsType[smithy.APIError](err); ok {
		for _, code := range codes {
			if apiErr.ErrorCode() == code {
				return true
			}
		}
	}

	return false
}

// RetryWhenNotFound retries the specified function when it returns a retry.NotFoundError.
func retryWhenNotFound(ctx context.Context, timeout time.Duration, f func() (any, error)) (any, error) { //nolint: ireturn
	return retryWhen(ctx, &RetryWhenConfig[any]{Timeout: timeout, Function: f}, NotFound)
}

func initS3Client(client *client.CloudAvenue) (s3Client Client, diags diag.Diagnostics) {
	v1Client, err := client.CAVSDK.V1.S3()
	if err != nil {
		diags.AddError(
			"Error initializing S3 client",
			err.Error(),
		)

		return Client{}, diags
	}

	// Temporary bridge until cloudavenue-sdk-go-v2 migrates the S3 domain:
	// reuse the CloudAvenue SDK v1-resolved S3 settings to build an AWS SDK v2
	// client, while keeping the remaining upstream v1-only helper methods.
	cavSDKV1Config := v1Client.Config
	config := awsv2.Config{
		Region: awsv2.ToString(cavSDKV1Config.Region),
		// Reuse the CloudAvenue SDK v1 credential source until the upstream S3
		// domain is migrated using cloudavenue-sdk-go-v2.
		Credentials: awsv2.CredentialsProviderFunc(func(context.Context) (awsv2.Credentials, error) {
			cavSDKV1Credentials, err := cavSDKV1Config.Credentials.Get()
			if err != nil {
				return awsv2.Credentials{}, err
			}

			return awsv2.Credentials{
				AccessKeyID:     cavSDKV1Credentials.AccessKeyID,
				SecretAccessKey: cavSDKV1Credentials.SecretAccessKey,
				SessionToken:    cavSDKV1Credentials.SessionToken,
			}, nil
		}),
	}

	// Match the resolved CloudAvenue S3 endpoint configuration in the v2 client.
	endpoint := awsv2.ToString(cavSDKV1Config.Endpoint)
	clientV2 := awss3.NewFromConfig(config, func(options *awss3.Options) {
		options.BaseEndpoint = &endpoint
		options.UsePathStyle = true
	})
	return Client{Client: clientV2, cavSDKV1Client: v1Client}, diags
}
