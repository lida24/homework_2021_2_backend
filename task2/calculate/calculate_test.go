package calculate

import (
	"fmt"
	"testing"
)

var tests = []struct {
	expr   string
	expect float64
}{
	{"1 + 2", 3},
	{"2 - 4", -2},
	{"2 * 4", 8},
	{"2 / 4", 0.5},
	{"1.1 * 2.2", 2.42},
	{"1 - (4 + 3)", -6},
	{"- (4 + 4)", -8},
	{"(1 + 2) * 3", 9},
	{"(1 + 2) - 3", 0},
	{"5 - 9 * 8", -67},
	{"(3 + 2.6 / 7 - 0.3 * 234.7) + 23 * 25.1", 510.2614285714285714},
	{"12 - 97 * 98 / 2.3 + 4", -4117.0434782608695652},
	{"(4 / 3 * ( 3.5 + 3 - 5) / 2 + 45 * 8) *2", 722},
}

func TestFlagParser(t *testing.T) {
	for _, tt := range tests {
		res, err := Calculate(tt.expr)
		if err != nil {
			fmt.Println(err)
		}
		if !floatEquals(res, tt.expect) {
			t.Errorf("Find %s = %f, expected %f", tt.expr, res, tt.expect)
		}
	}
}

var EPSILON float64 = 0.00000001

func floatEquals(a, b float64) bool {
	if (a-b) < EPSILON && (b-a) < EPSILON {
		return true
	} else {
		return false
	}
}
