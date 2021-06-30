package mdbx

/*
#cgo CFLAGS: -pthread -W -Wall -Wno-unused-parameter -Wno-format-extra-args -Wbad-function-cast -Wno-missing-field-initializers -O2 -g
#cgo linux,pwritev CFLAGS: -DMDB_USE_PWRITEV
#cgo LDFLAGS: -L../lib -lmdbx -Wl,-rpath=../lib

#include "mdbx.h"
*/
import "C"

func cbool(b bool) C.bool {
	if b {
		return C.true
	}
	return C.false
}
