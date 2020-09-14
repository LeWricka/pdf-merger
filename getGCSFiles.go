package pdfmerger

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const bucket = "climbingplan.appspot.com"

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
		reader, err := client.Bucket(bucket).Object("Routines/" + fileName).NewReader(ctx)
		if err != nil {
			fmt.Errorf("NewReader error: %v", err)
		}
		defer reader.Close()

		data, error := ioutil.ReadAll(reader)
		if error != nil {
			fmt.Errorf("ioutil.ReadAll: %v", error)
		}

		tmpfile, err := ioutil.TempFile("/tmp", fileName)
		if err != nil {
			log.Fatal(err)
		}
		outFileNames = append(outFileNames, tmpfile.Name())
		//defer os.Remove(tmpfile.Name()) // clean up

		if _, err := tmpfile.Write(data); err != nil {
			log.Fatal(err)
		}
		if err := tmpfile.Close(); err != nil {
			log.Fatal(err)
		}

		//error = TempFile("tmp/",fileName, data, 0664)
		//if error != nil {
		//	fmt.Errorf("WriteFile: %v", error)
		//}
		fmt.Printf("Blob %s create.\n", fileName)
	}

	return outFileNames
}