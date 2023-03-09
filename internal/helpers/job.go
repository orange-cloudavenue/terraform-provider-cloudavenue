// Package helpers provides helpers for the CloudAvenue Terraform Provider.
package helpers

import (
	"context"
	"errors"
	"strings"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

// JobStatusMessage is a type for job status.
type JobStatusMessage string

const (
	// DONE is the done status message.
	DONE JobStatusMessage = "DONE"
	// FAILED is the failed status message.
	FAILED JobStatusMessage = "FAILED"
	// CREATED is the created status message.
	CREATED JobStatusMessage = "CREATED"
	// PENDING is the pending status message.
	PENDING JobStatusMessage = "PENDING"
	// INPROGRESS is the in progress status message.
	INPROGRESS JobStatusMessage = "IN_PROGRESS"
	// ERROR is the error status message.
	ERROR JobStatusMessage = "ERROR"
)

// GetJobStatus is a helper function to get the status of a job.
func GetJobStatus(
	ctx context.Context,
	client *client.CloudAvenue,
	jobID string,
) (JobStatusMessage, error) {
	jobStatus, httpR, err := client.APIClient.JobsApi.GetJobById(ctx, jobID)
	if err != nil {
		return "", err
	}
	defer func() {
		err = errors.Join(err, httpR.Body.Close())
	}()

	// Find the action name with failed status if global status is failed
	if jobStatus[0].Status == string(FAILED) {
		for _, action := range jobStatus[0].Actions {
			if action.Status == string(FAILED) {
				return parseJobStatus(jobStatus[0].Status), errors.New("Error in action : " + action.Name)
			}
		}
	}

	return parseJobStatus(jobStatus[0].Status), nil
}

// parseJobStatus return the status of a job.
func parseJobStatus(str string) JobStatusMessage {
	switch str {
	case "DONE":
		return DONE
	case "FAILED":
		return FAILED
	case "CREATED":
		return CREATED
	case "PENDING":
		return PENDING
	case "IN_PROGRESS":
		return INPROGRESS
	default:
		return ""
	}
}

// String is a stringer interface for jobStatus.
func (j JobStatusMessage) String() string {
	return strings.ToLower(string(j))
}

// IsDone is a helper function to check if a job is done.
func (j JobStatusMessage) IsDone() bool {
	return j == DONE
}

// JobStatePending is a helper function to return an array of pending states.
func JobStatePending() []string {
	return []string{CREATED.String(), INPROGRESS.String(), PENDING.String()}
}

// JobStateDone is a helper function to return an array of done states.
func JobStateDone() []string {
	return []string{DONE.String()}
}
