package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	topimages "ws/top10"
)

func submit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method)
	if r.Method == "GET" {
		fmt.Println("GET")
	} else {
		r.ParseForm()
		fmt.Println("Subreddit:", r.PostFormValue("reddit"))
	}

	check := strings.Contains(r.PostFormValue("reddit"), "https://old.reddit.com/r/")
	if check == false {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/Submit.html", 301)
	savedImages = topimages.Body(r.PostFormValue("reddit"))

}

func download(w http.ResponseWriter, r *http.Request) {

	var files [20]string
	files = savedImages

	var test = strings.TrimSuffix(files[0], "\\")
	fmt.Println(test + "GOTCHA")
	output := "images.zip"

	if err := ZipFiles(output, files); err != nil {
		panic(err)
	}
	data, err2 := ioutil.ReadFile("images.zip")
	if err2 != nil {

		log.Fatal(err2)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename="+output)

	http.ServeContent(w, r, "images.zip", time.Now(), bytes.NewReader(data))

}

func ZipFiles(filename string, files [20]string) error {

	newZipFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range files {

		if file == "" {
			break
		}
		if err = AddFileToZip(zipWriter, file); err != nil {
			//broken here 4/5/19
			return err
		}
	}
	for _, element := range savedImages {
		defer os.Remove(element)
	}

	return nil
}

//AddFileToZip does what it sounds like it does
func AddFileToZip(zipWriter *zip.Writer, filename string) error {

	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

var savedImages [20]string

func main() {
	http.Handle("/", http.FileServer(http.Dir("./html")))
	http.HandleFunc("/Submit", submit)
	http.HandleFunc("/Download", download)
	err := http.ListenAndServe(":8080", nil) // set listen port

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
