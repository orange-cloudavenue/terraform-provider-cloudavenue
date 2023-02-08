package provider

import "context"

// getJobStatus is a helper function to get the status of a job.
func getJobStatus(ctx context.Context, client *CloudAvenueClient, jobID string) (string, error) {
	jobStatus, _, err := client.JobsApi.ApiCustomersV10JobsJobIdGet(ctx, jobID)
	if err != nil {
		return "", err
	}
	return jobStatus[0].Status, nil
}
