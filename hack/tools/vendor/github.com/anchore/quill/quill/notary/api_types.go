package notary

import "time"

// Common

type submissionResponseDescriptor struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// Request

type submissionRequest struct {
	Sha256         string `json:"sha256"`
	SubmissionName string `json:"submissionName"`
}

type submissionResponse struct {
	Data submissionResponseData `json:"data"`
}

type submissionResponseData struct {
	submissionResponseDescriptor
	Attributes submissionResponseAttributes `json:"attributes"`
}

type submissionResponseAttributes struct {
	AwsAccessKeyID     string `json:"awsAccessKeyId"`
	AwsSecretAccessKey string `json:"awsSecretAccessKey"`
	AwsSessionToken    string `json:"awsSessionToken"`
	Bucket             string `json:"bucket"`
	Object             string `json:"object"`
}

// List

type submissionListResponse struct {
	Data []submissionListResponseData `json:"data"`
}

type submissionListResponseData struct {
	submissionResponseDescriptor
	Attributes submissionListResponseAttributes `json:"attributes"`
}

type submissionListResponseAttributes struct {
	CreatedDate string `json:"createdDate"`
	Name        string `json:"name"`
	Status      string `json:"status"`
}

// Status

type submissionStatusResponse struct {
	Data submissionStatusResponseData `json:"data"`
}

type submissionStatusResponseData struct {
	submissionResponseDescriptor
	Attributes submissionStatusResponseAttributes `json:"attributes"`
}

type submissionStatusResponseAttributes struct {
	Status      string    `json:"status"`
	Name        string    `json:"name"`
	CreatedDate time.Time `json:"createdDate"`
}

// Logs

type submissionLogsResponse struct {
	Data submissionLogsResponseData `json:"data"`
}

type submissionLogsResponseData struct {
	submissionResponseDescriptor
	Attributes submissionLogsResponseAttributes `json:"attributes"`
}

type submissionLogsResponseAttributes struct {
	DeveloperLogURL string `json:"developerLogUrl"`
}
