package stmt

// note that Type here is qualified as expr.Type, the expression of a type

//go:generate go run github.com/spencer-p/craftinginterpreters/cmd/genexpr
/// import github.com/spencer-p/craftinginterpreters/pkg/lox/expr
/// import github.com/spencer-p/craftinginterpreters/pkg/lox/tok
/// Expression: Expr expr.Type
/// Print: Expr expr.Type
/// Var: Name tok.Token, Initializer expr.Type
