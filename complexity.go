package gocognito

import (
	"go/ast"
	"go/token"
)

// Complexity calculates the cognitive complexity of a function.
func Complexity(fn *ast.FuncDecl) int {
	v := complexityVisitor{
		name: fn.Name,
	}

	ast.Walk(&v, fn)
	return v.complexity
}

type complexityVisitor struct {
	complexity int
	nesting    int
	name       *ast.Ident
	// to recognize "else if" and just "if"
	justAfterElse bool
	// to count sequences of binary logic operators exclusively
	visitedBinaryExpr map[*ast.BinaryExpr]bool
}

func (v *complexityVisitor) incNesting() {
	v.nesting++
}

func (v *complexityVisitor) decNesting() {
	v.nesting--
}

func (v *complexityVisitor) incComplexity() {
	v.complexity++
}

func (v *complexityVisitor) nestIncComplexity() {
	v.complexity += (v.nesting + 1)
}

func (v *complexityVisitor) markVisited(e *ast.BinaryExpr) {
	if v.visitedBinaryExpr == nil {
		v.visitedBinaryExpr = make(map[*ast.BinaryExpr]bool)
	}

	v.visitedBinaryExpr[e] = true
}

func (v *complexityVisitor) isVisited(e *ast.BinaryExpr) bool {
	if v.visitedBinaryExpr == nil {
		return false
	}

	return v.visitedBinaryExpr[e]
}

func (v *complexityVisitor) extractOpSequence(be *ast.BinaryExpr) []token.Token {
	v.markVisited(be)
	var seq []token.Token
	if xbe, ok := be.X.(*ast.BinaryExpr); ok {
		xseq := v.extractOpSequence(xbe)
		seq = append(seq, xseq...)
	}
	seq = append(seq, be.Op)
	if ybe, ok := be.Y.(*ast.BinaryExpr); ok {
		yseq := v.extractOpSequence(ybe)
		seq = append(seq, yseq...)
	}
	return seq
}

// Visit implements the ast.Visitor interface.
func (v *complexityVisitor) Visit(n ast.Node) ast.Visitor {
	if n == nil {
		return nil
	}
	switch n := n.(type) {
	case *ast.ForStmt:
		return v.walkForStmt(n)
	case *ast.RangeStmt:
		return v.walkRangeStmt(n)
	case *ast.IfStmt:
		return v.walkIfStmt(n)
	case *ast.SwitchStmt:
		return v.walkSwitchStmt(n)
	case *ast.TypeSwitchStmt:
		return v.walkTypeSwitchStmt(n)
	case *ast.SelectStmt:
		return v.walkSelectStmt(n)
	case *ast.FuncLit:
		return v.walkFuncLit(n)
	case *ast.BranchStmt:
		return v.walkBranchStmt(n)
	case *ast.BinaryExpr:
		return v.walkBinaryExpr(n)
	case *ast.CallExpr:
		return v.walkCallExpr(n)
	}
	return v
}

func (v *complexityVisitor) walkForStmt(n *ast.ForStmt) ast.Visitor {
	v.nestIncComplexity()

	ast.Walk(v, n.Init)
	ast.Walk(v, n.Cond)
	ast.Walk(v, n.Post)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkRangeStmt(n *ast.RangeStmt) ast.Visitor {
	v.nestIncComplexity()

	ast.Walk(v, n.Key)
	ast.Walk(v, n.Value)
	ast.Walk(v, n.X)
	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkIfStmt(n *ast.IfStmt) ast.Visitor {
	if v.justAfterElse {
		v.incComplexity()
	} else {
		v.nestIncComplexity()
	}

	ast.Walk(v, n.Init)
	ast.Walk(v, n.Cond)
	if v.justAfterElse {
		// else if
		ast.Walk(v, n.Body)
	} else {
		v.incNesting()
		ast.Walk(v, n.Body)
		v.decNesting()
	}

	// branch here to decrease the number of types of nodes to handle in Visit func
	if _, ok := n.Else.(*ast.BlockStmt); ok {
		v.incComplexity()
		ast.Walk(v, n.Else)
	} else if _, ok := n.Else.(*ast.IfStmt); ok {
		v.justAfterElse = true
		ast.Walk(v, n.Else)
		v.justAfterElse = false
	}

	return nil
}

func (v *complexityVisitor) walkSwitchStmt(n *ast.SwitchStmt) ast.Visitor {
	v.nestIncComplexity()

	ast.Walk(v, n.Init)
	ast.Walk(v, n.Tag)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkTypeSwitchStmt(n *ast.TypeSwitchStmt) ast.Visitor {
	v.nestIncComplexity()

	ast.Walk(v, n.Init)
	ast.Walk(v, n.Assign)

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkSelectStmt(n *ast.SelectStmt) ast.Visitor {
	v.nestIncComplexity()

	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkFuncLit(n *ast.FuncLit) ast.Visitor {
	v.incNesting()
	ast.Walk(v, n.Body)
	v.decNesting()
	return nil
}

func (v *complexityVisitor) walkBranchStmt(n *ast.BranchStmt) ast.Visitor {
	if n.Label != nil {
		v.incComplexity()
	}
	// need to traverse all children
	return v
}

func (v *complexityVisitor) walkBinaryExpr(n *ast.BinaryExpr) ast.Visitor {
	if n.Op != token.LAND && n.Op != token.LOR {
		return v
	}
	if !v.isVisited(n) {
		ops := v.extractOpSequence(n)

		var currentOp token.Token
		for _, op := range ops {
			if currentOp != op {
				v.incComplexity()
				currentOp = op
			}
		}
	}
	return v
}

func (v *complexityVisitor) walkCallExpr(n *ast.CallExpr) ast.Visitor {
	if callIdent, ok := n.Fun.(*ast.Ident); ok {
		obj, name := callIdent.Obj, callIdent.Name
		if obj == v.name.Obj && name == v.name.Name {
			// called by same function directly (direct recursion)
			v.incComplexity()
		}
	}
	return v
}
