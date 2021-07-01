package mdbx

/*
#include <stdlib.h>
#include <stdio.h>
#include "mdbx.h"
*/
import "C"

import (
	"errors"
	"runtime"
	"sync"
	"unsafe"
)

const success = C.MDBX_SUCCESS

const (
	// Flags for Env.Open.
	//
	// See mdbx_env_open

	NoSubdir    = C.MDBX_NOSUBDIR   // Argument to Open is a file, not a directory.
	Readonly    = C.MDBX_RDONLY     // Used in several functions to denote an object as readonly.
	WriteMap    = C.MDBX_WRITEMAP   // Use a writable memory map.
	NoMetaSync  = C.MDBX_NOMETASYNC // Don't fsync metapage after commit.
	MapAsync    = C.MDBX_MAPASYNC   // Flush asynchronously when using the WriteMap flag.
	NoTLS       = C.MDBX_NOTLS      // Danger zone. When unset reader locktable slots are tied to their thread.
	NoReadahead = C.MDBX_NORDAHEAD  // Disable readahead. Requires OS support.
	NoMemInit   = C.MDBX_NOMEMINIT  // Disable MDBX memory initialization.
)

// DBI is a handle for a database in an Env.
//
// See MDBX_dbi
type DBI C.MDBX_dbi

// Env is opaque structure for a database environment.  A DB environment
// supports multiple databases, all residing in the same shared-memory map.
//
// See MDBX_env.
type Env struct {
	_env *C.MDBX_env

	// closeLock is used to allow the Txn finalizer to check if the Env has
	// been closed, so that it may know if it must abort.
	closeLock sync.RWMutex

	ckey *C.MDBX_val
	cval *C.MDBX_val
}

// NewEnv allocates and initializes a new Env.
//
// See mdbx_env_create.
func NewEnv() (*Env, error) {
	env := new(Env)
	ret := C.mdbx_env_create(&env._env)
	if ret != success {
		return nil, operrno("mdbx_env_create", ret)
	}
	env.ckey = (*C.MDBX_val)(C.malloc(C.size_t(unsafe.Sizeof(C.MDBX_val{}))))
	env.cval = (*C.MDBX_val)(C.malloc(C.size_t(unsafe.Sizeof(C.MDBX_val{}))))

	runtime.SetFinalizer(env, (*Env).Close)
	return env, nil
}

// Open an environment handle. If this function fails Close() must be called to
// discard the Env handle.
//
// See mdbx_env_open.
func (env *Env) Open(path string) error {
	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))
	ret := C.mdbx_env_open(env._env, cpath, C.MDBX_NOSUBDIR|C.MDBX_COALESCE|C.MDBX_LIFORECLAIM|C.MDBX_NORDAHEAD, 0664)
	return operrno("mdbx_env_open", ret)
}

var errNotOpen = errors.New("enivornment is not open")
var errNegSize = errors.New("negative size")

func (env *Env) close() bool {
	if env._env == nil {
		return false
	}

	env.closeLock.Lock()
	C.mdbx_env_close(env._env)
	env._env = nil
	env.closeLock.Unlock()

	C.free(unsafe.Pointer(env.ckey))
	C.free(unsafe.Pointer(env.cval))
	env.ckey = nil
	env.cval = nil
	return true
}

// Close shuts down the environment, releases the memory map, and clears the
// finalizer on env.
//
// See mdbx_env_close.
func (env *Env) Close() error {
	if env.close() {
		runtime.SetFinalizer(env, nil)
		return nil
	}
	return errors.New("environment is already closed")
}

func (env *Env) SetGeometry(size_lower, size_now, size_upper, growth_step, shrink_threshold, pagesize int) error {
	if size_upper < 0 || pagesize < 0 {
		return errNegSize
	}
	ret := C.mdbx_env_set_geometry(env._env,
		C.long(size_lower), C.long(size_now), C.long(size_upper),
		C.long(growth_step), C.long(shrink_threshold), C.long(pagesize))
	return operrno("mdbx_env_set_geometry", ret)
}

// Path returns the path argument passed to Open.  Path returns a non-nil error
// if env.Open() was not previously called.
//
// See mdbx_env_get_path.
func (env *Env) Path() (string, error) {
	var cpath *C.char
	ret := C.mdbx_env_get_path(env._env, &cpath)
	if ret != success {
		return "", operrno("mdbx_env_get_path", ret)
	}
	if cpath == nil {
		return "", errNotOpen
	}
	return C.GoString(cpath), nil
}

// SetMaxReaders sets the maximum number of reader slots in the environment.
//
// See mdbx_env_set_maxreaders.
func (env *Env) SetMaxReaders(size int) error {
	if size < 0 {
		return errNegSize
	}
	ret := C.mdbx_env_set_maxreaders(env._env, C.uint(size))
	return operrno("mdb_env_set_maxreaders", ret)
}

// MaxReaders returns the maximum number of reader slots for the environment.
//
// See mdbx_env_get_maxreaders.
func (env *Env) MaxReaders() (int, error) {
	var max C.uint
	ret := C.mdbx_env_get_maxreaders(env._env, &max)
	return int(max), operrno("mdbx_env_get_maxreaders", ret)
}

