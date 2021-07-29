package expectations_test

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
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

func TestSupportsBasicTypes(t *testing.T) {
	eT := expectations.NewTWithLogger(t, &SilenceLoggerMock{})

	var i int = 5
	eT.ExpectThat(i).IsGreater(i - 1)
	var i8 int8 = 5
	eT.ExpectThat(i8).IsGreater(i8 - 1)
	var i16 int16 = 5
	eT.ExpectThat(i16).IsGreater(i16 - 1)
	var i32 int32 = 5
	eT.ExpectThat(i32).IsGreater(i32 - 1)
	var i64 int64 = 5
	eT.ExpectThat(i64).IsGreater(i64 - 1)
	var u uint = 5
	eT.ExpectThat(u).IsGreater(u - 1)
	var up uintptr = 5
	eT.ExpectThat(up).IsGreater(up - 1)
	var b byte = 5
	eT.ExpectThat(b).IsGreater(b - 1)
	var r rune = 5
	eT.ExpectThat(r).IsGreater(r - 1)
	var f32 float32 = 5
	eT.ExpectThat(f32).IsGreater(f32 - 1)
	var f64 float64 = 5
	eT.ExpectThat(f64).IsGreater(f64 - 1)
	var c64 complex64 = 5
	eT.ExpectThat(c64).Equals(c64)
	var c128 complex128 = 5
	eT.ExpectThat(c128).Equals(c128)
	var sLong string = "foo"
	var sShort string = "fo"
	eT.ExpectThat(sLong).IsGreater(sShort)
}

type NumberTestCase struct {
	Fn            func(interface{}) *expectations.Expectation
	ExpectedValue interface{}
	Succeeds      bool
}

func TestIntegerExpectations(t *testing.T) {
	tMock := &TMock{}
	eT := expectations.NewTWithLogger(tMock, &SilenceLoggerMock{})
	actualValue := 2
	expect := eT.ExpectThat(actualValue)

	testCases := []NumberTestCase{
		NumberTestCase{expect.Equals, 2, true},
		NumberTestCase{expect.Equals, 1, false},
		NumberTestCase{expect.Equals, nil, false},
		NumberTestCase{expect.Equals, "foo", false},
		NumberTestCase{expect.DoesNotEqual, 1, true},
		NumberTestCase{expect.DoesNotEqual, 2, false},
		NumberTestCase{expect.DoesNotEqual, nil, true},
		NumberTestCase{expect.DoesNotEqual, "foo", false}, // reject to compare different types
		NumberTestCase{expect.IsGreater, 1, true},
		NumberTestCase{expect.IsGreater, 2, false},
		NumberTestCase{expect.IsGreater, 1.2, false},
		NumberTestCase{expect.IsGreaterOrEqual, 1, true},
		NumberTestCase{expect.IsGreaterOrEqual, 2, true},
		NumberTestCase{expect.IsGreaterOrEqual, 3, false},
		NumberTestCase{expect.IsLower, 3, true},
		NumberTestCase{expect.IsLower, 2, false},
		NumberTestCase{expect.IsLowerOrEqual, 3, true},
		NumberTestCase{expect.IsLowerOrEqual, 2, true},
		NumberTestCase{expect.IsLowerOrEqual, 1, false},
	}

	for _, testCase := range testCases {
		tMock.reset()
		testCase.Fn(testCase.ExpectedValue)
		if testCase.Succeeds == tMock.HasBeenCalled {
			t.Errorf("Test failed: %v %v %v should be %v", actualValue, functionName(testCase.Fn), testCase.ExpectedValue, testCase.Succeeds)
		}
		expect.Reset()
	}
}

func TestShowNumbersHaveNotTheSameType(t *testing.T) {
	tMock := &TMock{}

	loggerMock := LoggerMock{}
	et := expectations.NewTWithLogger(tMock, &loggerMock)

	var actualValue uint16 = 1
	et.ExpectThat(actualValue).Equals(1)
	if !strings.Contains(loggerMock.logs, "You try to compare different types") {
		t.Errorf("Expected message to indicate different types")
	}

	loggerMock.Reset()

	et.ExpectThat(actualValue).DoesNotEqual(2)
	if !strings.Contains(loggerMock.logs, "You try to compare different types") {
		t.Errorf("Expected message to indicate different types")
	}
}

