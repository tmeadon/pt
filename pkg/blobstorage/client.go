package blobstorage

import (
	"context"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var container *azblob.ContainerClient

type cloudBlob struct {
	Name         string
	LastModified time.Time
}

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

func DownloadNewest(dstFilePath string) error {
	if container == nil {
		return fmt.Errorf("container client has not been initialised")
	}
	blobs := listAllBlobs()
	newestBlob := findNewestBlob(&blobs)
	err := downloadBlob(newestBlob.Name, dstFilePath)
	return err
}

func listAllBlobs() []cloudBlob {
	allBlobs := make([]cloudBlob, 0)

	pager := container.ListBlobsFlat(&azblob.ContainerListBlobsFlatOptions{})

	for pager.NextPage(context.Background()) {
		pageItems := pager.PageResponse().ListBlobsFlatSegmentResponse.Segment.BlobItems
		for _, p := range pageItems {
			allBlobs = append(allBlobs, cloudBlob{*p.Name, *p.Properties.LastModified})
		}
	}

	return allBlobs
}

func findNewestBlob(blobs *[]cloudBlob) cloudBlob {
	newestBlob := (*blobs)[0]
	for _, b := range *blobs {
		if b.LastModified.After(newestBlob.LastModified) {
			newestBlob = b
		}
	}
	return newestBlob
}

func downloadBlob(name string, dstFilePath string) error {
	// check if file already exists
	if _, err := os.Stat(dstFilePath); err == nil {
		return fmt.Errorf("file already exists at path %s", dstFilePath)
	}

	// create the file
	file, err := os.Create(dstFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// create a blob client and download the file
	blobClient, err := container.NewBlockBlobClient(name)
	if err != nil {
		return err
	}

	err = blobClient.DownloadToFile(context.Background(), 0, 0, file, azblob.DownloadOptions{})
	return err
}
