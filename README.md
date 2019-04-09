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
	eT.ExpectThat("Hello World").EndsWith("joe")

	values := []int{1, 2, 3}
	eT.ExpectThat(values).Slice().Contains(1, 2, 55, 66)
```

```
expectations_test.go
--------------------
--- TestDemo in line 15: Expect Hello World to end with joe
--- TestDemo in line 17: Expect [1 2 3] to contain [1 2 55 66] but was missing [55 66]

--- FAIL: TestDemo (0.00s)
```

## Usage

```go
func TestDemo(t *testing.T) {
  eT := expectations.NewT(t)
  eT.ExpectThat(5).Equals(5)
  eT.ExpectThat(5).DoesNotEqual(1)
  eT.ExpectThat(5).IsGreater(4)
  eT.ExpectThat(5).IsGreaterOrEqual(4)
  eT.ExpectThat(5).IsLower(7)
  eT.ExpectThat(5).IsLowerOrEqual(5)

  eT.ExpectThat(5).IsGreater(2).IsLower(7)

  eT.ExpectThat("Hello World").String().Equals("Hello World")
  eT.ExpectThat("Hello World").String().EqualsIgnoringCase("hello world")
  eT.ExpectThat("Hello World").String().DoesNotEqual("Bye World")
  eT.ExpectThat("Hello World").String().Contains("Hello")
  eT.ExpectThat("Hello World").String().StartsWith("Hello")
  eT.ExpectThat("Hello World").String().EndsWith("World")
  eT.ExpectThat("Hello World").String().DoesNotContain("John", "Doe")

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

# License

The code is published under [Apache License Version 2.0](LICENSE)
