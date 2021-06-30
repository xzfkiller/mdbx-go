# mdbx-go

Go bindings to the libmdbx: https://erthink.github.io/libmdbx/

Most of articles in internet about LMDB are applicable to MDBX. But mdbx has more features.

For deeper DB understanding please read through [mdbx.h](https://github.com/erthink/libmdbx/blob/master/mdbx.h)

## Min Requirements

C language Compilers compatible with GCC or CLANG (mingw 10 on windows)
Golang: 1.14

## Packages

```go
import "github.com/xzfkiller/mdbx-go/mdbx"
```

Core bindings allowing low-level access to MDBX.

See make test for more information.


Use NoReadahead if Data > RAM

### Advantages of BoltDB

- Nested databases allow for hierarchical data organization.

- Far more databases can be accessed concurrently.

- No `Bucket` object - means less allocations and higher performance

- Operating systems that do not support sparse files do not use up excessive space due to a large pre-allocation of file
  space.

- As a pure Go package bolt can be easily cross-compiled using the `go`
  toolchain and `GOOS`/`GOARCH` variables.

- Its simpler design and implementation in pure Go mean it is free of many caveats and gotchas which are present using
  the MDBX package. For more information about caveats with the MDBX package, consult its
  [documentation](https://erthink.github.io/libmdbx/) so they are aware of these caveats. And even better if read
  through [mdbx.h](https://github.com/erthink/libmdbx/blob/master/mdbx.h).

### Advantages of MDBX

- Keys can contain multiple values using the DupSort flag.

- Updates can have sub-updates for atomic batching of changes.

- Databases typically remain open for the application lifetime. This limits the number of concurrently accessible
  databases. But, this minimizes the overhead of database accesses and typically produces cleaner code than an
  equivalent BoltDB implementation.

- Significantly faster than BoltDB. The raw speed of MDBX easily surpasses BoltDB. Additionally, MDBX provides
  optimizations ranging from safe, feature-specific optimizations to generally unsafe, extremely situational ones.
  Applications are free to enable any optimizations that fit their data, access, and reliability models.

- MDBX allows multiple applications to access a database simultaneously. Updates from concurrent processes are
  synchronized using a database lock file.

- As a C library, applications in any language can interact with MDBX databases. Mission critical Go applications can
  use a database while Python scripts perform analysis on the side.
