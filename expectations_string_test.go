package expectations_test

import (
	"testing"

	"github.com/laliluna/expectations"
)

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
