package probe

import (
	"lucy/logger"
	"lucy/probe/internal/detector"
	"lucy/prompt"
	"lucy/tools"
	"lucy/types"
	"os"
	"path"
	"sync"
	"sync/atomic"
)

// getExecutableInfo uses the new detector-based architecture to find server executables
var getExecutableInfo = tools.Memoize(
	func() *types.ExecutableInfo {
		var valid []*types.ExecutableInfo
		workPath := getServerWorkPath()
		jars, err := findJar(workPath)
		if err != nil {
			logger.Warn(err)
			logger.Info("cannot read the current directory, most features will be disabled")
		}

		// Use the new detector-based approach
		for _, jar := range jars {
			exec := detector.Executable(jar)
			if exec == nil {
				continue
			}
			valid = append(valid, exec)
		}

		if len(valid) == 0 {
			logger.Info("no server jar found, trying to find under libraries")
			jarPaths := findJarRecursive(path.Join(workPath, "libraries"))
			if len(jarPaths) == 0 {
				// The following code is commented out due to the aggressive search
				// being too slow and inaccurate. It is kept here for future reference.
				//
				// logger.Info("still no server jar found, attempting even more aggressive search")
				// logger.Info("note that this may take a long time, and the accuracy is not guaranteed")
				// jarPaths = findJarRecursive(workPath)

				return nil
			}
			mu := sync.Mutex{}
			wg := sync.WaitGroup{}
			for _, jarPath := range jarPaths {
				wg.Add(1)
				go func(jarPath string) {
					exec := detector.Executable(jarPath)
					if exec == nil {
						wg.Done()
						return
					}
					mu.Lock()
					valid = append(valid, exec)
					mu.Unlock()
					wg.Done()
				}(jarPath)
			}
			wg.Wait()
		}

		if len(valid) == 0 {
			logger.Info("no server under current directory")
			return detector.UnknownExecutable
		} else if len(valid) == 1 {
			return valid[0]
		}

		var choice int
		noExecUnderCd := true
		for _, exec := range valid {
			if tools.UnderCd(exec.Path) {
				noExecUnderCd = false
			}
		}
		if noExecUnderCd {
			choice = prompt.SelectExecutable(
				valid,
				[]prompt.Note{prompt.NoteSuspectPrePackagedServer},
			)
		} else {
			choice = prompt.SelectExecutable(valid, nil)
		}
		return valid[choice]
	},
)

func findJar(dir string) (jarFiles []string, err error) {
	jarFiles = []string{}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if path.Ext(entry.Name()) == ".jar" {
			jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
		}
	}

	return jarFiles, nil
}

const fileCountThreshold = 50000

func findJarRecursive(dir string) (jarFiles []string) {
	jarFiles = []string{}
	entries, _ := os.ReadDir(dir)
	var wg sync.WaitGroup
	var fileCount int32
	var mu sync.Mutex

	// TODO: Use semaphore to limit the number of goroutines
	for _, entry := range entries {
		if atomic.LoadInt32(&fileCount) >= fileCountThreshold {
			logger.Info("file count threshold reached, stopping search")
			break
		}
		if entry.IsDir() {
			wg.Add(1)
			go func(subDir string) {
				defer wg.Done()
				subJarFiles := findJarRecursive(subDir)
				mu.Lock()
				jarFiles = append(jarFiles, subJarFiles...)
				mu.Unlock()
			}(path.Join(dir, entry.Name()))
		} else {
			atomic.AddInt32(&fileCount, 1)
			if path.Ext(entry.Name()) == ".jar" {
				mu.Lock()
				jarFiles = append(jarFiles, path.Join(dir, entry.Name()))
				mu.Unlock()
			}
		}
	}

	wg.Wait()
	return
}
