package pdfmerger

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func getGCSFilesNew(fileNames []string) []string {
	var outFileNames []string
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	for _, fileName := range fileNames {
		bucketName := os.Getenv("BUCKET")
		if bucketName == "" {
			log.Printf("Error reading bucket var")
		}
		pdfsPath := os.Getenv("PDFS_PATH")
		if pdfsPath == "" {
			log.Printf("Error reading pdfspath var")
		}
		reader, bucketFetchError := client.Bucket(bucketName).Object(pdfsPath + fileName).NewReader(ctx)
		if bucketFetchError != nil {
			log.Printf("Error fetching bucket file"+fileName)
		}
		defer reader.Close()

		data, readError := ioutil.ReadAll(reader)
		if readError != nil {
			log.Printf("ioutil.ReadAll: %v", readError)
		}

		tmpfile, fileCreationError := ioutil.TempFile("/tmp", fileName)
		if fileCreationError != nil {
			log.Printf("Error creating tmp file: %v", fileCreationError)
		}
		outFileNames = append(outFileNames, tmpfile.Name())

		if _, err := tmpfile.Write(data); err != nil {
			log.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Blob %s create.\n", fileName)
	}

	return outFileNames
}