// SetMaxDBs sets the maximum number of named databases for the environment.
//
// See mdbx_env_set_maxdbs.
func (env *Env) SetMaxDBs(size int) error {
	if size < 0 {
		return errNegSize
	}
	ret := C.mdbx_env_set_maxdbs(env._env, C.MDBX_dbi(size))
	return operrno("mdbx_env_set_maxdbs", ret)
}

// BeginTxn is an unsafe, low-level method to initialize a new transaction on
// env.  The Txn returned by BeginTxn is unmanaged and must be terminated by
// calling either its Abort or Commit methods to ensure that its resources are
// released.
//
// BeginTxn does not call runtime.LockOSThread.  Unless the Readonly flag is
// passed goroutines must call runtime.LockOSThread before calling BeginTxn and
// the returned Txn must not have its methods called from another goroutine.
// Failure to meet these restrictions can have undefined results that may
// include deadlocking your application.
//
// Instead of calling BeginTxn users should prefer calling the View and Update
// methods, which assist in management of Txn objects and provide OS thread
// locking required for write transactions.
//
// A finalizer detects unreachable, live transactions and logs thems to
// standard error.  The transactions are aborted, but their presence should be
// interpreted as an application error which should be patched so transactions
// are terminated explicitly.  Unterminated transactions can adversly effect
// database performance and cause the database to grow until the map is full.
//
// See mdbx_txn_begin.
func (env *Env) BeginTxn(parent *Txn, flags uint) (*Txn, error) {
	txn, err := beginTxn(env, parent, flags)
	if txn != nil {
		runtime.SetFinalizer(txn, func(v interface{}) { v.(*Txn).finalize() })
	}
	return txn, err
}

// RunTxn creates a new Txn and calls fn with it as an argument.  Run commits
// the transaction if fn returns nil otherwise the transaction is aborted.
// Because RunTxn terminates the transaction goroutines should not retain
// references to it or its data after fn returns.
//
// RunTxn does not call runtime.LockOSThread.  Unless the Readonly flag is
// passed the calling goroutine should ensure it is locked to its thread and
// any goroutines started by fn must not call methods on the Txn object it is
// passed.
//
// See mdbx_txn_begin.
func (env *Env) RunTxn(flags uint, fn TxnOp) error {
	return env.run(false, flags, fn)
}

// View creates a readonly transaction with a consistent view of the
// environment and passes it to fn.  View terminates its transaction after fn
// returns.  Any error encountered by View is returned.
//
// Unlike with Update transactions, goroutines created by fn are free to call
// methods on the Txn passed to fn provided they are synchronized in their
// accesses (e.g. using a mutex or channel).
//
// Any call to Commit, Abort, Reset or Renew on a Txn created by View will
// panic.
func (env *Env) View(fn TxnOp) error {
	return env.run(false, Readonly, fn)
}

// Update calls fn with a writable transaction.  Update commits the transaction
// if fn returns a nil error otherwise Update aborts the transaction and
// returns the error.
//
// Update calls runtime.LockOSThread to lock the calling goroutine to its
// thread and until fn returns and the transaction has been terminated, at
// which point runtime.UnlockOSThread is called.  If the calling goroutine is
// already known to be locked to a thread, use UpdateLocked instead to avoid
// premature unlocking of the goroutine.
//
// Neither Update nor UpdateLocked cannot be called safely from a goroutine
// where it isn't known if runtime.LockOSThread has been called.  In such
// situations writes must either be done in a newly created goroutine which can
// be safely locked, or through a worker goroutine that accepts updates to
// apply and delivers transaction results using channels.  See the package
// documentation and examples for more details.
//
// Goroutines created by the operation fn must not use methods on the Txn
// object that fn is passed.  Doing so would have undefined and unpredictable
// results for your program (likely including data loss, deadlock, etc).
//
// Any call to Commit, Abort, Reset or Renew on a Txn created by Update will
// panic.
func (env *Env) Update(fn TxnOp) error {
	return env.run(true, 0, fn)
}

// UpdateLocked behaves like Update but does not lock the calling goroutine to
// its thread.  UpdateLocked should be used if the calling goroutine is already
// locked to its thread for another purpose.
//
// Neither Update nor UpdateLocked cannot be called safely from a goroutine
// where it isn't known if runtime.LockOSThread has been called.  In such
// situations writes must either be done in a newly created goroutine which can
// be safely locked, or through a worker goroutine that accepts updates to
// apply and delivers transaction results using channels.  See the package
// documentation and examples for more details.
//
// Goroutines created by the operation fn must not use methods on the Txn
// object that fn is passed.  Doing so would have undefined and unpredictable
// results for your program (likely including data loss, deadlock, etc).
//
// Any call to Commit, Abort, Reset or Renew on a Txn created by UpdateLocked
// will panic.
func (env *Env) UpdateLocked(fn TxnOp) error {
	return env.run(false, 0, fn)
}

func (env *Env) run(lock bool, flags uint, fn TxnOp) error {
	if lock {
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()
	}
	txn, err := beginTxn(env, nil, flags)
	if err != nil {
		return err
	}
	return txn.runOpTerm(fn)
}

// CloseDBI closes the database handle, db.  Normally calling CloseDBI
// explicitly is not necessary.
//
// It is the caller's responsibility to serialize calls to CloseDBI.
//
// See mdbx_dbi_close.
func (env *Env) CloseDBI(db DBI) {
	C.mdbx_dbi_close(env._env, C.MDBX_dbi(db))
}
