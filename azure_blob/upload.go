package azure_blob

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

const (
	MaxFileSize = 25000000
)

type AzureBlobClient struct {
	serviceURL string
	creds      azcore.TokenCredential
	client     *azblob.Client
}

func NewAzureBlobClient(serviceURL string, creds azcore.TokenCredential) (*AzureBlobClient, error) {
	client, err := azblob.NewClient(serviceURL, creds, &azblob.ClientOptions{})
	if err != nil {
		return nil, err
	}
	return &AzureBlobClient{serviceURL: serviceURL, creds: creds, client: client}, nil
}

func (a *AzureBlobClient) UploadFile(ctx context.Context, containerName string, path string, buffer []byte) error {
	if len(buffer) <= 0 || len(buffer) > MaxFileSize {
		return fmt.Errorf("invalid resume size")
	}
	_, err := a.client.UploadBuffer(ctx, containerName, path, buffer, &azblob.UploadFileOptions{})
	if err != nil {
		return err
	}
	return nil
}

func (a *AzureBlobClient) DownloadFile(ctx context.Context, containerName string, path string) (buffer []byte, err error) {
	_, err = a.client.DownloadBuffer(ctx, containerName, path, buffer, &azblob.DownloadBufferOptions{})
	if err != nil {
		return nil, err
	}
	return buffer, nil
}

func (a *AzureBlobClient) UploadResume(ctx context.Context, userId string, hackathonId string, resumeBytes []byte) error {
	return a.UploadFile(ctx, "resumes", getResumeBlobName(hackathonId, userId), resumeBytes)
}

func (a *AzureBlobClient) DownloadResume(ctx context.Context, userId string, hackathonId string) ([]byte, error) {
	return a.DownloadFile(ctx, "resumes", getResumeBlobName(hackathonId, userId))
}

func getResumeBlobName(hackathonId string, userId string) string {
	return fmt.Sprintf("resumes/%s/%s", hackathonId, userId)
}
