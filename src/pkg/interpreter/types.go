package interpreter

type LoxArray []any

type MapPair struct {
	key   string
	value any
}
type LoxMap map[int]MapPair
