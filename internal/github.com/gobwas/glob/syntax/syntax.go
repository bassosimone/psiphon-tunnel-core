package syntax

import (
	"github.com/Psiphon-Labs/psiphon-tunnel-core/internal/github.com/gobwas/glob/syntax/ast"
	"github.com/Psiphon-Labs/psiphon-tunnel-core/internal/github.com/gobwas/glob/syntax/lexer"
)

func Parse(s string) (*ast.Node, error) {
	return ast.Parse(lexer.NewLexer(s))
}

func Special(b byte) bool {
	return lexer.Special(b)
}
