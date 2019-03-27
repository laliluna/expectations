package expectations

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// FailFunction expects to be a testing.T in your code but is used to allow testing our own code with a mock
type FailFunction interface {
	Fail()
}

type Logger interface {
	Log(message string)
}

type defaultLogger struct{}

// Log writes a message to stdout
func (defaultLogger) Log(message string) {
	fmt.Println(message)
}

// Expectation holds the actual value and is linked to methods allowing to compare it with the expected value
type Expectation struct {
	T      FailFunction
	Logger Logger
	Value  interface{}
	failed bool
}

type Et struct {
	T      FailFunction
	Logger Logger
}

// NewT creates a struct containing a reference to the testing.T and a default Logger
func NewT(t FailFunction) Et {
	return Et{T: t, Logger: defaultLogger{}}
}

// NewTWithLogger creates a struct containing a reference to the testing.T and custome Logger
func NewTWithLogger(t FailFunction, l Logger) Et {
	return Et{T: t, Logger: l}
}

// Expect builds an Expectation which allows to compare the value to expected values
func (aEt *Et) Expect(value interface{}) *Expectation {
	return &Expectation{aEt.T, aEt.Logger, value, false}
}

// Reset sets the failed flag to false so that further expectations can be executed
func (e *Expectation) Reset() {
	e.failed = false
}

// Equals fails test if expected is not equal to value
func (e *Expectation) Equals(expected interface{}) *Expectation {
	if e.failed {
		return e
	}
	if e.Value != expected {
		e.failed = true
		fail(e.T, e.Logger, fmt.Sprintf("Expect %v to equal %v", e.Value, expected))
	}
	return e
}

// DoesNotEqual fails test if expected is equal to value
func (e *Expectation) DoesNotEqual(expected interface{}) *Expectation {
	if e.failed {
		return e
	}
	if e.Value == expected {
		e.failed = true
		fail(e.T, e.Logger, fmt.Sprintf("Expect %v to not equal %v", e.Value, expected))
	}
	return e
}

// ToBeGreater fails test if expected is not greater than value
func (e *Expectation) ToBeGreater(referencedValue interface{}) *Expectation {
	if e.failed {
		return e
	}
	result := doCompare(referencedValue, e.Value)
	if result != greater {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to be greater than %v", result == notComparable, e.Value, referencedValue))
	}
	return e
}

// ToBeGreaterOrEqual fails test if expected is not greater than or equal to value
func (e *Expectation) ToBeGreaterOrEqual(referencedValue interface{}) *Expectation {
	if e.failed {
		return e
	}
	result := doCompare(referencedValue, e.Value)
	if result != greater && result != equal {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to be greater than or equal to %v", result == notComparable, e.Value, referencedValue))
	}
	return e
}

// ToBeLower fails test if expected is not lower than referencedValue
func (e *Expectation) ToBeLower(referencedValue interface{}) *Expectation {
	if e.failed {
		return e
	}
	result := doCompare(referencedValue, e.Value)
	if result != lower {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to be lower than %v", result == notComparable, e.Value, referencedValue))
	}
	return e
}

// ToBeLowerOrEqual fails test if expected is not lower than or equal to referencedValue
func (e *Expectation) ToBeLowerOrEqual(referencedValue interface{}) *Expectation {
	if e.failed {
		return e
	}
	result := doCompare(referencedValue, e.Value)
	if result != lower && result != equal {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to be lower than or equal to %v", result == notComparable, e.Value, referencedValue))
	}
	return e
}

// ===================== Strings ==============================

// StringExpectation allows to express expectations on strings
type StringExpectation struct {
	T      FailFunction
	Logger Logger
	Value  interface{}
	failed bool
}

// ExpectString builds an Expectation for strings which allows to compare the value to expected values
func (aEt *Et) ExpectString(value interface{}) *StringExpectation {
	return &StringExpectation{aEt.T, aEt.Logger, value, false}
}

// Reset sets the failed flag to false so that further expectations can be executed
func (e *StringExpectation) Reset() {
	e.failed = false
}

