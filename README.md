# Expectation library for Go language

A small library to validate data in tests and print expressive fail messages.

The API is heavily inspired by http://joel-costigliola.github.io/assertj/

![alt "Travis Build Status"](https://travis-ci.org/laliluna/expectations.svg?branch=master)

## Introduction

```go
import (
	"testing"
	"github.com/laliluna/expectations"
)

func TestDemo(t *testing.T) {
  eT := expectations.NewT(t)
  eT.ExpectThat(5).IsGreater(6)
  eT.ExpectThat("Hello World").String().EndsWith("joe")

  values := []int{1, 2, 3}
  eT.ExpectThat(values).Slice().Contains(1, 2, 55, 66)
}
```

```
expectations_test.go
--------------------
--- TestDemo in line 15: Expect 5 to be greater than 6
--- TestDemo in line 16: Expect Hello World to end with joe
--- TestDemo in line 18: Expect [1 2 3] to contain [1 2 55 66] but was missing [55 66]

--- FAIL: TestDemo (0.00s)
```

## Usage

```go
func TestDemo(t *testing.T) {
 	eT := expectations.NewT(t)

	eT.ExpectThat(5).Equals(5)
	eT.ExpectThat(5).IsGreater(4)
	eT.ExpectThat(5).DoesNotEqual(1)
	eT.ExpectThat(5).IsGreater(4)
	eT.ExpectThat(5).IsGreaterOrEqual(4)
	eT.ExpectThat(5).IsLower(7)
	eT.ExpectThat(5).IsLowerOrEqual(5)
	eT.ExpectThat(nil).IsNil()
	var foo interface{}
	foo = 5
	eT.ExpectThat(foo).IsNotNil()

	// Chaining
	eT.ExpectThat(5).IsGreater(2).IsLower(7)

	eT.ExpectThat("Hello World").String().Equals("Hello World")
	eT.ExpectThat("Hello World").String().EqualsIgnoringCase("hello world")
	eT.ExpectThat("Hello World").String().DoesNotEqual("Bye World")
	eT.ExpectThat("Hello World").String().Contains("Hello")
	eT.ExpectThat("Hello World").String().StartsWith("Hello")
	eT.ExpectThat("Hello World").String().EndsWith("World")
	eT.ExpectThat("Hello World").String().DoesNotContain("John", "Doe")
	eT.ExpectThat("Hello World").String().IsNotNil()

	numbers := []float32{1.1, 2.2, 3.3}
	eT.ExpectThat(numbers).Slice().Contains(float32(1.1), float32(3.3))
	eT.ExpectThat(numbers).Slice().DoesNotContain(float64(1.1), float32(1.22), float32(3.22))

	numberArray := [3]float32{1.1, 2.2, 3.3}
	eT.ExpectThat(numberArray).Slice().Contains(float32(1.1))
}
```

You can chain assertions.

```go
eT.ExpectThat(5).DoesNotEqual(1).IsGreater(4)
```

## Comparing different types

Different types will always fail.

```go
var apple, pear interface{}
apple = uint(5)
pear = int(5)
eT.ExpectThat(apple).DoesNotEqual(pear) // is not equal but we still reject the comparison
eT.ExpectThat(apple).Equals(pear) 
```
```
expectations_test.go
--------------------
--- TestDemo in line 15: You try to compare different types 5 (int) to 5 (uint)
--- TestDemo in line 15: You try to compare different types 5 (int) to 5 (uint)

```
# License

The code is published under [Apache License Version 2.0](LICENSE)
