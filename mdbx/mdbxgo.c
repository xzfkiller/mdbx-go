/* mdbxgo.c
 * Helper utilities for github.com/xzfkiller/mdbx-go/mdbx
 * */
#include "mdbx.h"
#include "mdbxgo.h"
#include "_cgo_export.h"

#define MDBXGO_SET_VAL(val, size, data) \
    *(val) = (MDBX_val){.iov_len = (size), .iov_base = (data)}

int mdbxgo_mdb_del(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn) {
    MDBX_val key, val;
    MDBXGO_SET_VAL(&key, kn, kdata);
    MDBXGO_SET_VAL(&val, vn, vdata);
    return mdbx_del(txn, dbi, &key, &val);
}

int mdbxgo_mdb_get(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, MDBX_val *val) {
    MDBX_val key;
    MDBXGO_SET_VAL(&key, kn, kdata);
    return mdbx_get(txn, dbi, &key, val);
}

int mdbxgo_mdb_put2(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, char *vdata, size_t vn, unsigned int flags) {
    MDBX_val key, val;
    MDBXGO_SET_VAL(&key, kn, kdata);
    MDBXGO_SET_VAL(&val, vn, vdata);
    return mdbx_put(txn, dbi, &key, &val, flags);
}

int mdbxgo_mdb_put1(MDBX_txn *txn, MDBX_dbi dbi, char *kdata, size_t kn, MDBX_val *val, unsigned int flags) {
    MDBX_val key;
    MDBXGO_SET_VAL(&key, kn, kdata);
    return mdbx_put(txn, dbi, &key, val, flags);
}

