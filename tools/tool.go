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

package tools

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"sync"
	"time"

	"github.com/russross/blackfriday/v2"
)

// TernaryFunc gives a if expr == true, b if expr == false. For a simple
// bool expression, use Ternary instead.
func TernaryFunc[T any](expr func() bool, a T, b T) T {
	if expr() {
		return a
	}
	return b
}

// Ternary gives a if v == true, b if v == false. For a function parameter, use
// TernaryFunc instead.
func Ternary[T any](v bool, a T, b T) T {
	if v {
		return a
	}
	return b
}

// Memoize is only used for functions that do not take any arguments and return
// a value (typically a struct) that can be treated as a constant.
func Memoize[T any](f func() T) func() T {
	var res T
	var once sync.Once
	return func() T {
		once.Do(
			func() {
				res = f()
			},
		)
		return res
	}
}

func MemoizeE[T any](f func() (T, error)) func() (T, error) {
	var res T
	var err error
	var once sync.Once
	return func() (T, error) {
		once.Do(
			func() {
				res, err = f()
			},
		)
		return res, err
	}
}

// Insert inserts a value into a slice at a slice[pos]. If the pos is out of
// bounds, the slice remains unchanged.
func Insert[T any](slice []T, pos int, value ...T) []T {
	if pos < 0 || pos > len(slice) {
		return slice
	}
	return append(slice[:pos], append(value, slice[pos:]...)...)
}

// CloseReader closes a reader and runs failAction() if error occurs. Call this
// with a defer statement.
func CloseReader(reader io.ReadCloser, failAction func(error)) {
	err := reader.Close()
	if err != nil {
		failAction(err)
	}
}

const (
	NetworkTestTimeout = 5 // seconds
	NetworkTestRetries = 3
)

// NetworkTest is a simple the network connection test. You can use this before
// any operation that strictly requires a network connection.
//
// A nil value means the connection is successful.
func NetworkTest() (err error) {
	retry := NetworkTestRetries
	client := http.Client{
		Timeout: NetworkTestTimeout * time.Second,
	}
Retry:
	_, err = client.Get("https://example.com")
	if err != nil {
		retry--
		if retry > 0 {
			goto Retry
		}
		return err
	}
	return nil
}

func MarkdownToPlainText(md string) (s string) {
	s = string(blackfriday.Run([]byte(md)))
	return
}

// Decorate applies a series of decorators to a function. This is used to
// prevent nested function calls for better readability.
func Decorate[T interface{}](f T, decorators ...func(T) T) T {
	for _, decorator := range decorators {
		f = decorator(f)
	}
	return f
}

// UnderCd checks if the path is under the current working directory (non-recursive).
func UnderCd(path string) bool {
	abs, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	cd, err := os.Getwd()
	if err != nil {
		return false
	}

	parent := filepath.Dir(abs)
	return parent == cd
}

// KeyValue works together with SortAndExtract to sort a slice of Item
// with their corresponding Index.
type KeyValue[T, Ti any] struct {
	Item  T
	Index Ti
}

func SortAndExtract[T, Ti any](
	arr []KeyValue[T, Ti],
	cmp func(a, b KeyValue[T, Ti]) int,
) (res []T) {
	slices.SortFunc(arr, cmp)
	for _, item := range arr {
		res = append(res, item.Item)
	}
	return res
}
