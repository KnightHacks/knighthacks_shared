package azure_blob

import (
	"context"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

type Resume struct {
	FileName string
	Data     []byte
}

const (
	MAX_FILE_SIZE = 25000000
)

/* upload resume to azure and return ID */
func UploadResume(ctx context.Context, resume *Resume, serviceURL string) (string, error) {
	if len(resume.Data) <= 0 || len(resume.Data) > MAX_FILE_SIZE {
		return "", fmt.Errorf("invalid resume size")
	}

	// write resume data to a temporary file
	err := os.WriteFile(resume.FileName, resume.Data, 0666)
	if err != nil {
		return "", err
	}

	fileHandler, err := os.Open(resume.FileName)
	if err != nil {
		return "", err
	}
	defer fileHandler.Close()

	// delete file on local machine after upload
	defer func() error {
		err = os.Remove(resume.FileName)
		if err != nil {
			return err
		}
		return nil
	}()

	// TODO: placeholder - need function to obtain credentials
	var cred azcore.TokenCredential

	client, err := azblob.NewClient(serviceURL, cred, nil)
	if err != nil {
		return "", err
	}

	// upload file
	resp, err := client.UploadFile(ctx, "testcontainer", "example/path/blobname", fileHandler, &azblob.UploadFileOptions{})
	if err != nil {
		return "", err
	}

	return *resp.RequestID, err
}

func DownloadResume(ctx context.Context) {

}
