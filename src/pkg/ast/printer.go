package ast

import "fmt"

type AstPrinter struct {}

func NewAstPrinter() *AstPrinter {
	return &AstPrinter{};
}

func (p *AstPrinter) Print(e Expression) (any, error) {
	return e.Accept(p);
}

func (AstPrinter) VisitBinaryExpression(b *BinaryExpression) (any, error) {
	return fmt.Sprintf("%s %s %s", b.left, b.operator.GetLexeme(), b.right), nil;
}

func (AstPrinter) VisitUnaryExpression(u *UnaryExpression) (any, error) {
	return fmt.Sprintf("%s%s", u.operator.GetLexeme(), u.expr), nil;
}

func (AstPrinter) VisitGroupedExpression(g *GroupingExpression) (any, error) {
	return fmt.Sprintf("(%s)", g.expr), nil;
}

func (AstPrinter) VisitLiteralExpression(l *LiteralExpression) (any, error) {
	switch v := l.value.(type) {
		case float64: return fmt.Sprintf("%.2f\n", v), nil;
		case bool:    return fmt.Sprintf("%t\n", v), nil;
		case nil:     return "nil", nil;
		case string: 	return v, nil;
		default: 		  return fmt.Sprint(v), nil;
	}
}

