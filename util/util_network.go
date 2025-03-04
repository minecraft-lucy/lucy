/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package util

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"

	"lucy/lucyerrors"
	"lucy/tools"
)

// DownloadFile
//
// All downloaded files are stored in .lucy/downloads/{subdir}/{filename}
// Current policy for path is the slug of the package
func DownloadFile(
url string,
subdir string,
filename string,
) (out *os.File, err error) {
	if _, err := os.Stat(ProgramPath); os.IsNotExist(err) {
		return nil, lucyerrors.NoLucyError
	}

	out, err = os.Create(path.Join(DownloadPath, subdir, filename))
	if os.IsNotExist(err) {
		os.MkdirAll(path.Join(DownloadPath, subdir), os.ModePerm)
		out, _ = os.Create(path.Join(DownloadPath, subdir, filename))
	}
	defer out.Close()

	res, err := http.Get(url)
	defer res.Body.Close()

	fmt.Println("Downloading", url)

	// TODO: Move to output
	termWidth, _, _ := term.GetSize(int(os.Stdout.Fd()))
	bar := progressbar.NewOptions64(
		res.ContentLength,
		progressbar.OptionShowCount(),
		progressbar.OptionShowElapsedTimeOnFinish(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(
			progressbar.Theme{
				Saucer:        "[bold][magenta]█[reset]",
				SaucerHead:    "[bold][magenta]█[reset]",
				SaucerPadding: " ",
				BarStart:      "[bold][ [reset]",
				BarEnd:        "[bold] ][reset]",
			},
		),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetWidth(
			tools.TernaryFunc(
				func() bool { return termWidth/3 > 40 },
				termWidth/3,
				40,
			),
		),
	)
	writer := io.MultiWriter(out, bar)
	io.Copy(writer, res.Body)
	fmt.Println()

	return out, err
}

// MultiSourceDownload expects the urls hosts the same file. However, it does
// not verify the checksums to allow more loose file recognition policies in its
// callers.
//
// Download is concurrent. Other threads will be cancelled when one thread
// complete downloaded winThreshold of the file.
//
// Note that if the urls' speed are close, urls[0] will be selected since its
// goroutine is started first.
//
// Pros:
//   - Guaranteed to download the file from the fastest source.
//
// Cons:
//   - Wastes bandwidth
func MultiSourceDownload(urls []string, path string) {
	const winThreshold = 0.2 // 20% of the file
	var wg sync.WaitGroup
	var mu sync.Mutex
	var win bool
	var data *[]byte
	var winUrl string

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Get(url)
			if err != nil {
				return
			}
			defer resp.Body.Close()

			// TODO: totalSize might be -1 when size it not known, handle this case
			totalSize := resp.ContentLength
			thresholdSize := int64(float64(totalSize) * winThreshold)
			buffer := make([]byte, 2048)
			var downloadedSize int64

			for {
				n, err := resp.Body.Read(buffer)
				if err != nil && err != io.EOF {
					return
				}
				if n == 0 {
					break
				}
				downloadedSize += int64(n)
				if win && winUrl != url {
					println(
						"canceling:",
						url,
						"("+strconv.FormatInt(downloadedSize, 10)+"/"+
						strconv.FormatInt(totalSize, 10), "bytes)",
					)
					return
				}
				if downloadedSize >= thresholdSize {
					mu.Lock()
					if !win {
						println(
							"winning:",
							url,
							"("+strconv.FormatInt(downloadedSize, 10)+"/"+
							strconv.FormatInt(totalSize, 10), "bytes)",
						)
						win = true
						data = &buffer
						winUrl = url
					}
					mu.Unlock()
				}
			}
		}(url)
	}

	wg.Wait()
	println("winning url: ", winUrl)

	file, _ := os.Create(path)
	defer file.Close()
	_, err := file.Write(*data)
	if err != nil {
		panic(err)
	}

	println("Downloaded to", path)
}
