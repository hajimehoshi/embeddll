//go:generate x86_64-w64-mingw32-gcc -shared -o hi_amd64.dll -Wl,--no-insert-timestamp hi.c
//go:generate file2byteslice -input hi_amd64.dll -output hidll_amd64.go -package main -var hiDLL
//go:generate rm hi_amd64.dll

//go:generate i686-w64-mingw32-gcc -shared -o hi_386.dll -Wl,--no-insert-timestamp hi.c
//go:generate file2byteslice -input hi_386.dll -output hidll_386.go -package main -var hiDLL
//go:generate rm hi_386.dll

//go:generate gofmt -s -w .

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows"
)

var (
	dllfn  string
	modHi  *windows.LazyDLL
	procHi *windows.LazyProc
)

func createTmpDLL(content []byte) (string, error) {
	f, err := ioutil.TempFile("", "hi.*.dll")
	if err != nil {
		return "", err
	}
	defer f.Close()

	fn := f.Name()

	if _, err := f.Write(content); err != nil {
		return "", err
	}

	return fn, nil
}

func initialize() error {
	fn, err := createTmpDLL(hiDLL)
	if err != nil {
		return err
	}
	dllfn = fn

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
