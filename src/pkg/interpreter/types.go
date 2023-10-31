package interpreter

type LoxArray []any

type MapPair struct {
	Key   string
	Value any
}
type LoxMap map[int]MapPair
