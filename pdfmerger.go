package pdfmerger

import (
	"encoding/json"
	"fmt"
	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func Merge(w http.ResponseWriter, request *http.Request) {
	var files []string
	error := json.NewDecoder(request.Body).Decode(&files)
	if error != nil {
		log.Println(error)
	}
	log.Println(files)
	filesToBeMerged := getGCSFilesNew(files)

	mergedFile := createOutputFile()
	error = api.MergeCreateFile(filesToBeMerged, mergedFile, pdfcpu.NewDefaultConfiguration())
	if error != nil {
		log.Println(error)
	}

	sendFile(w, mergedFile)
	for _, fileName := range files {
		log.Println(fileName + " removed")
		os.Remove(fileName)
	}
}

func createOutputFile() string {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("/tmp", "out.pdf")
	if err != nil {
		log.Println(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		log.Println(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Println(err)
	}

	return tmpfile.Name()
}

func sendFile(w http.ResponseWriter, outfile string) {
	f, err := os.Open(outfile)
	if err != nil {
		fmt.Printf("error opening out file %s\n", err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	w.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}
}
