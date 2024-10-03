### Comparing Embedded Key/Value Stores
Compare some embedded key value stores.

### BuntDB
I've have quite a few surprising problems with buntdb. 
- The benchmark fails to GET a value that was set during a transaction
  when using SyncAlways.
- Panic from an unbounded index during development of Querator
- The secondary indexes are useless unless you are serializing to JSON
- I do not understand why the value storage is `string` instead of `[]byte` which
  is what any marshalling library is going to provide, including standard JSON encoders
- The entire btree is in memory, which means you have to be careful of the total database
  size, else you will run out of memory.
- The code is amazingly simple and easy to follow, very few dependencies.

### BadgerDB
- Has a ton of external dependencies, and the code is quite complex.
- Based on some of the latest db research, is very fast read/write

### BoltDB
- Is based upon Btree+ which loads/unloads pages from disk at will
- Uses a write ahead log
- Is very slow compared to other implementations (8-10ms for a Set() operation)

### Results
```
goos: darwin
goarch: amd64
pkg: github.com/thrawn01/db-benchmarks
cpu: VirtualApple @ 2.50GHz
BenchmarkDatabases
BenchmarkDatabases/BuntDB
BenchmarkDatabases/BuntDB/BuntDB-Set
BenchmarkDatabases/BuntDB/BuntDB-Set-10               387495          3571 ns/op
BenchmarkDatabases/BuntDB/BuntDB-Get
BenchmarkDatabases/BuntDB/BuntDB-Get-10              2280698           529.8 ns/op
BenchmarkDatabases/BadgerDB
BenchmarkDatabases/BadgerDB/BadgerDB-Set
BenchmarkDatabases/BadgerDB/BadgerDB-Set-10           136003          8731 ns/op
BenchmarkDatabases/BadgerDB/BadgerDB-Get
BenchmarkDatabases/BadgerDB/BadgerDB-Get-10           575146          1889 ns/op
BenchmarkDatabases/BBolt
BenchmarkDatabases/BBolt/BBolt-Set
BenchmarkDatabases/BBolt/BBolt-Set-10                    127      10030328 ns/op
BenchmarkDatabases/BBolt/BBolt-Get
BenchmarkDatabases/BBolt/BBolt-Get-10                2023285           581.7 ns/op
PASS
```
