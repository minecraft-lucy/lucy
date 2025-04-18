package util

import (
	"fmt"
	"lucy/logger"
	"lucy/tools"
	"os"
	"path"
	"slices"
	"time"
)

func init() {
	if err := os.MkdirAll(CacheDir, os.ModePerm); err != nil {
		CacheOn = false
		logger.Warn(
			fmt.Errorf(
				"failed to create cache directory, disabling caching: %w",
				err,
			),
		)
	}
	CacheOn = true
	clearExpiredCache()
	maintainCacheLimit()
}

var (
	CacheDir = path.Join(os.TempDir(), ProgramName)
	CacheOn  bool
)

const (
	CacheLiveTime = 30 * time.Minute
	CacheLimit    = 30 * 1024 * 1024 // 30MB
)

func Cache(filename string, content []byte) (err error) {
	if !CacheOn {
		return nil
	}
	filepath := path.Join(CacheDir, filename)
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer tools.CloseReader(file, logger.Warn)
	_, err = file.Write(content)
	if err != nil {
		return err
	}
	return nil
}

func GetCache(filename string) (hit bool, file *os.File, err error) {
	if !CacheOn {
		return false, nil, nil
	}
	filepath := path.Join(CacheDir, filename)
	file, err = os.Open(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, file, nil
}

func DeleteCache(filename string) (err error) {
	if !CacheOn {
		return nil
	}
	filepath := path.Join(CacheDir, filename)
	err = os.Remove(filepath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func ClearAllCache() (err error) {
	err = os.RemoveAll(CacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return nil
}

func GetCacheDir() string {
	return CacheDir
}

func DoCache() bool {
	return CacheOn
}

func clearExpiredCache() {
	files, err := os.ReadDir(CacheDir)
	if err != nil {
		logger.Warn(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := path.Join(CacheDir, file.Name())
		info, err := os.Stat(filepath)
		if err != nil {
			logger.Warn(err)
			continue
		}
		if time.Since(info.ModTime()) > CacheLiveTime {
			err = os.Remove(filepath)
			if err != nil {
				logger.Warn(err)
			}
		}
	}
}

func maintainCacheLimit() {
	files, err := os.ReadDir(CacheDir)
	if err != nil {
		logger.Warn(err)
		return
	}
	var size int64
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			logger.Warn(err)
			continue
		}
		size += info.Size()
	}
	if size > CacheLimit {
		type fileInfo struct {
			path    string
			modTime time.Time
			size    int64
		}

		fileInfos := make([]fileInfo, 0, len(files))
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			info, err := file.Info()
			if err != nil {
				logger.Warn(err)
				continue
			}
			fileInfos = append(
				fileInfos, fileInfo{
					path:    path.Join(CacheDir, file.Name()),
					modTime: info.ModTime(),
					size:    info.Size(),
				},
			)
		}

		slices.SortFunc(
			fileInfos, func(a, b fileInfo) int {
				return int(a.modTime.Sub(b.modTime))
			},
		)

		for _, info := range fileInfos {
			if size <= CacheLimit {
				break
			}
			err := os.Remove(info.path)
			if err != nil {
				logger.Warn(err)
				continue
			}
			size -= info.size
		}
	}
}
