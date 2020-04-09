package expr

type Expr interface {
	Visit(Visitor) interface{}
}

type Visitor interface {
	VisitExpr1(*Expr1) interface{}
	VisitExpr2(*Expr2) interface{}
}

type Expr1 struct {
	x int
	y string
}

func (e *Expr1) Visit(v Visitor) interface{} {
	return v.VisitExpr1(e)
}

type Expr2 struct {
	a bool
	b []byte
}

func (e *Expr2) Visit(v Visitor) interface{} {
	return v.VisitExpr2(e)
}

