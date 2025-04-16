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

package remote

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"lucy/lucytypes"
)

var AvailableSources = map[lucytypes.Platform][]lucytypes.Source{
	lucytypes.Fabric: {lucytypes.CurseForge, lucytypes.Modrinth},
	lucytypes.Forge:  {lucytypes.CurseForge, lucytypes.Modrinth},
	lucytypes.Mcdr:   {lucytypes.McdrCatalogue},
}

var SpeedTestUrls = map[lucytypes.Source]string{
	lucytypes.CurseForge: "https://mediafilez.forgecdn.net/files/4834/896/fabric-api-0.87.2%2B1.19.4.jar",
	lucytypes.Modrinth:   "https://cdn.modrinth.com/data/P7dR8mSH/versions/nyAmoHlr/fabric-api-0.87.2%2B1.19.4.jar",
}

const slow float64 = 0x7FF0000000000000 // inf

// SelectSource is an alternative to MultiSourceDownload. It fetches a fixed url
// from SpeedTestUrls and measures the download speed of each source. The source
// with the fastest download speed is returned.
//
// Pros:
//   - Fastest source can be stored for later use.
//   - Saves bandwidth
//
// Cons:
//   - Speed test might not be representative
func SelectSource(platform lucytypes.Platform) lucytypes.Source {
	slowest := slow
	fastestSource := lucytypes.UnknownSource
	wg := sync.WaitGroup{}
	for _, source := range AvailableSources[platform] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			speed := testDownloadSpeed(SpeedTestUrls[source])
			if speed < slowest {
				fastestSource = source
			}
			fmt.Printf("Speed for %s: %f\n", source, speed)
		}()
	}

	wg.Wait()
	if fastestSource == lucytypes.UnknownSource {
		panic("No available source")
	}

	return fastestSource
}

func testDownloadSpeed(url string) (elapsedTime float64) {
	startTime := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		return slow
	}
	defer resp.Body.Close()

	chunkSize := 2048

	buffer := make([]byte, chunkSize)
	for i := 0; i < 10; i++ {
		n, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			return slow
		}
		if n == 0 {
			break
		}
	}

	elapsedTime = time.Since(startTime).Seconds()
	return elapsedTime
}
