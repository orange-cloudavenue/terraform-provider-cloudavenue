package provider

import (
	"context"
	"strings"
)

type jobStatusMessage string

const (
	DONE       jobStatusMessage = "DONE"
	FAILED     jobStatusMessage = "FAILED"
	CREATED    jobStatusMessage = "CREATED"
	PENDING    jobStatusMessage = "PENDING"
	INPROGRESS jobStatusMessage = "IN_PROGRESS"
	ERROR      jobStatusMessage = "ERROR"
)

// getJobStatus is a helper function to get the status of a job.
func getJobStatus(
	ctx context.Context,
	client *CloudAvenueClient,
	jobID string,
) (jobStatusMessage, error) {
	jobStatus, _, err := client.JobsApi.ApiCustomersV10JobsJobIdGet(ctx, jobID)
	if err != nil {
		return "", err
	}
	return parseJobStatus(jobStatus[0].Status), nil
}

// parseJobStatus return the status of a job.
func parseJobStatus(str string) jobStatusMessage {
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
func (j jobStatusMessage) string() string {
	return strings.ToLower(string(j))
}

// isDone is a helper function to check if a job is done.
func (j jobStatusMessage) isDone() bool {
	return j == DONE
}

// jobStatePending is a helper function to return an array of pending states.
func jobStatePending() []string {
	return []string{CREATED.string(), INPROGRESS.string(), PENDING.string()}
}

// jobStateDone is a helper function to return an array of done states.
func jobStateDone() []string {
	return []string{DONE.string()}
}
