package mdbx

/*
#cgo CFLAGS: -pthread -W -Wall -Wno-unused-parameter -Wno-format-extra-args -Wbad-function-cast -Wno-missing-field-initializers -O2 -g
#cgo linux,pwritev CFLAGS: -DMDB_USE_PWRITEV
#cgo LDFLAGS: -L. -lmdbx -Wl,-rpath=.
*/
import "C"