func TestFloatExpectations(t *testing.T) {
	tMock := &TMock{}
	eT := expectations.NewTWithLogger(tMock, &SilenceLoggerMock{})
	actualValue := 2.2
	expect := eT.ExpectThat(actualValue)

	testCases := []NumberTestCase{
		NumberTestCase{expect.Equals, 2.2, true},
		NumberTestCase{expect.Equals, 1.0, false},
		NumberTestCase{expect.Equals, nil, false},
		NumberTestCase{expect.Equals, "foo", false},
		NumberTestCase{expect.DoesNotEqual, 1.1, true},
		NumberTestCase{expect.DoesNotEqual, 2.2, false},
		NumberTestCase{expect.DoesNotEqual, nil, true},
		NumberTestCase{expect.DoesNotEqual, "foo", false}, // reject to compare type
		NumberTestCase{expect.IsGreater, 1.0, true},
		NumberTestCase{expect.IsGreater, 2.2, false},
		NumberTestCase{expect.IsGreater, 3.2, false},
		NumberTestCase{expect.IsGreaterOrEqual, 1.0, true},
		NumberTestCase{expect.IsGreaterOrEqual, 2.2, true},
		NumberTestCase{expect.IsGreaterOrEqual, 3.0, false},
		NumberTestCase{expect.IsLower, 3.0, true},
		NumberTestCase{expect.IsLower, 2.0, false},
		NumberTestCase{expect.IsLowerOrEqual, 3.0, true},
		NumberTestCase{expect.IsLowerOrEqual, 2.2, true},
		NumberTestCase{expect.IsLowerOrEqual, 1.0, false},
	}

	for _, testCase := range testCases {
		tMock.reset()
		testCase.Fn(testCase.ExpectedValue)
		if testCase.Succeeds == tMock.HasBeenCalled {
			t.Errorf("Test failed: %v %v %v should be %v", actualValue, functionName(testCase.Fn), testCase.ExpectedValue, testCase.Succeeds)
		}
		expect.Reset()
	}
}

func TestNil(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThat("foo").IsNil()
	if !tMock.HasBeenCalled {
		t.Errorf("Expect 5 to be not nil")
	}
	tMock.reset()

	et.ExpectThat(5).IsNil()
	if !tMock.HasBeenCalled {
		t.Errorf("Expect 5 to be not nil")
	}
	tMock.reset()

	et.ExpectThat(nil).IsNil()
	if tMock.HasBeenCalled {
		t.Errorf("Expect nil to be nil")
	}
	tMock.reset()

	et.ExpectThat("foo").String().IsNil()
	if !tMock.HasBeenCalled {
		t.Errorf("Expect foo to be not nil")
	}
	tMock.reset()

}

func TestNilPointer(t *testing.T) {
	et := expectations.NewT(t)
	type data struct{}
	var x *data
	et.ExpectThat(x).IsNil()
}

func TestNilChannel(t *testing.T) {
	et := expectations.NewT(t)
	var x chan int
	et.ExpectThat(x).IsNil()
}

func TestNilMap(t *testing.T) {
	et := expectations.NewT(t)
	var x map[string]int
	et.ExpectThat(x).IsNil()
}
func TestNilSlice(t *testing.T) {
	et := expectations.NewT(t)
	var x []string
	et.ExpectThat(x).IsNil()
}
func TestNotNil(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThat("foo").IsNotNil()
	if tMock.HasBeenCalled {
		t.Errorf("Expect foo to be not nil")
	}
	tMock.reset()

	et.ExpectThat(5).IsNotNil()
	if tMock.HasBeenCalled {
		t.Errorf("Expect 5 to be not nil")
	}
	tMock.reset()

	et.ExpectThat(nil).IsNotNil()
	if !tMock.HasBeenCalled {
		t.Errorf("Expect nil to be nil")
	}
	tMock.reset()

	et.ExpectThat("foo").String().IsNotNil()
	if tMock.HasBeenCalled {
		t.Errorf("Expect 'foo' to be not nil")
	}
	tMock.reset()
}

type StringTestCase struct {
	Fn            func(interface{}) *expectations.StringExpectation
	ExpectedValue interface{}
	Succeeds      bool
}

func TestStringExpectations(t *testing.T) {

	tMock := &TMock{}
	et := expectations.NewT(tMock)
	actualValue := "FooBoo"
	expect := et.ExpectThat(actualValue).String()

	et.ExpectThat(5).String()
	if !tMock.HasBeenCalled {
		t.Errorf("Expect String to fail if value is not a string")
	}

	testCases := []StringTestCase{
		StringTestCase{expect.Equals, actualValue, true},
		StringTestCase{expect.Equals, "Something else", false},
		StringTestCase{expect.Equals, nil, false},
		StringTestCase{expect.EqualsIgnoringCase, "fooboo", true},
		StringTestCase{expect.EqualsIgnoringCase, "Something else", false},
		StringTestCase{expect.DoesNotEqual, "SomethingElse", true},
		StringTestCase{expect.DoesNotEqual, "FooBoo", false},
		StringTestCase{expect.DoesNotEqual, nil, true},
		StringTestCase{expect.StartsWith, "Foo", true},
		StringTestCase{expect.StartsWith, "Boo", false},
		StringTestCase{expect.StartsWith, "foo", false},
		StringTestCase{expect.EndsWith, "Boo", true},
		StringTestCase{expect.EndsWith, "Foo", false},
		StringTestCase{expect.EndsWith, "boo", false},
		StringTestCase{expect.StartsWith, "Boo", false},
		StringTestCase{expect.StartsWith, "foo", false},
	}

	for _, testCase := range testCases {
		tMock.reset()
		testCase.Fn(testCase.ExpectedValue)
		if testCase.Succeeds == tMock.HasBeenCalled {
			t.Errorf("Test failed: %v %v %v should be %v", actualValue, functionName(testCase.Fn), testCase.ExpectedValue, testCase.Succeeds)
		}
		expect.Reset()
	}
}

