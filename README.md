# funcstat

funcstat outputs function stat.

```
$ go install github.com/gostaticanalysis/funcstat/cmd/funcstat@latest
$ funcstat encoding/json | head
package,file,line,name,lines,bytes,params,returns,cyclomatic complexity
encoding/json,decode.go,96,Unmarshal,87,3970,2,1,1
encoding/json,decode.go,132,Error,6,309,0,1,1
encoding/json,decode.go,149,Error,3,191,0,1,1
encoding/json,decode.go,159,Error,10,263,0,1,1
encoding/json,decode.go,170,unmarshal,15,320,1,1,1
encoding/json,decode.go,191,String,2,102,0,1,1
encoding/json,decode.go,194,Float64,4,132,0,2,1
encoding/json,decode.go,199,Int64,4,127,0,2,1
encoding/json,decode.go,222,readIndex,4,116,0,1,1
```
