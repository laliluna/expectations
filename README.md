# Expectation library for Go language

A small library to validate data in tests and print expressive fail messages.

The API is heavily inspired by http://joel-costigliola.github.io/assertj/

![alt "Travis Build Status"](https://travis-ci.org/laliluna/expectations.svg?branch=master)

## Installation

```
go get github.com/laliluna/expectations
```

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
## Release notes

Since 0.6 

	eT.ExpectThatString("Hello World")

is preferred over		

	eT.ExpectThat("Hello World").String()

and 

	eT.ExpectThatSlice(numbers)

over

	eT.ExpectThat(numbers).Slice()


## Usage

```go
import (
	"testing"

	"github.com/laliluna/expectations"
)

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

	// String
	eT.ExpectThatString("Hello World").Equals("Hello World")
	eT.ExpectThatString("Hello World").EqualsIgnoringCase("hello world")
	eT.ExpectThatString("Hello World").DoesNotEqual("Bye World")
	eT.ExpectThatString("Hello World").Contains("Hello")
	eT.ExpectThatString("Hello World").StartsWith("Hello")
	eT.ExpectThatString("Hello World").EndsWith("World")
	eT.ExpectThatString("Hello World").DoesNotContain("John", "Doe")
	eT.ExpectThatString("Hello World").IsNotNil()

	// Slices and arrays
	numbers := []float32{1.1, 2.2, 3.3}
	eT.ExpectThatSlice(numbers).Contains(float32(1.1), float32(3.3))
	eT.ExpectThatSlice(numbers).DoesNotContain(float64(1.1), float32(1.22), float32(3.22))
	eT.ExpectThatSlice(numbers).IsNotEmpty() // IsEmpty

	eT.ExpectThatSlice(numbers).HasSize(3).First().Equals(float32(1.1)) // Second | Third | Nth

	numberArray := [3]float32{1.1, 2.2, 3.3}
	eT.ExpectThatSlice(numberArray).Contains(float32(1.1))
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
