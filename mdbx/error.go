package mdbx

/*
#include "mdbx.h"
*/
import "C"

import (
	"os"
	"syscall"
)

// OpError is an error returned by the C API.  Not all errors returned by
// mdbx-go have type OpError but typically they do.  The Errno field will
// either have type Errno or syscall.Errno.
type OpError struct {
	Op    string
	Errno error
}

// Error implements the error interface.
func (err *OpError) Error() string {
	return err.Op + ": " + err.Errno.Error()
}

// The most common error codes do not need to be handled explicity.  Errors can
// be checked through helper functions IsNotFound, IsMapFull, etc, Otherwise
// they should be checked using the IsErrno function instead of direct
// comparison because they will typically be wrapped with an OpError.
const (
	// Error codes defined by MDBX.  See the list of MDBX return codes for more
	// information about each
	//
	//		http://symas.com/mdb/doc/group__errors.html

	KeyExist        Errno = C.MDBX_KEYEXIST
	NotFound        Errno = C.MDBX_NOTFOUND
	PageNotFound    Errno = C.MDBX_PAGE_NOTFOUND
	Corrupted       Errno = C.MDBX_CORRUPTED
	Panic           Errno = C.MDBX_PANIC
	VersionMismatch Errno = C.MDBX_VERSION_MISMATCH
	Invalid         Errno = C.MDBX_INVALID
	MapFull         Errno = C.MDBX_MAP_FULL
	DBsFull         Errno = C.MDBX_DBS_FULL
	ReadersFull     Errno = C.MDBX_READERS_FULL
	TxnFull         Errno = C.MDBX_TXN_FULL
	CursorFull      Errno = C.MDBX_CURSOR_FULL
	PageFull        Errno = C.MDBX_PAGE_FULL
	Incompatible    Errno = C.MDBX_INCOMPATIBLE
	BadRSlot        Errno = C.MDBX_BAD_RSLOT
	BadTxn          Errno = C.MDBX_BAD_TXN
	BadValSize      Errno = C.MDBX_BAD_VALSIZE
	BadDBI          Errno = C.MDBX_BAD_DBI
)

// Errno is an error type that represents the (unique) errno values defined by
// LMDB.  Other errno values (such as EINVAL) are represented with type
// syscall.Errno.  On Windows, LMDB return codes are translated into portable
// syscall.Errno constants (e.g. syscall.EINVAL, syscall.EACCES, etc.).
//
// Most often helper functions such as IsNotFound may be used instead of
// dealing with Errno values directly.
//
//		lmdb.IsNotFound(err)
//		lmdb.IsErrno(err, lmdb.TxnFull)
//		lmdb.IsErrnoSys(err, syscall.EINVAL)
//		lmdb.IsErrnoFn(err, os.IsPermission)
type Errno C.int

// minimum and maximum values produced for the Errno type. syscall.Errnos of
// other values may still be produced.
const minErrno, maxErrno C.int = C.MDBX_KEYEXIST, 22 //C.MDBX_LAST_ERRCODE

func (e Errno) Error() string {
	return C.GoString(C.mdbx_strerror(C.int(e)))
}

// _operrno is for use by tests that can't import C
func _operrno(op string, ret int) error {
	return operrno(op, C.int(ret))
}

// IsNotFound returns true if the key requested in Txn.Get or Cursor.Get does
// not exist or if the Cursor reached the end of the database without locating
// a value (EOF).
func IsNotFound(err error) bool {
	return IsErrno(err, NotFound)
}

// IsNotExist returns true the path passed to the Env.Open method does not
// exist.
func IsNotExist(err error) bool {
	return IsErrnoFn(err, os.IsNotExist)
}

// IsMapFull returns true if the environment map size has been reached.
func IsMapFull(err error) bool {
	return IsErrno(err, MapFull)
}

// IsErrno returns true if err's errno is the given errno.
func IsErrno(err error, errno Errno) bool {
	return IsErrnoFn(err, func(err error) bool { return err == errno })
}

// IsErrnoSys returns true if err's errno is the given errno.
func IsErrnoSys(err error, errno syscall.Errno) bool {
	return IsErrnoFn(err, func(err error) bool { return err == errno })
}

// IsErrnoFn calls fn on the error underlying err and returns the result.  If
// err is an *OpError then err.Errno is passed to fn.  Otherwise err is passed
// directly to fn.
func IsErrnoFn(err error, fn func(error) bool) bool {
	if err == nil {
		return false
	}
	if err, ok := err.(*OpError); ok {
		return fn(err.Errno)
	}
	return fn(err)
}

func operrno(op string, ret C.int) error {
	if ret == C.MDBX_SUCCESS {
		return nil
	}
	if minErrno <= ret && ret <= maxErrno {
		return &OpError{Op: op, Errno: Errno(ret)}
	}
	return &OpError{Op: op, Errno: syscall.Errno(ret)}
}
