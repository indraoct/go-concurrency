# Go Concurrency Example: Download JKT48 Member Photos
 This repository demonstrates the use of concurrency in Go to download images of JKT48 members. The code is structured to efficiently download images from specified URLs concurrently, making use of goroutines and wait groups and chanel.

## Table of Contents
- Overview
- Prerequisites
- Installation
- Usage
- Code Explanation
- License
- Overview
- The main objective of this project is to showcase how Go's concurrency model can be leveraged to download multiple images simultaneously. This example downloads images from three different categories of JKT48 members.

## Prerequisites
- Go installed on your local machine (version 1.20 or higher)
- Internet connection to download the images
- Installation
Clone this repository:
```
git clone https://github.com/indraoct/go-concurrency.git
```
Navigate to the project directory:
```
cd go-concurrency/example_1
```
Usage
To run the code, simply execute the following command:

```
go run main.go
```
This will start the concurrent download process for the images.

## Code Explanation
### Structure
**dataDownload**: A struct to hold information about each download task, including the base URL, range of images, prefix, and destination folder.
**main()**: The entry point of the application where download tasks are defined and initiated.
**iterateDownloadImages()**: A function to iterate over a range of image indices, constructing URLs, and calling the download function.
**downloadImages()**: A function to handle the actual downloading and saving of images to the specified folder.
Detailed Code Explanation
Data Structure
```go
type dataDownload struct {
BaseUrl string
First   int
LastInt int
Prefix  string
Folder  string
}
```

This struct holds the necessary details for downloading a series of images:

**BaseUrl**: The base URL where images are hosted.
**First**: The starting index for the images.
**LastInt**: The ending index for the images.
**Prefix**: The prefix for the image filenames.
**Folder**: The local folder where the images will be saved.

Main Function

```go
func main() {
var dDownload []dataDownload

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
        wg.Add(1)
        go iterateDownloadImages(data.First, data.LastInt, data.BaseUrl, data.Prefix, data.Folder, &wg)
    }

    wg.Wait()

    fmt.Println("All processes completed")
}
```

The main function initializes the download tasks and starts them concurrently using goroutines. The sync.WaitGroup is used to wait for all goroutines to finish.

Iterate and Download Functions

```go
func iterateDownloadImages(firstIt, lastIt int, baseUrl, prefix, folder string, wg *sync.WaitGroup) {
defer wg.Done()

    fmt.Printf("Process %s started\n", baseUrl)

    for i := firstIt; i <= lastIt; i++ {
        iteration := prefix + fmt.Sprintf("%04d", i)
        url := baseUrl + "/" + iteration + ".png"
        result, err := downloadImages(url, folder, iteration)
        if err != nil {
            fmt.Println(err)
        } else {
            fmt.Println(result)
        }
    }

    fmt.Printf("Process %s ended\n", baseUrl)
}

func downloadImages(url, folder, iterate string) (string, error) {
filePath := filepath.Join(folder, iterate+".png")

    if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
        return "", fmt.Errorf("error creating directory: %v", err)
    }

    file, err := os.Create(filePath)
    if err != nil {
        return "", fmt.Errorf("error creating file: %v", err)
    }
    defer file.Close()

    resp, err := http.Get(url)
    if err != nil {
        return "", fmt.Errorf("error downloading image: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("server returned non-200 status code: %d", resp.StatusCode)
    }

    if _, err := io.Copy(file, resp.Body); err != nil {
        return "", fmt.Errorf("error saving image: %v", err)
    }

    return fmt.Sprintf("Image downloaded successfully: %s", filePath), nil
}
```

**iterateDownloadImages**: Iterates through the image indices, constructs URLs, and calls downloadImages to download each image.
**downloadImages**: Handles the downloading of an image from a URL and saves it to the specified folder. It ensures directories are created if they don't exist and handles any errors that occur during the process.


## Other Example

I'm using chanel to implement download member JKT48 photos : 

- The main code:

```go
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
	
    fmt.Println("Semua proses berakhir")

```

The main function initializes the download tasks and starts them concurrently using goroutines. The chanel is used to block for all goroutines to finish.


- download with iteration code :
```go
func iterateDownloadImages(firstIt, lastIt int, baseUrl, prefix, folder string, ch chan string) {
	var data string

	fmt.Println("Process ", baseUrl, " start")

	defer func() {
		fmt.Println("Process ", data, " end")
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

```

look at the **defer func** part. There is a variable **ch** that receive data from the download logic (string) to state that
the process for 1 folder download is finish. So the part **retrieveData = <-ch** on the main function will release the blocking.

## License
This project is licensed under the MIT License.