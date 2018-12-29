//go:generate x86_64-w64-mingw32-gcc -c -o hi_amd64.o hi.c
//go:generate x86_64-w64-mingw32-gcc -shared -o hi_amd64.dll -Wl,--no-insert-timestamp hi_amd64.o
//go:generate file2byteslice -input hi_amd64.dll -output hidll_amd64.go -package main -var hiDLL
//go:generate rm hi_amd64.dll hi_amd64.o

//go:generate i686-w64-mingw32-gcc -c -o hi_386.o hi.c
//go:generate i686-w64-mingw32-gcc -shared -o hi_386.dll -Wl,--no-insert-timestamp hi_386.o
//go:generate file2byteslice -input hi_386.dll -output hidll_386.go -package main -var hiDLL
//go:generate rm hi_386.dll hi_386.o

//go:generate gofmt -s -w .

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"golang.org/x/sys/windows"
)

var (
	dllfn  string
	modHi  *windows.LazyDLL
	procHi *windows.LazyProc
)

func initialize() error {
	const FILE_FLAG_DELETE_ON_CLOSE = 0x04000000

	dir, err := ioutil.TempDir("", "embeddll")
	if err != nil {
		return err
	}
	dllfn = filepath.Join(dir, "hi.dll")
	if err := ioutil.WriteFile(dllfn, hiDLL, 0777); err != nil {
		return err
	}

	modHi = windows.NewLazyDLL(dllfn)
	procHi = modHi.NewProc("hi")
	return nil
}

func terminate() error {
	if err := windows.FreeLibrary(windows.Handle(modHi.Handle())); err != nil {
		return err
	}
	if err := os.Remove(dllfn); err != nil {
		return err
	}
	return nil
}

func hi(f func() uintptr) (int, error) {
	r, _, err := procHi.Call(windows.NewCallback(f))
	if err != nil && err.(windows.Errno) != 0 {
		return 0, err
	}
	return int(r), nil
}

func main() {
	if err := initialize(); err != nil {
		panic(err)
	}
	defer func() {
		if err := terminate(); err != nil {
			panic(err)
		}
	}()

	v, err := hi(func() uintptr {
		return 42
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(v)
}
