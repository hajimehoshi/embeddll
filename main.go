//go:generate x86_64-w64-mingw32-gcc -c -o hi.o hi.c
//go:generate x86_64-w64-mingw32-gcc -shared -o hi.dll hi.o
//go:generate file2byteslice -input hi.dll -output hidll.go -package main -var hiDLL
//go:generate gofmt -s -w .
//go:generate rm hi.dll hi.o

package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var (
	modHi  *windows.LazyDLL
	procHi *windows.LazyProc
)

func init() {
	dir, err := ioutil.TempDir("", "windlltest")
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
