package tools

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/russross/blackfriday/v2"
)

const CRLF = "\r\n"

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

func TernaryLazy[T any](v bool, a func() T, b func() T) T {
	if v {
		return a()
	}
	return b()
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
	networkTestTimeout = 5 // seconds
	networkTestRetries = 3
)

// factoryNetworkTest is a simple the network connection test. You can use this before
// any operation that strictly requires a network connection.
//
// A nil value means the connection is successful.
func factoryNetworkTest(url string, retry int, timeout int) func() (err error) {
	return func() (err error) {
		retry := networkTestRetries
		client := http.Client{Timeout: networkTestTimeout * time.Second}
	Retry:
		_, err = client.Get(url)
		if err != nil {
			retry--
			if retry > 0 {
				goto Retry
			}
			return err
		}
		return nil
	}
}

var GoogleTest = factoryNetworkTest(
	"https://www.google.com",
	networkTestRetries,
	networkTestRetries,
)

var GithubTest = factoryNetworkTest(
	"https://github.com",
	networkTestRetries,
	networkTestRetries,
)

var RegularTest = factoryNetworkTest(
	"https://www.example.com",
	networkTestRetries,
	networkTestRetries,
)

func MarkdownToPlainText(md string) string {
	parser := blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
	doc := parser.Parse([]byte(md))

	var output string
	var listDepth, listItemNumber int
	var inCodeBlock bool
	var lastType blackfriday.NodeType

	preNodeHandler := func(node *blackfriday.Node) {
		switch node.Type {
		case blackfriday.Heading:
			if lastType != blackfriday.Document &&
				!strings.HasSuffix(output, "\n\n") &&
				node.HeadingData.Level != 1 {
				output += "\n\n"
			}
			for range node.HeadingData.Level {
				output += "█"
			}
			output += " "
		case blackfriday.List:
			listDepth++
			if node.ListData.ListFlags&blackfriday.ListTypeOrdered != 0 {
				listItemNumber = 1
			}
		case blackfriday.Item:
			prefix := strings.Repeat("  ", listDepth-1)
			if node.Parent.ListData.ListFlags&blackfriday.ListTypeOrdered != 0 {
				output += prefix + strconv.Itoa(listItemNumber) + ". "
				listItemNumber++
			} else {
				output += prefix + "• "
			}
		case blackfriday.CodeBlock:
			inCodeBlock = true
			output += "\n```\n"
		case blackfriday.Paragraph:
			if lastType != blackfriday.Document &&
				lastType != blackfriday.Item &&
				lastType != blackfriday.Heading &&
				!strings.HasSuffix(output, "\n\n") {
				output += "\n\n"
			}
		}
	}

	postNodeHandler := func(node *blackfriday.Node) {
		switch node.Type {
		case blackfriday.Heading:
			output += "\n\n"
		case blackfriday.List:
			listDepth--
			if listDepth == 0 && !strings.HasSuffix(output, "\n\n") {
				output += "\n"
			}
		case blackfriday.CodeBlock:
			inCodeBlock = false
			output += "\n```\n"
		case blackfriday.Paragraph:
			if !strings.HasSuffix(output, "\n") {
				output += "\n"
			}
		}
		lastType = node.Type
	}

	textHandler := func(node *blackfriday.Node, entering bool) {
		switch node.Type {
		case blackfriday.Text:
			// preserve what's in the code block
			if inCodeBlock {
				output += string(node.Literal)
			} else {
				output += string(node.Literal)
			}
		case blackfriday.Softbreak:
			// preserve linebreak in the code block
			if inCodeBlock {
				output += "\n"
			} else {
				output += " "
			}
		case blackfriday.Hardbreak:
			output += "\n"
		case blackfriday.Strong:
			if entering {
				node.LinkData.NoteID = len(output)
			} else {
				start := node.LinkData.NoteID
				boldText := output[start:]
				output = output[:start] + Bold(boldText)
			}
		case blackfriday.Emph:
			if entering {
				node.LinkData.NoteID = len(output)
			} else {
				start := node.LinkData.NoteID
				emphText := output[start:]
				output = output[:start] + Underline(emphText)
			}
		case blackfriday.Link:
			if !entering {
				// url
				output += " (" + string(node.LinkData.Destination) + ")"
			}
		case blackfriday.Image:
			if !entering {
				// omit image
				output += "[Image]"
			}
		case blackfriday.Code:
			output += "`" + string(node.Literal) + "`"
		}
	}

	doc.Walk(
		func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
			if entering {
				preNodeHandler(node)
			}

			textHandler(node, entering)

			if !entering {
				postNodeHandler(node)
			}

			return blackfriday.GoToNext
		},
	)

	// 清理多余的空行
	output = strings.ReplaceAll(output, "\n\n\n", "\n\n")
	return strings.TrimSpace(output)
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
