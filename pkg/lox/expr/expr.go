package expr

// note that Type here is qualified as expr.Type, the expression of a type

//go:generate go run github.com/spencer-p/craftinginterpreters/cmd/genexpr
/// import github.com/spencer-p/craftinginterpreters/pkg/lox/tok
/// Binary: Left Type, Right Type, Op tok.Token
/// Grouping: Expr Type
/// Literal: Value interface{}
/// Unary: Op tok.Token, Right Type
/// Variable: Name tok.Token
/// Assign: Name tok.Token, Value Type
