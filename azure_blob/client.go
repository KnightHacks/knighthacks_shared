package azure_blob

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/KnightHacks/knighthacks_shared/utils"
	"time"
)

type AzureBlobClient struct {
	credential *azblob.SharedKeyCredential
}

func (a *AzureBlobClient) CreatePreSignedURL(hackathonId string, userId string) (string, error) {
	sasQueryParams, err := sas.BlobSignatureValues{
		Protocol:      sas.ProtocolHTTPS,
		StartTime:     time.Now().UTC(),
		ExpiryTime:    time.Now().UTC().Add(5 * time.Minute),
		Permissions:   to.Ptr(sas.BlobPermissions{Create: true, Write: true, Add: true}).String(),
		ContainerName: "resumes",
		BlobName:      fmt.Sprintf("%s/%s.pdf", hackathonId, userId),
	}.SignWithSharedKey(a.credential)

	if err != nil {
		return "", nil
	}

	sasURL := fmt.Sprintf("https://%s.blob.core.windows.net/?%s", a.credential.AccountName(), sasQueryParams.Encode())
	return sasURL, nil
}

func (a *AzureBlobClient) GetResumeURL(hackathonID string, userID string) string {
	return fmt.Sprintf("https://%S.blob.core.windows.net/resumes/%s/%s.pdf", a.credential.AccountName(), hackathonID, userID)
}

func NewAzureBlobClient(credential *sas.SharedKeyCredential) (*AzureBlobClient, error) {
	return &AzureBlobClient{credential: credential}, nil
}

func NewSharedCredentialFromEnv() (*sas.SharedKeyCredential, error) {
	accountName := utils.GetEnvOrDie("AZURE_ACCOUNT_NAME")
	accountKey := utils.GetEnvOrDie("AZURE_ACCOUNT_KEY")

	credential, err := azblob.NewSharedKeyCredential(accountName, accountKey)
	if err != nil {
		return nil, err
	}
	return credential, nil
}
