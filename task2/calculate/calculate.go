package calculate

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func Calculate(str string) (float64, error) {
	root, err := parser.ParseExpr(str)

	if err != nil {
		return -1, err
	}
	return eval(root)
}

func eval(str ast.Expr) (float64, error) {
	switch t := str.(type) {
	case *ast.ParenExpr:
		return eval(str.(*ast.ParenExpr).X)
	case *ast.BasicLit:
		return basic(str.(*ast.BasicLit))
	case *ast.BinaryExpr:
		return binary(str.(*ast.BinaryExpr))
	case *ast.UnaryExpr:
		return unary(str.(*ast.UnaryExpr))
	default:
		_ = t
		return -1, errors.New("can't evaluate this expression")
	}
}

func basic(lit *ast.BasicLit) (float64, error) {
	switch lit.Kind {
	case token.FLOAT:
		i, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return -1, err
		} else {
			return i, err
		}
	case token.INT:
		i, err := strconv.ParseInt(lit.Value, 10, 64)
		if err != nil {
			return -1, err
		} else {
			return float64(i), err
		}
	default:
		return -1, errors.New("unknown token")
	}
}

func binary(str *ast.BinaryExpr) (result float64, err error) {
	x, error1 := eval(str.X)
	y, error2 := eval(str.Y)

	result = -1

	if error1 == nil && error2 == nil {
		switch str.Op {
		case token.ADD:
			result = x + y
		case token.SUB:
			result = x - y
		case token.MUL:
			result = x * y
		case token.QUO:
			result = x / y
		default:
			err = errors.New("unknown operator")
		}
	} else {
		if error1 != nil {
			err = error1
		} else {
			err = error2
		}
	}
	return
}

func unary(str *ast.UnaryExpr) (float64, error) {
	x, err := eval(str.X)

	result := -1.0

	if err == nil {
		switch str.Op {
		case token.ADD:
			result = x
		case token.SUB:
			result = -x
		default:
			err = errors.New("unknown operator")
		}
	} else {
		return -1, err
	}
	return float64(result), err
}
