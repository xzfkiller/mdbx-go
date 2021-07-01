/* mdbxgo.h
 * Helper utilities for github.com/xzfkiller/mdbx-go/mdbx.  These functions have
 * no compatibility guarantees and may be modified or deleted without warning.
 * */
#ifndef _MDBXGO_H_
#define _MDBXGO_H_

#include "mdbx.h"

/* Proxy functions for mdbx get/put operations. The functions are defined to
 * take char* values instead of void* to keep cgo from cheking their data for
 * nested pointers and causing a couple of allocations per argument.
 *
 * For more information see github issues for more information about the
 * problem and the decision.
 *      https://github.com/golang/go/issues/14387
 *      https://github.com/golang/go/issues/15048
 *      https://github.com/bmatsuo/lmdb-go/issues/63
 * */
int mdbxgo_mdb_del(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn);
int mdbxgo_mdb_get(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, MDBX_val *val);
int mdbxgo_mdb_put1(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, MDBX_val *val, unsigned int flags);
int mdbxgo_mdb_put2(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags);

/* ConstCString wraps a null-terminated (const char *) because Go's type system
 * does not represent the 'const' qualifier directly on a function argument and
 * causes warnings to be emitted during linking.
 * */
typedef struct{ const char *p; } mdbxgo_ConstCString;

#endif
