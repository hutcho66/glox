package interpreter

import "github.com/hutcho66/glox/src/pkg/token"

type LoxObject interface {
	get(name *token.Token) (any, error)
}
