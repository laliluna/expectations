# Expectation library for Go language

This is an initial version, to discuss the API. Expect API changes with the next release. Feel invited to try out the library.

The API is heavily inspired by http://joel-costigliola.github.io/assertj/

## Introduction

```go
import (
	"testing"
	"laliluna.de/expectations"
)

func TestDemo(t *testing.T) {
	eT := expectations.NewT(t)
	eT.ExpectString("Hello World").EndsWith("joe")

	values := []int{1, 2, 3}
	eT.ExpectSlice(values).Contains(1, 2, 55, 66)
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
  eT := expectations.NewT(t)eT.Expect(5).Equals(5)
  eT.Expect(5).DoesNotEqual(1)
  eT.Expect(5).ToBeGreater(4)
  eT.Expect(5).ToBeGreaterOrEqual(4)
  eT.Expect(5).ToBeLower(7)
  eT.Expect(5).ToBeLowerOrEqual(5)

  eT.Expect(5).ToBeGreater(2).ToBeLower(7)

  eT.ExpectString("Hello World").Equals("Hello World")
  eT.ExpectString("Hello World").EqualsIgnoringCase("hello world")
  eT.ExpectString("Hello World").DoesNotEqual("Bye World")
  eT.ExpectString("Hello World").Contains("Hello")
  eT.ExpectString("Hello World").StartsWith("Hello")
  eT.ExpectString("Hello World").EndsWith("World")
  eT.ExpectString("Hello World").DoesNotContain("John", "Doe")

  numbers := []float32{1.1, 2.2, 3.3}
  eT.ExpectSlice(numbers).Contains(float32(1.1), float32(3.3))
  eT.ExpectSlice(numbers).DoesNotContain(float64(1.1), float32(1.22), float32(3.22))

  numberArray := [3]float32{1.1, 2.2, 3.3}
  eT.ExpectSlice(numberArray).Contains(float32(1.1))
}
```

You can chain assertions.

```go
eT.Expect(5).DoesNotEqual(1).ToBeGreater(4)
```

# License

The code is published under [Apache License Version 2.0](LICENSE)
