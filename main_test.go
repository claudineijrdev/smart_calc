package main

import "testing"

func TestCalc(t *testing.T) {
	testCases := []struct {
		inputs []string
		expect int
		err    bool
	}{
		{[]string{"16", "+", "9"}, 25, false},
		// {[]string{"1", "2"}, 3, false},
		// {[]string{"1", "+", "3"}, 4, false},
		// {[]string{"1", "-", "3"}, -2, false},
		// {[]string{"1", "+", "3", "-", "2"}, 2, false},
		// {[]string{"8"}, 8, false},
		// {[]string{"-2", "+", "4", "-", "5", "+", "6"}, 3, false},
		// {[]string{"9", "+++", "10", "--", "8"}, 27, false},
		// {[]string{"3", "---", "5"}, -2, false},
		// {[]string{"14", " ", " ", "", "-", "12"}, 2, false},
	}

	for _, tc := range testCases {
		result, err := calc(tc.inputs)
		if tc.err {
			if err == nil {
				t.Errorf("Expected error for inputs: %v", tc.inputs)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for inputs: %v", tc.inputs)
			}
			if result != tc.expect {
				t.Errorf("Expected %d, give %d for inputs: %v", tc.expect, result, tc.inputs)
			}
		}
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		inputs         string
		expect         int
		shouldContinue bool
		err            bool
		errMsg         string
	}{
		{"/exit", 0, false, false, ""},
		{"/help", 0, true, true, "The program calculates expressions"},
		{"/command", 0, true, true, "Unknown command"},
		{"16 + 9", 25, true, false, ""},
		{"16+9", 25, true, false, ""},
		{"16 + 9 - 10", 15, true, false, ""},
		{"15 +- 9", 6, true, false, ""},
		{"15 -+ 9", 6, true, false, ""},
		{"15 -- 9", 24, true, false, ""},
		{"15 --- 9", 6, true, false, ""},
		{"3 * 4", 12, true, false, ""},
		{"3 * 4 + 2 * 3", 18, true, false, ""},
		{"3 + 4 * 2 ", 11, true, false, ""},
		{"9 / 3", 3, true, false, ""},
		{"9 / 3 + 2 * 3", 9, true, false, ""},
		{"2 ^ 3", 8, true, false, ""},
		{"2^3", 8, true, false, ""},
		{"2 ^ 3 + 3 * 2", 14, true, false, ""},
		{"2 * (3 + 4)", 14, true, false, ""},

		{"8 * 3 + 12 * (4 - 2)", 48, true, false, ""},
		{"2 - 2 + 3", 3, true, false, ""},
		{"4 * (2 + 3", 0, true, true, "Invalid expression"},
		{"-10", -10, true, false, ""},
		{"a=4", 0, true, false, ""},
		{"b=5", 0, true, false, ""},
		{"c=6", 0, true, false, ""},
		{"a*2+b*3+c*(2+3)", 53, true, false, ""},
		{"1 +++ 2 * 3 -- 4", 11, true, false, ""},
		{"3 *** 5", 0, true, true, "Invalid expression"},
		{"4+3)", 0, true, true, "Invalid expression"},
		{"BIG = 9000", 0, true, false, ""},
		{"BIG", 9000, true, false, ""},
		{"91 / 13", 7, true, false, ""},
	}
	variables := make(map[string]int)
	for _, tc := range testCases {
		result, shouldContinue, _, err := run(tc.inputs, &variables)
		if tc.err {
			if err == nil {
				t.Errorf("Expected error for inputs: %v", tc.inputs)
			}
			if err.Error() != tc.errMsg {
				t.Errorf("Expected error message %s, give %s for inputs: %v", tc.errMsg, err.Error(), tc.inputs)
			}
		} else {
			if err != nil {
				t.Errorf("Unexpected error for inputs: %v", tc.inputs)
			}
			if result != tc.expect {
				t.Errorf("[Result] Expected %d, give %d for inputs: %v", tc.expect, result, tc.inputs)
			}
		}
		if shouldContinue != tc.shouldContinue {
			t.Errorf("[ShouldContinue]Expected %t, give %t for inputs: %v", tc.shouldContinue, shouldContinue, tc.inputs)
		}
	}
}

// func TestNormalizeOperations(t *testing.T) {
// 	testCases := []struct {
// 		input  string
// 		expect string
// 	}{
// 		{"+", "+"},
// 		{"-", "-"},
// 		{"++", "+"},
// 		{"+-", "-"},
// 		{"-+", "-"},
// 		{"--", "+"},
// 		{"---", "-"},
// 		{"+--", "+"},
// 		{"-+-", "+"},;
// 		{"-++", "-"},
// 	}

// 	for _, tc := range testCases {
// 		result, error := normalizeOperations(tc.input)
// 		if result != tc.expect {
// 			t.Errorf("Expected %s for input: %s", tc.expect, tc.input)
// 		}
// 	}
// }
