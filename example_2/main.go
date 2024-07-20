package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type dataDownload struct {
	BaseUrl string
	First   int
	LastInt int
	Prefix  string
	Folder  string
}

func main() {

	var (
		ch           = make(chan string)
		retrieveData string
		dDownload    []dataDownload
	)

	dDownload = []dataDownload{
		{
			BaseUrl: "https://raw.githubusercontent.com/indraoct/go-concurrency/main/files/zee",
			First:   1,
			LastInt: 23,
			Prefix:  "frame_",
			Folder:  "/Users/indraoctama/Downloads/zee/",
		},
		{
			BaseUrl: "https://raw.githubusercontent.com/indraoct/go-concurrency/main/files/magic_hour",
			First:   1,
			LastInt: 96,
			Prefix:  "magic_hour_",
			Folder:  "/Users/indraoctama/Downloads/magic_hour/",
		},
		{
			BaseUrl: "https://raw.githubusercontent.com/indraoct/go-concurrency/main/files/ponytail",
			First:   1,
			LastInt: 151,
			Prefix:  "ponytail_",
			Folder:  "/Users/indraoctama/Downloads/ponytail/",
		},
	}

	for _, data := range dDownload {
		go iterateDownloadImages(data.First, data.LastInt, data.BaseUrl, data.Prefix, data.Folder, ch)
	}

	for i := 0; i < len(dDownload); i++ {
		retrieveData = <-ch //blocking the chanel processing
		fmt.Println(retrieveData)
	}
	fmt.Println("All processes completed")

}

func iterateDownloadImages(firstIt, lastIt int, baseUrl, prefix, folder string, ch chan string) {
	var data string

	fmt.Println("Process ", baseUrl, " start ", time.Now().Format(time.DateTime))

	defer func() {
		fmt.Println("Process ", data, " end ", time.Now().Format(time.DateTime))
		ch <- data
	}()

	count := 1
	for i := firstIt; i <= lastIt; i++ {
		iteration := prefix + fmt.Sprintf("%04d", i)
		url := baseUrl + "/" + iteration + ".png"
		result, err := downloadImages(url, folder, iteration)
		if err != nil {
			data = err.Error()
		} else {
			data = result + ", count data " + folder + " : " + strconv.Itoa(count) + " success"
		}

		fmt.Println(data)

		count++
	}

}

func downloadImages(url, folder, iterate string) (result string, err error) {
	// The URL of the image
	imageURL := url

	// The path where the image will be saved
	filePath := folder + iterate + ".png"

	// Create the directory if it doesn't exist
	if err = os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		err = fmt.Errorf("Error creating directory:", err)
		return
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("Error: %v creating file:", err)
		return
	}
	defer file.Close()

	// Get the image data
	resp, err := http.Get(imageURL)
	if err != nil {
		err = fmt.Errorf("Error: %v downloading image:", err)
		return
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("Response Status: %v Error: unable to download image, status code", resp.StatusCode)
		return
	}

	// Copy the data from the response to the file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		err = fmt.Errorf("Error: %v  saving image:", err)
		return
	}

	result = fmt.Sprintf("Image downloaded %v successfully!", filePath)
	return
}
