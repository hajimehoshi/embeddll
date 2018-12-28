//go:generate x86_64-w64-mingw32-gcc -c -o hi_amd64.o hi.c
//go:generate x86_64-w64-mingw32-gcc -shared -o hi_amd64.dll -Wl,--no-insert-timestamp hi_amd64.o
//go:generate i686-w64-mingw32-gcc -c -o hi_386.o hi.c
//go:generate i686-w64-mingw32-gcc -shared -o hi_386.dll -Wl,--no-insert-timestamp hi_386.o
//go:generate file2byteslice -input hi_amd64.dll -output hidll_amd64.go -package main -var hiDLL
//go:generate file2byteslice -input hi_386.dll -output hidll_386.go -package main -var hiDLL
//go:generate gofmt -s -w .
//go:generate rm hi_amd64.dll hi_386.dll hi_amd64.o hi_386.o

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"

	"golang.org/x/sys/windows"
)

var (
	modHi  *windows.LazyDLL
	procHi *windows.LazyProc
)

func init() {
	dir, err := ioutil.TempDir("", "embeddll")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "hi.dll"), hiDLL, 0777); err != nil {
		panic(err)
	}

	modHi = windows.NewLazyDLL(filepath.Join(dir, "hi"))
	procHi = modHi.NewProc("hi")
}

func hi() (int, error) {
	r, _, err := procHi.Call(0, 0, 0)
	if err != nil && err.(windows.Errno) != 0 {
		return 0, err
	}
	return int(r), nil
}

func main() {
	v, err := hi()
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