// Equals fails test if expected is not equal to value
func (e *StringExpectation) Equals(expected interface{}) *StringExpectation {
	if e.failed {
		return e
	}
	result := compareEquality(expected, e.Value)
	if result != equal {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to equal %v", result == notComparable, e.Value, expected))
	}
	return e
}

// EqualsIgnoringCase fails test if expected is not equal to value
func (e *StringExpectation) EqualsIgnoringCase(expected interface{}) *StringExpectation {
	if e.failed {
		return e
	}
	valueString, valueOk := e.Value.(string)
	expectedString, expectedOk := expected.(string)
	if !(valueOk && expectedOk) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to equal ignoring case %v", showTypeInfos, e.Value, expected))
	} else if strings.ToLower(valueString) != strings.ToLower(expectedString) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to equal ignoring case %v", hideTypeInfos, e.Value, expected))
	}
	return e
}

// DoesNotEqual fails test if expected is equal to value
func (e *StringExpectation) DoesNotEqual(expected interface{}) *StringExpectation {
	if e.failed {
		return e
	}
	if expected == e.Value {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to not equal %v", hideTypeInfos, e.Value, expected))
	}
	return e
}

// StartsWith checks if expected starts with value
func (e *StringExpectation) StartsWith(prefix interface{}) *StringExpectation {
	if e.failed {
		return e
	}
	valueString, valueOk := e.Value.(string)
	prefixString, prefixOk := prefix.(string)
	if !(valueOk && prefixOk) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to start with %v", showTypeInfos, e.Value, prefix))
	} else if !strings.HasPrefix(valueString, prefixString) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to start with %v", hideTypeInfos, e.Value, prefix))
	}
	return e
}

// EndsWith checks if expected starts with value
func (e *StringExpectation) EndsWith(suffix interface{}) *StringExpectation {
	if e.failed {
		return e
	}
	valueString, valueOk := e.Value.(string)
	suffixString, suffixOk := suffix.(string)
	if !(valueOk && suffixOk) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to end with %v", showTypeInfos, e.Value, suffix))
	} else if !strings.HasSuffix(valueString, suffixString) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to end with %v", hideTypeInfos, e.Value, suffix))
	}
	return e
}

// Contains checks if expected contains all expected values
func (e *StringExpectation) Contains(expectedValues ...string) *StringExpectation {
	if e.failed {
		return e
	}
	valueString, valueOk := e.Value.(string)
	if !(valueOk) {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to contain %v", showTypeInfos, e.Value, expectedValues))
		return e
	}

	var lackingValues []string
	var failed bool
	for _, expectedValue := range expectedValues {
		if !strings.Contains(valueString, expectedValue) {
			failed = true
			lackingValues = append(lackingValues, expectedValue)
		}
	}

	if failed {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to contain %v but was missing %v", hideTypeInfos, e.Value, expectedValues, lackingValues))
	}

	return e
}

// DoesNotContain checks if expected does not contain any of the expected values
func (e *StringExpectation) DoesNotContain(expectedValues ...string) *StringExpectation {
	if e.failed {
		return e
	}
	valueString, valueOk := e.Value.(string)
	if !(valueOk) {
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to not contain %v", showTypeInfos, e.Value, expectedValues))
		return e
	}

	var foundValues []string
	var failed bool
	for _, expectedValue := range expectedValues {
		if strings.Contains(valueString, expectedValue) {
			failed = true
			foundValues = append(foundValues, expectedValue)
		}
	}

	if failed {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to not contain %v but it includes %v", hideTypeInfos, e.Value, expectedValues, foundValues))
	}

	return e
}

// ===================== Slices ==============================

// SliceExpectation allows to express expectations on strings
type SliceExpectation struct {
	T      FailFunction
	Logger Logger
	Value  interface{}
	failed bool
}

// ExpectSlice builds an Expectation for slices which allows to compare the value to expected values
func (aEt *Et) ExpectSlice(value interface{}) *SliceExpectation {
	return &SliceExpectation{aEt.T, aEt.Logger, value, false}
}

