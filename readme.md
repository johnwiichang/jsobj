# JSObj - JavaScript Object Analyser

JSObj is a Golang library suitable for parsing the JavaScript object for `map` or `array` / `slice` object.

JSObj 是一个适用于解析 JavaScript 对象为 `map` 或 `array` / `slice` 对象的 Golang 库。

---

[![Build Status](https://travis-ci.org/johnwiichang/jsobj.svg?branch=main)](https://travis-ci.org/johnwiichang/jsobj) [![Go Report Card](https://goreportcard.com/badge/github.com/johnwiichang/jsobj)](https://goreportcard.com/report/github.com/johnwiichang/jsobj) [![Open Source Helpers](https://www.codetriage.com/johnwiichang/jsobj/badges/users.svg)](https://www.codetriage.com/johnwiichang/jsobj) 

## 开始使用 <small>QuickStart</small>

```bash
go get "github.com/johnwiichang/jsobj"
```

## 方法 <small>Methods</small>
Provides multiple methods for handling the JavaScript statement to complete the analysis of the JavaScript statement.

提供用于处理 JavaScript 语句的多个方法以完成 JavaScript 语句的解析。

### ReadObjects() ([]interface{}, error)
Read several JavaScript objects. Suitable parameters in the method: `(Obj1, obj2)`

读取数个 JavaScript 对象。适用于方法中的参数：`(obj1, obj2)`

### ReadObject() (interface{}, error)
Read a JavaScript object.

读取一个 JavaScript 对象。

### ReadWord() (Word, error)
Read a JavaScript word.

读取一个 JavaScript 字。

> JavaScript word may be a token, or it is possible to be a string.
> 
> JavaScript 字可能是一个 Token，也有可能是一个字符串。

### Read(...rune) (string, error)
Get a set of strings until a character.

获取一组字符串一直到某个字符。

### Location() int
Get the current character position.

获取当前字符位置。

### EOF() bool
Whether it has arrived at the end of the file.

是否已经抵达文件末尾。

> Some non-standard JavaScript objects may be supported.
> 
> 可能支持部分非标准的 JavaScript 对象表达。
