package logic

import "github.com/pavlo67/workshop/common/selectors"

func AND(value0, value1 interface{}) *selectors.Term {
	return &selectors.Term{selectors.Operand(value0), selectors.Operand(value1), selectors.And}
}

//func OR(value0, value1 interface{}) *selectors.Term {
//	return &selectors.Term{selectors.TermUnary{ValueUnary: value0}, []selectors.TermRight{{selectors.TermUnary{ValueUnary: value1}, selectors.Or}}}
//}
//
//func NOT(value interface{}) *selectors.Term {
//	return &selectors.Term{selectors.TermUnary{ValueUnary: value, OperationUnary: selectors.Not}, nil}
//}