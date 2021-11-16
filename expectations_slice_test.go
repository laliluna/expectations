package expectations_test

import (
	"strings"
	"testing"

	"github.com/laliluna/expectations"
)

type ArrayTestCase struct {
	Fn            func(...interface{}) *expectations.SliceExpectation
	ExpectedValue []interface{}
	Succeeds      bool
}

func TestSliceExpectations(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)
	actualValue := []int{1, 2, 3}
	expect := et.ExpectThatSlice(actualValue)

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

func TestSliceIsEmpty(t *testing.T) {
	et := expectations.NewT(t)
	et.ExpectThatSlice([]int{}).IsEmpty()
}

func TestSliceIsEmptyFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThatSlice([]int{1}).IsEmpty()
	if !tMock.HasBeenCalled {
		t.Error("Slice should not be empty")
	}
}

func TestSliceIsNotEmpty(t *testing.T) {
	et := expectations.NewT(t)
	et.ExpectThatSlice([]int{1}).IsNotEmpty()
}

func TestSliceIsNotEmptyFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThatSlice([]int{}).IsNotEmpty()
	if !tMock.HasBeenCalled {
		t.Error("Slice should not be empty")
	}
}

func TestDeprecatedSliceIsEmpty(t *testing.T) {
	et := expectations.NewT(t)
	et.ExpectThat([]int{}).Slice().IsEmpty()
}

func TestEmptySliceContainsFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThatSlice([]int{}).Contains(1)
	if !tMock.HasBeenCalled {
		t.Error("Slice should not contain")
	}
}
func TestEmptySliceDoesNotContains(t *testing.T) {
	et := expectations.NewT(t)

	et.ExpectThatSlice([]int{}).DoesNotContain(1)
}

func TestSliceHasSizeFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThatSlice([]int{1}).HasSize(2)
	if !tMock.HasBeenCalled {
		t.Error("Slice should have different size")
	}
}

func TestSliceHasSize(t *testing.T) {
	et := expectations.NewT(t)

	et.ExpectThatSlice([]int{1, 2}).HasSize(2)
}

func TestSliceExposeElement(t *testing.T) {
	et := expectations.NewT(t)

	et.ExpectThatSlice([]int{1, 2, 3}).First().Equals(1)
	et.ExpectThatSlice([]int{1, 2, 3}).Second().Equals(2)
	et.ExpectThatSlice([]int{1, 2, 3}).Third().Equals(3)
}

func TestSliceStopsOnFirstFailure(t *testing.T) {
	tMock := &TMock{}
	loggerMock := LoggerMock{}
	et := expectations.NewTWithLogger(tMock, &loggerMock)

	et.ExpectThatSlice([]int{1}).HasSize(2).First().Equals(1)
	if !tMock.HasBeenCalled {
		t.Error("Should fail with has size")

	}

	if !strings.Contains(loggerMock.logs, "Expect len of [1] []int to be 2 and not 1") {
		t.Errorf("Expected '%v' should contain 'Expect len of [1] []int to be 2 and not 1'", loggerMock.logs)
	}
}

func TestSliceSecondShouldNotPanic(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThatSlice([]int{1}).Second()
	if !tMock.HasBeenCalled {
		t.Error("Second should not panic with index ouf of range")
	}
}
