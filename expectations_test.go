package expectations_test

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"laliluna.de/expectations"
)

func TestUsageDemo(t *testing.T) {
	eT := expectations.NewT(t)
	eT.Expect(5).Equals(5)
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

func TestSupportsBasicTypes(t *testing.T) {
	eT := expectations.NewT(t)

	var i int = 5
	eT.Expect(i).ToBeGreater(i - 1)
	var i8 int8 = 5
	eT.Expect(i8).ToBeGreater(i8 - 1)
	var i16 int16 = 5
	eT.Expect(i16).ToBeGreater(i16 - 1)
	var i32 int32 = 5
	eT.Expect(i32).ToBeGreater(i32 - 1)
	var i64 int64 = 5
	eT.Expect(i64).ToBeGreater(i64 - 1)
	var u uint = 5
	eT.Expect(u).ToBeGreater(u - 1)
	var up uintptr = 5
	eT.Expect(up).ToBeGreater(up - 1)
	var b byte = 5
	eT.Expect(b).ToBeGreater(b - 1)
	var r rune = 5
	eT.Expect(r).ToBeGreater(r - 1)
	var f32 float32 = 5
	eT.Expect(f32).ToBeGreater(f32 - 1)
	var f64 float64 = 5
	eT.Expect(f64).ToBeGreater(f64 - 1)
	var c64 complex64 = 5
	eT.Expect(c64).Equals(c64)
	var c128 complex128 = 5
	eT.Expect(c128).Equals(c128)
	var sLong string = "foo"
	var sShort string = "fo"
	eT.Expect(sLong).ToBeGreater(sShort)
}

type NumberTestCase struct {
	Fn            func(interface{}) *expectations.Expectation
	ExpectedValue interface{}
	Succeeds      bool
}

func TestIntegerExpectations(t *testing.T) {
	tMock := &TMock{}
	eT := expectations.NewT(tMock)
	actualValue := 2
	expect := eT.Expect(actualValue)

	testCases := []NumberTestCase{
		NumberTestCase{expect.Equals, 2, true},
		NumberTestCase{expect.Equals, 1, false},
		NumberTestCase{expect.Equals, nil, false},
		NumberTestCase{expect.Equals, "foo", false},
		NumberTestCase{expect.DoesNotEqual, 1, true},
		NumberTestCase{expect.DoesNotEqual, 2, false},
		NumberTestCase{expect.DoesNotEqual, nil, true},
		NumberTestCase{expect.DoesNotEqual, "foo", true},
		NumberTestCase{expect.ToBeGreater, 1, true},
		NumberTestCase{expect.ToBeGreater, 2, false},
		NumberTestCase{expect.ToBeGreater, 1.2, false},
		NumberTestCase{expect.ToBeGreaterOrEqual, 1, true},
		NumberTestCase{expect.ToBeGreaterOrEqual, 2, true},
		NumberTestCase{expect.ToBeGreaterOrEqual, 3, false},
		NumberTestCase{expect.ToBeLower, 3, true},
		NumberTestCase{expect.ToBeLower, 2, false},
		NumberTestCase{expect.ToBeLowerOrEqual, 3, true},
		NumberTestCase{expect.ToBeLowerOrEqual, 2, true},
		NumberTestCase{expect.ToBeLowerOrEqual, 1, false},
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

func TestFloatExpectations(t *testing.T) {
	tMock := &TMock{}
	eT := expectations.NewT(tMock)
	actualValue := 2.2
	expect := eT.Expect(actualValue)

	testCases := []NumberTestCase{
		NumberTestCase{expect.Equals, 2.2, true},
		NumberTestCase{expect.Equals, 1.0, false},
		NumberTestCase{expect.Equals, nil, false},
		NumberTestCase{expect.Equals, "foo", false},
		NumberTestCase{expect.DoesNotEqual, 1.1, true},
		NumberTestCase{expect.DoesNotEqual, 2.2, false},
		NumberTestCase{expect.DoesNotEqual, nil, true},
		NumberTestCase{expect.DoesNotEqual, "foo", true},
		NumberTestCase{expect.ToBeGreater, 1.0, true},
		NumberTestCase{expect.ToBeGreater, 2.2, false},
		NumberTestCase{expect.ToBeGreater, 3.2, false},
		NumberTestCase{expect.ToBeGreaterOrEqual, 1.0, true},
		NumberTestCase{expect.ToBeGreaterOrEqual, 2.2, true},
		NumberTestCase{expect.ToBeGreaterOrEqual, 3.0, false},
		NumberTestCase{expect.ToBeLower, 3.0, true},
		NumberTestCase{expect.ToBeLower, 2.0, false},
		NumberTestCase{expect.ToBeLowerOrEqual, 3.0, true},
		NumberTestCase{expect.ToBeLowerOrEqual, 2.2, true},
		NumberTestCase{expect.ToBeLowerOrEqual, 1.0, false},
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

type StringTestCase struct {
	Fn            func(interface{}) *expectations.StringExpectation
	ExpectedValue interface{}
	Succeeds      bool
}

func TestStringExpectations(t *testing.T) {

	tMock := &TMock{}
	et := expectations.NewT(tMock)
	actualValue := "FooBoo"
	expect := et.ExpectString(actualValue)

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
	expect := et.ExpectString(actualValue)

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
	expect := et.ExpectSlice(actualValue)

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

type LoggerMock struct{ logs string }

func (lm *LoggerMock) Log(message string) {
	fmt.Println(message)
	lm.logs += message
}

func TestStopOnFirstFailure(t *testing.T) {
	tMock := &TMock{}

	loggerMock := LoggerMock{}
	et := expectations.NewTWithLogger(tMock, &loggerMock)
	et.Expect(2).Equals(3).ToBeLower(1)
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