// Reset sets the failed flag to false, so that further checks can be executed
func (e *SliceExpectation) Reset() {
	e.failed = false
}

// Contains checks if expected contains all expected values
func (e *SliceExpectation) Contains(expectedValues ...interface{}) *SliceExpectation {
	if e.failed {
		return e
	}
	kind := reflect.TypeOf(e.Value).Kind()
	if !(kind == reflect.Slice || kind == reflect.Array) {
		e.failed = true
		fail(e.T, e.Logger, fmt.Sprintf("Expect %v %T to be a slice", e.Value, e.Value))
		return e
	}

	typesMatch := checkTypesMatch(toSlice(e.Value), expectedValues)

	var lackingValues []interface{}
	for _, expectedValue := range expectedValues {
		if !doContain(e.Value, expectedValue) {
			lackingValues = append(lackingValues, expectedValue)
		}
	}

	if len(lackingValues) > 0 {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to contain %v but was missing %v", typesMatch, e.Value, expectedValues, lackingValues))
	}
	return e
}

// DoesNotContain checks if expected does not contain any of the expected values
func (e *SliceExpectation) DoesNotContain(expectedValues ...interface{}) *SliceExpectation {
	if e.failed {
		return e
	}

	if reflect.TypeOf(e.Value).Kind() != reflect.Slice {
		e.failed = true
		fail(e.T, e.Logger, fmt.Sprintf("Expect %v %T to be a slice", e.Value, e.Value))
		return e
	}

	typesMatch := checkTypesMatch(toSlice(e.Value), expectedValues)

	var additionalValues []interface{}
	for _, expectedValue := range expectedValues {
		if doContain(e.Value, expectedValue) {
			additionalValues = append(additionalValues, expectedValue)
		}
	}

	if len(additionalValues) > 0 {
		e.failed = true
		fail(e.T, e.Logger, buildFailMessage2("Expect %v to not contain %v but it includes %v", typesMatch, e.Value, expectedValues, additionalValues))
	}
	return e
}

func toSlice(value interface{}) []interface{} {
	sourceSlice := reflect.ValueOf(value)
	result := make([]interface{}, sourceSlice.Len(), sourceSlice.Cap())
	for i := 0; i < sourceSlice.Len(); i++ {
		value := sourceSlice.Index(i).Interface()
		result[i] = value
	}
	return result
}

func doContain(sliceValue, expectedValue interface{}) bool {
	for _, value := range toSlice(sliceValue) {
		if value == expectedValue {
			return true
		}
	}
	return false
}

func checkTypesMatch(values, expectedValues []interface{}) bool {
	valueType := reflect.TypeOf(values[0])
	for _, expectedValue := range expectedValues {
		if valueType != reflect.TypeOf(expectedValue) {
			return false
		}
	}
	return true
}

const (
	greater       = iota
	lower         = iota
	equal         = iota
	notComparable = iota
	comparable    = iota
	notEqual      = iota
)

func compareEquality(expected interface{}, actual interface{}) uint {
	if expected == actual {
		return equal
	}

	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return notComparable
	}

	return notEqual
}

func compareInt(expected int64, actual int64) uint {
	switch {
	case actual > expected:
		return greater
	case actual < expected:
		return lower
	default:
		return equal
	}
}

func compareUint(expected uint64, actual uint64) uint {
	switch {
	case actual > expected:
		return greater
	case actual < expected:
		return lower
	default:
		return equal
	}
}

func compareFloat(expected float64, actual float64) uint {
	switch {
	case actual > expected:
		return greater
	case actual < expected:
		return lower
	default:
		return equal
	}
}

func compareString(expected string, actual string) uint {
	switch {
	case actual > expected:
		return greater
	case actual < expected:
		return lower
	default:
		return equal
	}
}

