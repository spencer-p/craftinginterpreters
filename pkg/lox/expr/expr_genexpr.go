package expr

type Expr interface {
	Visit(Visitor) interface{}
}

type Visitor interface {
	VisitBinary(*Binary) interface{}
	VisitGrouping(*Grouping) interface{}
	VisitLiteral(*Literal) interface{}
	VisitUnary(*Unary) interface{}
}

type Binary struct {
	Left Expr
	Right Expr
	Op tok.Token
}

func (e *Binary) Visit(v Visitor) interface{} {
	return v.VisitBinary(e)
}

type Grouping struct {
	Expression Expr
}

func (e *Grouping) Visit(v Visitor) interface{} {
	return v.VisitGrouping(e)
}

type Literal struct {
	Value interface{}
}

func (e *Literal) Visit(v Visitor) interface{} {
	return v.VisitLiteral(e)
}

type Unary struct {
	Op tok.Token
	Right Expr
}

func (e *Unary) Visit(v Visitor) interface{} {
	return v.VisitUnary(e)
}

