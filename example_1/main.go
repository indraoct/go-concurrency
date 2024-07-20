package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
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
		dDownload []dataDownload
	)

	wg := sync.WaitGroup{}

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
		wg.Add(1) // adding waitGroup for every iteration
		go func() {
			iterateDownloadImages(data.First, data.LastInt, data.BaseUrl, data.Prefix, data.Folder)
			wg.Done() // end the wait group for each iteration
		}()
	}

	//waiting all the go routines wait group have been done
	wg.Wait()

	fmt.Println("All processes completed")

}

func iterateDownloadImages(firstIt, lastIt int, baseUrl, prefix, folder string) {

	fmt.Println("Process ", baseUrl, " start ", time.Now().Format(time.DateTime))
	defer func() {
		fmt.Println("Process ", baseUrl, " end ", time.Now().Format(time.DateTime))
	}()

	count := 1
	for i := firstIt; i <= lastIt; i++ {
		iteration := prefix + fmt.Sprintf("%04d", i)
		url := baseUrl + "/" + iteration + ".png"
		result, err := downloadImages(url, folder, iteration)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(result)
		}
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
