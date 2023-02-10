package provider

import (
	"context"
	"strings"
)

type JobStatus string

const (
	DONE       JobStatus = "DONE"
	FAILED     JobStatus = "FAILED"
	CREATED    JobStatus = "CREATED"
	PENDING    JobStatus = "PENDING"
	INPROGRESS JobStatus = "IN_PROGRESS"
	ERROR      JobStatus = "ERROR"
)

// getJobStatus is a helper function to get the status of a job.
func getJobStatus(ctx context.Context, client *CloudAvenueClient, jobID string) (JobStatus, error) {
	jobStatus, _, err := client.JobsApi.ApiCustomersV10JobsJobIdGet(ctx, jobID)
	if err != nil {
		return "", err
	}
	return parseJobStatus(jobStatus[0].Status), nil
}

// parseJobStatus return the status of a job.
func parseJobStatus(str string) JobStatus {
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

// Stringer interface for JobStatus
func (j JobStatus) String() string {
	return strings.ToLower(string(j))
}

// IsDone is a helper function to check if a job is done.
func (j JobStatus) IsDone() bool {
	return j == DONE
}

// StatePending is a helper function to return an array of pending states.
func StatePending() []string {
	return []string{CREATED.String(), INPROGRESS.String(), PENDING.String()}
}
