package vm_test

import (
	"monkey/object"
	"monkey/testutil"
	"monkey/vm"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 - 5", 45},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 * 10)", 100},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	runVmTest(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
		{"!(if (false) {5})", true},
		{"if ((if (false) {10})) {10} else {20}", 20},
	}

	runVmTest(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}

	runVmTest(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", vm.Null},
		{"if (false) { 10 }", vm.Null},
	}

	runVmTest(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[1, 2, 3]", []int{1, 2, 3}},
		{"[1 + 2, 3 * 4, 5 + 6]", []int{3, 12, 11}},
	}

	runVmTest(t, tests)
}

func TestHashLiterals(t *testing.T) {
	tests := []vmTestCase{
		{"{}", map[object.HashKey]int64{}},
		{"{1: 2, 2: 3}", map[object.HashKey]int64{
			(&object.Integer{Value: 1}).HashKey(): 2,
			(&object.Integer{Value: 2}).HashKey(): 3,
		}},
		{"{1 + 1: 2 * 2, 3 + 3: 4 * 4}", map[object.HashKey]int64{
			(&object.Integer{Value: 2}).HashKey(): 4,
			(&object.Integer{Value: 6}).HashKey(): 16,
		}},
	}

	runVmTest(t, tests)
}

func TestIndexExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", vm.Null},
		{"[1, 2, 3][99]", vm.Null},
		{"[1][-1]", vm.Null},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", vm.Null},
		{"{}[0]", vm.Null},
	}

	runVmTest(t, tests)
}

func TestGlobalLetStatements(t *testing.T) {
	tests := []vmTestCase{
		{"let one = 1; one", 1},
		{"let one = 1; let two = 2; one + two", 3},
		{"let one = 1; let two = one + one; one + two", 3},
	}

	runVmTest(t, tests)
}

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{"let fivePlusTen = fn(){5+10;}; fivePlusTen();", 15},
		{"let one = fn() { 1; }; let two = fn() { 2; }; one() + two();", 3},
		{"let a = fn() { 1; }; let b = fn() { a() + 1; }; let c = fn() { b() + 1; }; c();", 3},
	}

	runVmTest(t, tests)
}

func TestFunctionsWithReturnStatement(t *testing.T) {
	tests := []vmTestCase{
		{"let earlyExit = fn() { return 99; 100; }; earlyExit()", 99},
		{"let earlyExit = fn() { return 99; return 100; }; earlyExit()", 99},
	}

	runVmTest(t, tests)
}

func TestFunctionsWithNoReturnValue(t *testing.T) {
	tests := []vmTestCase{
		{"let noReturn = fn() { }; noReturn()", vm.Null},
		{"let noReturn = fn() { }; let noReturnTwo = fn() { noReturn() }; noReturn(); noReturnTwo();", vm.Null},
	}

	runVmTest(t, tests)
}

func TestFirstClassFunctions(t *testing.T) {
	tests := []vmTestCase{
		{"let returnsOne = fn() { 1; }; let returnsOneReturner = fn() { returnsOne; }; returnsOneReturner()();", 1},
	}

	runVmTest(t, tests)
}
func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let one = fn() { let one = 1; one; }; one();`,
			expected: 1,
		},
		{
			input:    `let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }; oneAndTwo();`,
			expected: 3,
		},
		{
			input:    `let oneAndTwo = fn() { let one = 1; let two = 2; one + two; }; let threeAndFour = fn() { let three = 3; let four = 4; three + four; }; oneAndTwo() + threeAndFour();`,
			expected: 10,
		},
		{
			input:    `let firstFoobar = fn() { let foobar = 50; foobar; }; let secondFoobar = fn() { let foobar = 100; foobar; }; firstFoobar() + secondFoobar();`,
			expected: 150,
		},
		{
			input: `let globalSeed = 50;
			let minusOne = fn() { let num = 1; globalSeed - num; };
			let minusTwo = fn() { let num = 2; globalSeed - num; };
			minusOne() + minusTwo();`,
			expected: 97,
		},
	}

	runVmTest(t, tests)
}

func TestCallingFunctionsWithArgumentsAndBindings(t *testing.T) {
	tests := []vmTestCase{
		{
			input:    `let identity = fn(a) { a; }; identity(4);`,
			expected: 4,
		},
		{
			input:    `let sum = fn(a, b) { a + b; }; sum(1, 2);`,
			expected: 3,
		},
		{
			input:    `let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2);`,
			expected: 3,
		},
		{
			input:    `let sum = fn(a, b) { let c = a + b; c; }; sum(1, 2) + sum(3, 4);`,
			expected: 10,
		},
		{
			input:    `let sum = fn(a, b) { let c = a + b; c; }; let outer = fn() { sum(1, 2) + sum(3, 4); }; outer();`,
			expected: 10,
		},
		{
			input: `let globalNum = 10;
			let sum = fn(a, b) {let c = a + b; c + globalNum};
			let outer = fn() { sum(1, 2) + sum(3, 4) + globalNum; };
			outer() + globalNum;`,
			expected: 50,
		},
	}

	runVmTest(t, tests)
}

func TestCallingFunctionsWithWrongArguments(t *testing.T) {
	tests := []struct{ input, expected string }{
		{
			input:    `fn() { 1; }(1);`,
			expected: "wrong number of arguments: want=0, got=1",
		},
		{
			input:    `fn(a) { a; }();`,
			expected: "wrong number of arguments: want=1, got=0",
		},
		{
			input:    `fn(a, b) { a + b; }(1);`,
			expected: "wrong number of arguments: want=2, got=1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			comp := testutil.Compile(t, tt.input)
			vm := vm.New(comp.Bytecode())
			err := vm.Run()
			assert.Error(t, err)
			assert.ErrorContainsf(t, err, tt.expected, "wrong VM error: want=%q, got=%q,", tt.expected, err)
		})
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []vmTestCase{
		{`len("")`, 0},
		{`len("four")`, 4},
		{`len("hello world")`, 11},
		{`len(1)`, &object.Error{Message: "argument to `len` not supported, got INTEGER"}},
		{`len("one", "two")`, &object.Error{Message: "wrong number of arguments. got=2, want=1"}},
		{`len([1, 2, 3])`, 3},
		{`len([])`, 0},
		{`puts("hello", "world!")`, vm.Null},
		{`first([1, 2, 3])`, 1},
		{`first([])`, vm.Null},
		{`first(1)`, &object.Error{Message: "argument to `first` must be an ARRAY, got INTEGER"}},
		{`last([1, 2, 3])`, 3},
		{`last([])`, vm.Null},
		{`last(1)`, &object.Error{Message: "argument to `last` must be an ARRAY, got INTEGER"}},
		{`rest([1, 2, 3])`, []int{2, 3}},
		{`rest([])`, vm.Null},
		{`push([], 1)`, []int{1}},
		{`push(1)`, &object.Error{Message: "wrong number of arguments. got=1, want=2"}},
		{`push(1, 1)`, &object.Error{Message: "argument to `push` must be an ARRAY, got INTEGER"}},
	}

	runVmTest(t, tests)
}

type vmTestCase struct {
	input    string
	expected any
}

func runVmTest(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			comp := testutil.Compile(t, tt.input)
			vm := vm.New(comp.Bytecode())
			err := vm.Run()
			require.NoError(t, err)
			stackElm := vm.LastPoppedStackElem()

			testutil.AssertObject(t, stackElm, tt.expected)
		})
	}
}
