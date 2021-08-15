package expectations_test

import (
	"testing"

	"github.com/laliluna/expectations"
)

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

func TestSliceIsEmpty(t *testing.T) {
	et := expectations.NewT(t)
	et.ExpectThat([]int{}).Slice().IsEmpty()
}

func TestSliceIsEmptyFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThat([]int{1}).Slice().IsEmpty()
	if !tMock.HasBeenCalled {
		t.Error("Slice should not be empty")
	}
}

func TestSliceIsNotEmpty(t *testing.T) {
	et := expectations.NewT(t)
	et.ExpectThat([]int{1}).Slice().IsNotEmpty()
}

func TestSliceIsNotEmptyFails(t *testing.T) {
	tMock := &TMock{}
	et := expectations.NewT(tMock)

	et.ExpectThat([]int{}).Slice().IsNotEmpty()
	if !tMock.HasBeenCalled {
		t.Error("Slice should not be empty")
	}
}
