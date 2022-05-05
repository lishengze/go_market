package extools

import "exterior-interactor/pkg/xmath"

var DecimalStringComparator = func(a, b interface{}) int {
	aAsserted := xmath.MustDecimal(a.(string))
	bAsserted := xmath.MustDecimal(b.(string))
	switch {
	case aAsserted.GreaterThan(bAsserted):
		return 1
	case aAsserted.LessThan(bAsserted):
		return -1
	default:
		return 0
	}
}
