// +build ignore

#include <stdlib.h>

typedef uintptr_t (*callback)();

__declspec(dllexport) int hi(callback cb) {
  return cb();
}
