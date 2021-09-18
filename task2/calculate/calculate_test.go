package calculate

import (
	"fmt"
	"strings"
	"testing"
)

func TestCalc(t *testing.T) {
	test(t, "1 + 2", 3)
	test(t, "2 - 4", -2)
	test(t, "2 * 4", 8)
	test(t, "2 / 4", 0.5)
	test(t, "1.1 * 2.2", 2.42)
	test(t, "1 - (4 + 3)", -6)
	test(t, "- (4 + 4)", -8)
	test(t, "(1 + 2) * 3", 9)
	test(t, " (1 + 2) - 3 ", 0)
	test(t, "5 - 9 * 8", -67)
	test(t, "(3 + 2.6 / 7 - 0.3 * 234.7) + 23 * 25.1", 510.2614285714285714)
	test(t, "12 - 97 * 98 / 2.3 + 4", -4117.0434782608695652)
	test(t, "(4 / 3 * ( 3.5 + 3 - 5) / 2 + 45 * 8) *2", 722)
}

func test(t *testing.T, expr string, expect float64) {
	r, err := Calculate(expr)
	fmt.Printf(strings.Replace(expr, "%", "%%", -1)+" = %f, expect %f\n", r, expect)
	if err != nil {
		fmt.Println(err)
	}
	if !floatEquals(r, expect) {
		t.Fail()
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
