package main

import (
    "archive/zip"
    "fmt"
    "net/http"
    "log"
    "os/exec"
    "bytes"
   "time"
    "io"
    "os"
    "io/ioutil"
)

func submit(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:",r.Method)
	if r.Method == "GET" {
		fmt.Println("GET")
	} else {
		r.ParseForm()
		fmt.Println("Subreddit:", r.PostFormValue("reddit"))
	}
	http.Redirect(w,r,"/Submit.html", 301)
	cmd := exec.Command("./top10",r.PostFormValue("reddit"))

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}

}

func download(w http.ResponseWriter, r *http.Request){

	files := []string{"images/0.jpg","images/1.jpg", "images/2.jpg","images/3.jpg","images/4.jpg","images/5.jpg","images/6.jpg","images/7.jpg","images/8.jpg","images/9.jpg"}
        output := "images.zip"

        if err := ZipFiles(output, files); err != nil {
                panic(err)
        }
    data, err2 := ioutil.ReadFile("images.zip")
    if(err2 != nil){
        log.Fatal(err2)
    }

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", "attachment; filename=" + "images.zip")
	http.ServeContent(w,r,"images.zip", time.Now(), bytes.NewReader(data))

}

//zip 
func ZipFiles(filename string, files []string) error {

    newZipFile, err := os.Create(filename)
    if err != nil {
        return err
    }
    defer newZipFile.Close()

    zipWriter := zip.NewWriter(newZipFile)
    defer zipWriter.Close()

    // Add files to zip
    for _, file := range files {
        if err = AddFileToZip(zipWriter, file); err != nil {
            return err
        }
    }
    return nil
}

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


func main() {
    http.Handle("/", http.FileServer(http.Dir("./html")))
    http.HandleFunc("/Submit", submit)
    http.HandleFunc("/Download", download)
err := http.ListenAndServe(":8080", nil) // set listen port

if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
