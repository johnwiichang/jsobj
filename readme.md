# JSObj - JavaScript Object Analyser

JSObj is a Golang library suitable for parsing the JavaScript object for `map` or `array` / `slice` object.

JSObj 是一个适用于解析 JavaScript 对象为 `map` 或 `array` / `slice` 对象的 Golang 库。

---

## 开始使用 <small>QuickStart</small>

```bash
go get "github.com/johnwiichang/jsobj"
```

## 例子 <small>Sample</small>

```go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/johnwiichang/jsobj"
)

func main() {
	js := `{
		$or: [
		   {key1: value1}, {key2:value2}
		]
	 }`
	obj, err := jsobj.Parse(js)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		bin, _ := json.Marshal(obj)
		fmt.Println(string(bin))
	}
}
```

### 输出 <small>Output</small>

```json
{"$or":[{"key1":"value1"},{"key2":"value2"}]}
```

> Some non-standard JavaScript objects may be supported.
> 
> 可能支持部分非标准的 JavaScript 对象表达。