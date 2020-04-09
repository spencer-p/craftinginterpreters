package expr

//go:generate go run github.com/spencer-p/craftinginterpreters/cmd/genexpr
/// import github.com/spencer-p/craftinginterpreters/pkg/lox/tok
/// Binary: Left Expr, Right Expr, Op tok.Token
/// Grouping: Expression Expr
/// Literal: Value interface{}
/// Unary: Op tok.Token, Right Expr