type StringArrayTestCase struct {
	Fn            func(...string) *expectations.StringExpectation
	ExpectedValue []string
	Succeeds      bool
}

func TestStringContainsExpectations(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)
	actualValue := "FooBoo"
	expect := et.ExpectThat(actualValue).String()

	testCases := []StringArrayTestCase{
		StringArrayTestCase{expect.Contains, []string{"oo", "ooBo"}, true},
		StringArrayTestCase{expect.Contains, []string{"oo", "x"}, false},
		StringArrayTestCase{expect.Contains, []string{"x"}, false},
		StringArrayTestCase{expect.DoesNotContain, []string{"ox", "obo"}, true},
		StringArrayTestCase{expect.DoesNotContain, []string{"x", "oo"}, false},
		StringArrayTestCase{expect.DoesNotContain, []string{"oo"}, false},
	}

	for _, testCase := range testCases {
		tMock.reset()
		testCase.Fn(testCase.ExpectedValue...)
		if testCase.Succeeds == tMock.HasBeenCalled {
			t.Errorf("Test failed: %v %v %v should be %v", actualValue, functionName(testCase.Fn), testCase.ExpectedValue, testCase.Succeeds)
		}
		expect.Reset()
	}
}

type ArrayTestCase struct {
	Fn            func(...interface{}) *expectations.SliceExpectation
	ExpectedValue []interface{}
	Succeeds      bool
}

func TestSliceExpectations(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)
	actualValue := []int{1, 2, 3}
	expect := et.ExpectThat(actualValue).Slice()

	testCases := []ArrayTestCase{
		ArrayTestCase{expect.Contains, []interface{}{1, 3}, true},
		ArrayTestCase{expect.Contains, []interface{}{1, 3.3}, false},
		ArrayTestCase{expect.Contains, []interface{}{1, 5}, false},
		ArrayTestCase{expect.Contains, []interface{}{5}, false},
		ArrayTestCase{expect.DoesNotContain, []interface{}{7, 11}, true},
		ArrayTestCase{expect.DoesNotContain, []interface{}{7, 2}, false},
		ArrayTestCase{expect.DoesNotContain, []interface{}{2}, false},
	}

	for _, testCase := range testCases {
		tMock.reset()
		testCase.Fn(testCase.ExpectedValue...)
		if testCase.Succeeds == tMock.HasBeenCalled {
			t.Errorf("Test failed: %v %v %v should be %v", actualValue, functionName(testCase.Fn), testCase.ExpectedValue, testCase.Succeeds)
		}
		expect.Reset()
	}
}

type SilenceLoggerMock struct{}

func (lm *SilenceLoggerMock) Log(message string) {

}

type LoggerMock struct{ logs string }

func (lm *LoggerMock) Log(message string) {
	fmt.Println(message)
	lm.logs += message
}

func (lm *LoggerMock) Reset() {
	lm.logs = ""
}

func TestStopOnFirstFailure(t *testing.T) {
	tMock := &TMock{}

	loggerMock := LoggerMock{}
	et := expectations.NewTWithLogger(tMock, &loggerMock)
	et.ExpectThat(2).Equals(3).IsLower(1)
	if !strings.Contains(loggerMock.logs, "to equal") {
		t.Errorf("Expected '%v' should contain 'to equal'", loggerMock.logs)
	}
	if strings.Contains(loggerMock.logs, "to be lower than") {
		t.Errorf("Expected '%v' should not contain 'to be lower than'", loggerMock.logs)
	}
}

type TMock struct {
	HasBeenCalled bool
}

func (t *TMock) Fail() {
	t.HasBeenCalled = true
}

func (t *TMock) reset() {
	t.HasBeenCalled = false
}

func functionName(fn interface{}) string {
	fnNameDetails := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return fnNameDetails[strings.LastIndex(fnNameDetails, ".")+1 : len(fnNameDetails)-3]
}
