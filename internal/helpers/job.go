package helpers

import (
	"context"
	"strings"

	"github.com/orange-cloudavenue/terraform-provider-cloudavenue/internal/client"
)

type JobStatusMessage string

const (
	DONE       JobStatusMessage = "DONE"
	FAILED     JobStatusMessage = "FAILED"
	CREATED    JobStatusMessage = "CREATED"
	PENDING    JobStatusMessage = "PENDING"
	INPROGRESS JobStatusMessage = "IN_PROGRESS"
	ERROR      JobStatusMessage = "ERROR"
)

// GetJobStatus is a helper function to get the status of a job.
func GetJobStatus(
	ctx context.Context,
	client *client.CloudAvenue,
	jobID string,
) (JobStatusMessage, error) {
	jobStatus, _, err := client.APIClient.JobsApi.ApiCustomersV10JobsJobIdGet(ctx, jobID)
	if err != nil {
		return "", err
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

// string is a stringer interface for jobStatus
func (j JobStatusMessage) String() string {
	return strings.ToLower(string(j))
}

// isDone is a helper function to check if a job is done.
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
