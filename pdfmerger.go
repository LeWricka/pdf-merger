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
	"path/filepath"
)

func Merge(w http.ResponseWriter, request *http.Request) {
	var files []string
	error := json.NewDecoder(request.Body).Decode(&files)
	if error != nil {
		panic(error)
	}
	log.Println(files)
	//printDirectory()
	outFiles := getGCSFilesNew(files)
	//if _, err := os.Stat("tmp"); os.IsNotExist(err) {
	//	error = os.Mkdir("tmp", 0755)
	//	if err != nil {
	//		fmt.Printf("Unable to create directory: %v", err)
	//	}
	//}
	//_,err := os.OpenFile("./tmp/out.pdf", os.O_CREATE, 0755)
	//data := [0]byte{}
	//error = ioutil.WriteFile("out.pdf", data, 0664)

	//if err != nil {
	//	fmt.Printf("Unable to write file: %v", err)
	//}

	outFile := createOutfile()
	error = api.MergeCreateFile(outFiles, outFile, pdfcpu.NewDefaultConfiguration())
	//printDirectory()
	if error != nil {
		log.Println(error)
		log.Println("error merging files")
	}

	SendPdf(w, outFile)
	for _, fileName := range files {
		log.Println(fileName + " removed")
			os.Remove(fileName)
	}
}

func createOutfile() string {
	content := []byte("temporary file's content")
	tmpfile, err := ioutil.TempFile("/tmp", "out.pdf")
	if err != nil {
		log.Fatal(err)
	}

	//defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}

	return tmpfile.Name()
}

func printDirectory() {
	log.Println("SHOW TREE:")

	err := filepath.Walk("/tmp",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			fmt.Println(path, info.Size())
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func SendPdf(w http.ResponseWriter, outfile string) {
	f, err := os.Open(outfile)
	if err != nil {
		fmt.Printf("error opening out file %s\n", err)
		w.WriteHeader(500)
		return
	}
	defer f.Close()

	//Set header
	w.Header().Set("Content-type", "application/pdf")

	//Stream to response
	if _, err := io.Copy(w, f); err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
	}

	//e := os.Remove("out.pdf")
	//if e != nil {
	//	log.Fatal(e)
	//}
}