func doCompare(expected interface{}, actual interface{}) uint {
	if expected == actual {
		return equal
	}
	if expected == nil || actual == nil {
		return notComparable
	}
	if reflect.TypeOf(expected) != reflect.TypeOf(actual) {
		return notComparable
	}

	switch expected.(type) {
	case int:
		return compareInt(int64(expected.(int)), int64(actual.(int)))
	case int8:
		return compareInt(int64(expected.(int8)), int64(actual.(int8)))
	case int16:
		return compareInt(int64(expected.(int16)), int64(actual.(int16)))
	case int32:
		return compareInt(int64(expected.(int32)), int64(actual.(int32)))
	case int64:
		return compareInt(expected.(int64), actual.(int64))
	case uint:
		return compareUint(uint64(expected.(uint)), uint64(actual.(uint)))
	case uint8:
		return compareUint(uint64(expected.(uint8)), uint64(actual.(uint8)))
	case uint16:
		return compareUint(uint64(expected.(uint16)), uint64(actual.(uint16)))
	case uint32:
		return compareUint(uint64(expected.(uint32)), uint64(actual.(uint32)))
	case uint64:
		return compareUint(expected.(uint64), actual.(uint64))
	case uintptr:
		return compareUint(uint64(expected.(uintptr)), uint64(actual.(uintptr)))
	case float32:
		return compareFloat(float64(expected.(float32)), float64(actual.(float32)))
	case float64:
		return compareFloat(expected.(float64), actual.(float64))
	case string:
		return compareString(expected.(string), actual.(string))
	}
	return notComparable
}

func buildFailMessage(message string, actual, expected interface{}, typeOfFailure uint) string {
	result := fmt.Sprintf(message, addTypeIfNeeded(actual, typeOfFailure), addTypeIfNeeded(expected, typeOfFailure))
	if typeOfFailure == notComparable {
		return "Types do not match. " + result
	}

	return result
}

func doMap(source []interface{}, fn func(interface{}) interface{}) []interface{} {
	result := make([]interface{}, len(source))
	for i := 0; i < len(source); i++ {
		result[i] = fn(source[i])
	}
	return result
}

const showTypeInfos = true
const hideTypeInfos = true

func buildFailMessage2(message string, showTypeInfos bool, args ...interface{}) string {
	formattedArgs := doMap(args, func(value interface{}) interface{} {
		if showTypeInfos {
			return fmt.Sprintf("%v", value)
		}
		return addTypes(value)
	})

	return fmt.Sprintf(message, formattedArgs...)
}

func addTypeIfNeeded(value interface{}, typeOfFailure uint) string {
	if typeOfFailure == notComparable {
		return fmt.Sprintf("%v (%T)", value, value)
	}
	return fmt.Sprintf("%v", value)
}

func addTypes(value interface{}) string {
	if reflect.TypeOf(value).Kind() == reflect.Slice {
		result := ""
		for _, item := range toSlice(value) {
			if result != "" {
				result += ", "
			}
			result += fmt.Sprintf("%v (%T)", item, item)
		}
		return fmt.Sprintf("[%v]", result)
	}

	return fmt.Sprintf("%v (%T)", value, value)
}

var lastFileName string

func fail(f FailFunction, l Logger, message string) {

	fileName, methodName, line := determineCodeLocation()
	if lastFileName != fileName {
		l.Log(fileName)
		l.Log(strings.Repeat("-", len(fileName)))
		lastFileName = fileName
	}
	l.Log(fmt.Sprintf("--- %v in line %v: %v\n", methodName, line, message))
	f.Fail()
}

func determineCodeLocation() (string, string, int) {
	fileName := getFrame(3).File[strings.LastIndex(getFrame(3).File, "/")+1:]
	methodName := getFrame(3).Function[strings.LastIndex(getFrame(3).Function, ".")+1:]
	line := getFrame(3).Line
	if strings.Contains(methodName, "-fm") && methodName[len(methodName)-3:] == "-fm" {
		methodName = getFrame(4).Function[strings.LastIndex(getFrame(4).Function, ".")+1:]
	}
	return fileName, methodName, line
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}