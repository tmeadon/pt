package blobstorage

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var container *azblob.ContainerClient

func ConnectContainer(sasToken string) error {
	con, err := azblob.NewContainerClientWithNoCredential(sasToken, nil)
	if err != nil {
		return err
	}

	container = con
	return nil
}

func Upload(srcPath string) error {
	// check src exists
	if _, err := os.Stat(srcPath); err != nil {
		return err
	}

	// open the src file
	file, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a blob client
	blobClient, err := container.NewBlockBlobClient(path.Base(srcPath))
	if err != nil {
		return err
	}

	// upload the file
	_, err = blobClient.UploadFile(context.Background(), file, azblob.UploadOption{})
	if err != nil {
		return err
	}

	fmt.Printf("Uploaded %s to %s\n", srcPath, blobClient.URL())
	return nil
}
