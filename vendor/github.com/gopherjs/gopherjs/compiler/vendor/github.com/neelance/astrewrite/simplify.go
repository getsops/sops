package astrewrite

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
)

type simplifyContext struct {
	info          *types.Info
	varCounter    int
	simplifyCalls bool
}

func Simplify(file *ast.File, info *types.Info, simplifyCalls bool) *ast.File {
	c := &simplifyContext{info: info, simplifyCalls: simplifyCalls}

	decls := make([]ast.Decl, len(file.Decls))
	for i, decl := range file.Decls {
		c.varCounter = 0
		switch decl := decl.(type) {
		case *ast.GenDecl:
			decls[i] = c.simplifyGenDecl(nil, decl)

		case *ast.FuncDecl:
			decls[i] = &ast.FuncDecl{
				Doc:  decl.Doc,
				Recv: decl.Recv,
				Name: decl.Name,
				Type: decl.Type,
				Body: c.simplifyBlock(decl.Body),
			}
		}
	}

	newFile := &ast.File{
		Doc:        file.Doc,
		Package:    file.Package,
		Name:       file.Name,
		Decls:      decls,
		Scope:      file.Scope,
		Imports:    file.Imports,
		Unresolved: file.Unresolved,
		Comments:   file.Comments,
	}
	c.info.Scopes[newFile] = c.info.Scopes[file]
	return newFile
}

func (c *simplifyContext) simplifyStmtList(stmts []ast.Stmt) []ast.Stmt {
	var newStmts []ast.Stmt
	for _, s := range stmts {
		c.simplifyStmt(&newStmts, s)
	}
	return newStmts
}

func (c *simplifyContext) simplifyGenDecl(stmts *[]ast.Stmt, decl *ast.GenDecl) *ast.GenDecl {
	if decl.Tok != token.VAR {
		return decl
	}

	specs := make([]ast.Spec, len(decl.Specs))
	for j, spec := range decl.Specs {
		switch spec := spec.(type) {
		case *ast.ValueSpec:
			var values []ast.Expr
			if spec.Values != nil {
				values = make([]ast.Expr, len(spec.Values))
				for i, v := range spec.Values {
					v2 := c.simplifyExpr(stmts, v)
					for _, initializer := range c.info.InitOrder {
						if initializer.Rhs == v {
							initializer.Rhs = v2
						}
					}
					values[i] = v2
				}
			}
			specs[j] = &ast.ValueSpec{
				Doc:     spec.Doc,
				Names:   spec.Names,
				Type:    spec.Type,
				Values:  values,
				Comment: spec.Comment,
			}
		default:
			specs[j] = spec
		}
	}

	return &ast.GenDecl{
		Doc:    decl.Doc,
		TokPos: decl.TokPos,
		Tok:    token.VAR,
		Lparen: decl.Lparen,
		Specs:  specs,
		Rparen: decl.Rparen,
	}
}

func (c *simplifyContext) simplifyStmt(stmts *[]ast.Stmt, s ast.Stmt) {
	if s == nil {
		return
	}

	switch s := s.(type) {
	case *ast.ExprStmt:
		*stmts = append(*stmts, &ast.ExprStmt{
			X: c.simplifyExpr2(stmts, s.X, true),
		})

	case *ast.BlockStmt:
		*stmts = append(*stmts, c.simplifyBlock(s))

	case *ast.LabeledStmt:
		c.simplifyStmt(stmts, s.Stmt)
		(*stmts)[len(*stmts)-1] = &ast.LabeledStmt{
			Label: s.Label,
			Colon: s.Colon,
			Stmt:  (*stmts)[len(*stmts)-1],
		}

	case *ast.AssignStmt:
		lhs := make([]ast.Expr, len(s.Lhs))
		for i, x := range s.Lhs {
			lhs[i] = c.simplifyExpr(stmts, x)
		}
		rhs := make([]ast.Expr, len(s.Rhs))
		for i, x := range s.Rhs {
			rhs[i] = c.simplifyExpr2(stmts, x, true)
		}
		*stmts = append(*stmts, &ast.AssignStmt{
			Lhs:    lhs,
			Tok:    s.Tok,
			TokPos: s.TokPos,
			Rhs:    rhs,
		})

	case *ast.DeclStmt:
		*stmts = append(*stmts, &ast.DeclStmt{
			Decl: c.simplifyGenDecl(stmts, s.Decl.(*ast.GenDecl)),
		})

	case *ast.IfStmt:
		if s.Init != nil {
			block := &ast.BlockStmt{}
			*stmts = append(*stmts, block)
			stmts = &block.List
			c.simplifyStmt(stmts, s.Init)
		}
		newS := &ast.IfStmt{
			If:   s.If,
			Cond: c.simplifyExpr(stmts, s.Cond),
			Body: c.simplifyBlock(s.Body),
			Else: c.toElseBranch(c.simplifyToStmtList(s.Else), c.info.Scopes[s.Else]),
		}
		c.info.Scopes[newS] = c.info.Scopes[s]
		*stmts = append(*stmts, newS)

	case *ast.SwitchStmt:
		c.simplifySwitch(stmts, s)

	case *ast.TypeSwitchStmt:
		if s.Init != nil {
			block := &ast.BlockStmt{}
			*stmts = append(*stmts, block)
			stmts = &block.List
			c.simplifyStmt(stmts, s.Init)
		}
		var assign ast.Stmt
		switch a := s.Assign.(type) {
		case *ast.ExprStmt:
			ta := a.X.(*ast.TypeAssertExpr)
			assign = &ast.ExprStmt{
				X: &ast.TypeAssertExpr{
					X:      c.simplifyExpr(stmts, ta.X),
					Lparen: ta.Lparen,
					Type:   ta.Type,
					Rparen: ta.Rparen,
				},
			}
		case *ast.AssignStmt:
			ta := a.Rhs[0].(*ast.TypeAssertExpr)
			assign = &ast.AssignStmt{
				Lhs:    a.Lhs,
				Tok:    a.Tok,
				TokPos: a.TokPos,
				Rhs: []ast.Expr{
					&ast.TypeAssertExpr{
						X:      c.simplifyExpr(stmts, ta.X),
						Lparen: ta.Lparen,
						Type:   ta.Type,
						Rparen: ta.Rparen,
					},
				},
			}
		default:
			panic("unexpected type switch assign")
		}
		clauses := make([]ast.Stmt, len(s.Body.List))
		for i, ccs := range s.Body.List {
			cc := ccs.(*ast.CaseClause)
			newClause := &ast.CaseClause{
				Case:  cc.Case,
				List:  cc.List,
				Colon: cc.Colon,
				Body:  c.simplifyStmtList(cc.Body),
			}
			if implicit, ok := c.info.Implicits[cc]; ok {
				c.info.Implicits[newClause] = implicit
			}
			clauses[i] = newClause
		}
		newS := &ast.TypeSwitchStmt{
			Switch: s.Switch,
			Assign: assign,
			Body: &ast.BlockStmt{
				List: clauses,
			},
		}
		c.info.Scopes[newS] = c.info.Scopes[s]
		*stmts = append(*stmts, newS)

	case *ast.ForStmt:
		newS := &ast.ForStmt{
			For:  s.For,
			Init: s.Init,
			Cond: s.Cond,
			Post: s.Post,
			Body: c.simplifyBlock(s.Body),
		}
		c.info.Scopes[newS] = c.info.Scopes[s]
		*stmts = append(*stmts, newS)

	// case *ast.ForStmt:
	// 	c.simplifyStmt(stmts, s.Init)
	// 	var condStmts []ast.Stmt
	// 	cond := c.newVar(&condStmts, s.Cond)
	// 	bodyStmts := s.Body.List
	// 	if len(condStmts) != 0 {
	// 		bodyStmts = append(append(condStmts, &ast.IfStmt{
	// 			Cond: &ast.UnaryExpr{
	// 				Op: token.NOT,
	// 				X:  cond,
	// 			},
	// 			Body: &ast.BlockStmt{
	// 				List: []ast.Stmt{&ast.BranchStmt{
	// 					Tok: token.BREAK,
	// 				}},
	// 			},
	// 		}), bodyStmts...)
	// 		cond = nil
	// 	}
	// 	*stmts = append(*stmts, &ast.ForStmt{
	// 		For:  s.For,
	// 		Cond: cond,
	// 		Post: s.Post,
	// 		Body: &ast.BlockStmt{
	// 			List: bodyStmts,
	// 		},
	// 	})

	case *ast.RangeStmt:
		var newS ast.Stmt
		switch t := c.info.TypeOf(s.X).Underlying().(type) {
		case *types.Chan:
			key := s.Key
			tok := s.Tok
			if key == nil {
				key = ast.NewIdent("_")
				tok = token.DEFINE
			}
			okVar := c.newIdent(types.Typ[types.Bool])
			if s.Tok == token.ASSIGN {
				*stmts = append(*stmts, &ast.DeclStmt{
					Decl: &ast.GenDecl{
						Tok: token.VAR,
						Specs: []ast.Spec{&ast.ValueSpec{
							Names: []*ast.Ident{okVar},
							Type:  ast.NewIdent("bool"),
						}},
					},
				})
			}
			newS = &ast.ForStmt{
				For: s.For,
				Body: &ast.BlockStmt{
					Lbrace: s.Body.Lbrace,
					List: append([]ast.Stmt{
						&ast.AssignStmt{
							Lhs:    []ast.Expr{key, okVar},
							TokPos: s.TokPos,
							Tok:    tok,
							Rhs: []ast.Expr{c.setType(&ast.UnaryExpr{
								Op: token.ARROW,
								X:  c.newVar(stmts, s.X),
							}, types.NewTuple(
								types.NewVar(token.NoPos, nil, "", t.Elem()),
								types.NewVar(token.NoPos, nil, "", types.Typ[types.Bool]),
							))},
						},
						&ast.IfStmt{
							Cond: c.setType(&ast.UnaryExpr{
								Op: token.NOT,
								X:  okVar,
							}, types.Typ[types.Bool]),
							Body: &ast.BlockStmt{
								List: []ast.Stmt{
									&ast.BranchStmt{Tok: token.BREAK},
								},
							},
						},
					}, c.simplifyStmtList(s.Body.List)...),
					Rbrace: s.Body.Rbrace,
				},
			}

		default:
			newS = &ast.RangeStmt{
				For:    s.For,
				Key:    s.Key,
				Value:  s.Value,
				TokPos: s.TokPos,
				Tok:    s.Tok,
				X:      s.X,
				Body:   c.simplifyBlock(s.Body),
			}
		}
		c.info.Scopes[newS] = c.info.Scopes[s]
		*stmts = append(*stmts, newS)

	case *ast.IncDecStmt:
		*stmts = append(*stmts, &ast.IncDecStmt{
			X:      c.simplifyExpr(stmts, s.X),
			TokPos: s.TokPos,
			Tok:    s.Tok,
		})

	case *ast.GoStmt:
		*stmts = append(*stmts, &ast.GoStmt{
			Go:   s.Go,
			Call: c.simplifyCall(stmts, s.Call),
		})

	case *ast.SelectStmt:
		clauses := make([]ast.Stmt, len(s.Body.List))
		for i, entry := range s.Body.List {
			cc := entry.(*ast.CommClause)
			var newComm ast.Stmt
			var bodyPrefix []ast.Stmt
			switch comm := cc.Comm.(type) {
			case *ast.ExprStmt:
				recv := comm.X.(*ast.UnaryExpr)
				if recv.Op != token.ARROW {
					panic("unexpected comm clause")
				}
				newComm = &ast.ExprStmt{
					X: &ast.UnaryExpr{
						Op:    token.ARROW,
						OpPos: recv.OpPos,
						X:     c.simplifyExpr(stmts, recv.X),
					},
				}
			case *ast.AssignStmt:
				recv := comm.Rhs[0].(*ast.UnaryExpr)
				if recv.Op != token.ARROW {
					panic("unexpected comm clause")
				}
				simplifyLhs := false
				for _, x := range comm.Lhs {
					if c.simplifyCalls && ContainsCall(x) {
						simplifyLhs = true
					}
				}
				lhs := comm.Lhs
				tok := comm.Tok
				if simplifyLhs {
					for i, x := range lhs {
						id := c.newIdent(c.info.TypeOf(x))
						bodyPrefix = append(bodyPrefix, simpleAssign(c.simplifyExpr(&bodyPrefix, x), comm.Tok, id))
						lhs[i] = id
					}
					tok = token.DEFINE
				}
				newComm = &ast.AssignStmt{
					Lhs: lhs,
					Tok: tok,
					Rhs: []ast.Expr{c.simplifyExpr(stmts, recv)},
				}
			case *ast.SendStmt:
				newComm = &ast.SendStmt{
					Chan:  c.simplifyExpr(stmts, comm.Chan),
					Arrow: comm.Arrow,
					Value: c.simplifyExpr(stmts, comm.Value),
				}
			case nil:
				newComm = nil
			default:
				panic("unexpected comm clause")
			}
			newCC := &ast.CommClause{
				Case:  cc.Case,
				Comm:  newComm,
				Colon: cc.Colon,
				Body:  append(bodyPrefix, c.simplifyStmtList(cc.Body)...),
			}
			c.info.Scopes[newCC] = c.info.Scopes[cc]
			clauses[i] = newCC
		}
		*stmts = append(*stmts, &ast.SelectStmt{
			Select: s.Select,
			Body: &ast.BlockStmt{
				List: clauses,
			},
		})

	case *ast.DeferStmt:
		*stmts = append(*stmts, &ast.DeferStmt{
			Defer: s.Defer,
			Call:  c.simplifyCall(stmts, s.Call),
		})

	case *ast.SendStmt:
		*stmts = append(*stmts, &ast.SendStmt{
			Chan:  c.simplifyExpr(stmts, s.Chan),
			Arrow: s.Arrow,
			Value: c.simplifyExpr(stmts, s.Value),
		})

	case *ast.ReturnStmt:
		*stmts = append(*stmts, &ast.ReturnStmt{
			Return:  s.Return,
			Results: c.simplifyExprList(stmts, s.Results),
		})

	default:
		*stmts = append(*stmts, s)
	}
}

func (c *simplifyContext) simplifyBlock(s *ast.BlockStmt) *ast.BlockStmt {
	if s == nil {
		return nil
	}
	newS := &ast.BlockStmt{
		Lbrace: s.Lbrace,
		List:   c.simplifyStmtList(s.List),
		Rbrace: s.Rbrace,
	}
	c.info.Scopes[newS] = c.info.Scopes[s]
	return newS
}

func (c *simplifyContext) simplifySwitch(stmts *[]ast.Stmt, s *ast.SwitchStmt) {
	wrapClause := &ast.CaseClause{}
	newS := &ast.SwitchStmt{
		Switch: s.Switch,
		Body:   &ast.BlockStmt{List: []ast.Stmt{wrapClause}},
	}
	c.info.Scopes[newS] = c.info.Scopes[s]
	c.info.Scopes[wrapClause] = c.info.Scopes[s]
	*stmts = append(*stmts, newS)
	stmts = &wrapClause.Body

	c.simplifyStmt(stmts, s.Init)

	nonDefaultClauses, defaultClause := c.simplifyCaseClauses(s.Body.List)
	tag := c.makeTag(stmts, s.Tag, len(nonDefaultClauses) != 0)
	*stmts = append(*stmts, unwrapBlock(c.switchToIfElse(tag, nonDefaultClauses, defaultClause))...)
}

func (c *simplifyContext) makeTag(stmts *[]ast.Stmt, tag ast.Expr, needsTag bool) ast.Expr {
	if tag == nil {
		id := ast.NewIdent("true")
		c.info.Types[id] = types.TypeAndValue{Type: types.Typ[types.Bool], Value: constant.MakeBool(true)}
		return id
	}
	if !needsTag {
		*stmts = append(*stmts, simpleAssign(ast.NewIdent("_"), token.ASSIGN, tag))
		return nil
	}
	return c.newVar(stmts, tag)
}

func (c *simplifyContext) simplifyCaseClauses(clauses []ast.Stmt) (nonDefaultClauses []*ast.CaseClause, defaultClause *ast.CaseClause) {
	var openClauses []*ast.CaseClause
	for _, cc := range clauses {
		clause := cc.(*ast.CaseClause)
		newClause := &ast.CaseClause{
			Case:  clause.Case,
			List:  clause.List,
			Colon: clause.Colon,
		}
		c.info.Scopes[newClause] = c.info.Scopes[clause]

		body := clause.Body
		hasFallthrough := false
		if len(body) != 0 {
			if b, isBranchStmt := body[len(body)-1].(*ast.BranchStmt); isBranchStmt && b.Tok == token.FALLTHROUGH {
				body = body[:len(body)-1]
				hasFallthrough = true
			}
		}
		openClauses = append(openClauses, newClause)
		for _, openClause := range openClauses {
			openClause.Body = append(openClause.Body, body...)
		}
		if !hasFallthrough {
			openClauses = nil
		}

		if len(clause.List) == 0 {
			defaultClause = newClause
			continue
		}
		nonDefaultClauses = append(nonDefaultClauses, newClause)
	}
	return
}

func (c *simplifyContext) switchToIfElse(tag ast.Expr, nonDefaultClauses []*ast.CaseClause, defaultClause *ast.CaseClause) ast.Stmt {
	if len(nonDefaultClauses) == 0 {
		if defaultClause != nil {
			return c.toElseBranch(c.simplifyStmtList(defaultClause.Body), c.info.Scopes[defaultClause])
		}
		return nil
	}

	clause := nonDefaultClauses[0]
	conds := make([]ast.Expr, len(clause.List))
	for i, cond := range clause.List {
		conds[i] = c.setType(&ast.BinaryExpr{
			X:  tag,
			Op: token.EQL,
			Y:  c.setType(&ast.ParenExpr{X: cond}, c.info.TypeOf(cond)),
		}, types.Typ[types.Bool])
	}

	var stmts []ast.Stmt
	ifStmt := &ast.IfStmt{
		If:   clause.Case,
		Cond: c.simplifyExpr(&stmts, c.disjunction(conds)),
		Body: &ast.BlockStmt{List: c.simplifyStmtList(clause.Body)},
		Else: c.switchToIfElse(tag, nonDefaultClauses[1:], defaultClause),
	}
	c.info.Scopes[ifStmt] = c.info.Scopes[clause]
	stmts = append(stmts, ifStmt)
	return c.toElseBranch(stmts, c.info.Scopes[clause])
}

func (c *simplifyContext) disjunction(conds []ast.Expr) ast.Expr {
	if len(conds) == 1 {
		return conds[0]
	}
	return c.setType(&ast.BinaryExpr{
		X:  conds[0],
		Op: token.LOR,
		Y:  c.disjunction(conds[1:]),
	}, types.Typ[types.Bool])
}

func (c *simplifyContext) simplifyToStmtList(s ast.Stmt) (stmts []ast.Stmt) {
	c.simplifyStmt(&stmts, s)
	return
}

func (c *simplifyContext) toElseBranch(stmts []ast.Stmt, scope *types.Scope) ast.Stmt {
	if len(stmts) == 0 {
		return nil
	}
	if len(stmts) == 1 {
		switch stmt := stmts[0].(type) {
		case *ast.IfStmt, *ast.BlockStmt:
			c.info.Scopes[stmt] = scope
			return stmt
		}
	}
	block := &ast.BlockStmt{
		List: stmts,
	}
	c.info.Scopes[block] = scope
	return block
}

func unwrapBlock(s ast.Stmt) []ast.Stmt {
	if s == nil {
		return nil
	}
	if block, ok := s.(*ast.BlockStmt); ok {
		return block.List
	}
	return []ast.Stmt{s}
}

func (c *simplifyContext) simplifyExpr(stmts *[]ast.Stmt, x ast.Expr) ast.Expr {
	return c.simplifyExpr2(stmts, x, false)
}

func (c *simplifyContext) simplifyExpr2(stmts *[]ast.Stmt, x ast.Expr, callOK bool) ast.Expr {
	x2 := c.simplifyExpr3(stmts, x, callOK)
	if t, ok := c.info.Types[x]; ok {
		c.info.Types[x2] = t
	}
	return x2
}

func (c *simplifyContext) simplifyExpr3(stmts *[]ast.Stmt, x ast.Expr, callOK bool) ast.Expr {
	switch x := x.(type) {
	case *ast.FuncLit:
		return &ast.FuncLit{
			Type: x.Type,
			Body: &ast.BlockStmt{
				List: c.simplifyStmtList(x.Body.List),
			},
		}

	case *ast.CompositeLit:
		elts := make([]ast.Expr, len(x.Elts))
		for i, elt := range x.Elts {
			if kv, ok := elt.(*ast.KeyValueExpr); ok {
				elts[i] = &ast.KeyValueExpr{
					Key:   kv.Key,
					Colon: kv.Colon,
					Value: c.simplifyExpr(stmts, kv.Value),
				}
				continue
			}
			elts[i] = c.simplifyExpr(stmts, elt)
		}
		return &ast.CompositeLit{
			Type:   x.Type,
			Lbrace: x.Lbrace,
			Elts:   elts,
			Rbrace: x.Rbrace,
		}

	case *ast.ParenExpr:
		return &ast.ParenExpr{
			Lparen: x.Lparen,
			X:      c.simplifyExpr(stmts, x.X),
			Rparen: x.Rparen,
		}

	case *ast.SelectorExpr:
		selExpr := &ast.SelectorExpr{
			X:   c.simplifyExpr(stmts, x.X),
			Sel: x.Sel,
		}
		if sel, ok := c.info.Selections[x]; ok {
			c.info.Selections[selExpr] = sel
		}
		return selExpr

	case *ast.IndexExpr:
		return &ast.IndexExpr{
			X:      c.simplifyExpr(stmts, x.X),
			Lbrack: x.Lbrack,
			Index:  c.simplifyExpr(stmts, x.Index),
			Rbrack: x.Rbrack,
		}

	case *ast.SliceExpr:
		return &ast.SliceExpr{
			X:      c.simplifyExpr(stmts, x.X),
			Lbrack: x.Lbrack,
			Low:    c.simplifyExpr(stmts, x.Low),
			High:   c.simplifyExpr(stmts, x.High),
			Max:    c.simplifyExpr(stmts, x.Max),
			Slice3: x.Slice3,
			Rbrack: x.Rbrack,
		}

	case *ast.TypeAssertExpr:
		return &ast.TypeAssertExpr{
			X:      c.simplifyExpr(stmts, x.X),
			Lparen: x.Lparen,
			Type:   x.Type,
			Rparen: x.Rparen,
		}

	case *ast.CallExpr:
		call := c.simplifyCall(stmts, x)
		if callOK || !c.simplifyCalls {
			return call
		}
		return c.newVar(stmts, call)

	case *ast.StarExpr:
		return &ast.StarExpr{
			Star: x.Star,
			X:    c.simplifyExpr(stmts, x.X),
		}

	case *ast.UnaryExpr:
		return &ast.UnaryExpr{
			OpPos: x.OpPos,
			Op:    x.Op,
			X:     c.simplifyExpr(stmts, x.X),
		}

	case *ast.BinaryExpr:
		if (x.Op == token.LAND || x.Op == token.LOR) && c.simplifyCalls && ContainsCall(x.Y) {
			v := c.newVar(stmts, x.X)
			cond := v
			if x.Op == token.LOR {
				cond = &ast.UnaryExpr{
					Op: token.NOT,
					X:  cond,
				}
			}
			var ifBody []ast.Stmt
			ifBody = append(ifBody, simpleAssign(v, token.ASSIGN, c.simplifyExpr2(&ifBody, x.Y, true)))
			*stmts = append(*stmts, &ast.IfStmt{
				Cond: cond,
				Body: &ast.BlockStmt{
					List: ifBody,
				},
			})
			return v
		}
		return &ast.BinaryExpr{
			X:     c.simplifyExpr(stmts, x.X),
			OpPos: x.OpPos,
			Op:    x.Op,
			Y:     c.simplifyExpr(stmts, x.Y),
		}

	default:
		return x
	}
}

func (c *simplifyContext) simplifyCall(stmts *[]ast.Stmt, x *ast.CallExpr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:      c.simplifyExpr(stmts, x.Fun),
		Lparen:   x.Lparen,
		Args:     c.simplifyArgs(stmts, x.Args),
		Ellipsis: x.Ellipsis,
		Rparen:   x.Rparen,
	}
}

func (c *simplifyContext) simplifyArgs(stmts *[]ast.Stmt, args []ast.Expr) []ast.Expr {
	if len(args) == 1 {
		if tuple, ok := c.info.TypeOf(args[0]).(*types.Tuple); ok && c.simplifyCalls {
			call := c.simplifyExpr2(stmts, args[0], true)
			vars := make([]ast.Expr, tuple.Len())
			for i := range vars {
				vars[i] = c.newIdent(tuple.At(i).Type())
			}
			*stmts = append(*stmts, &ast.AssignStmt{
				Lhs: vars,
				Tok: token.DEFINE,
				Rhs: []ast.Expr{call},
			})
			return vars
		}
	}
	return c.simplifyExprList(stmts, args)
}

func (c *simplifyContext) simplifyExprList(stmts *[]ast.Stmt, exprs []ast.Expr) []ast.Expr {
	if exprs == nil {
		return nil
	}
	simplifiedExprs := make([]ast.Expr, len(exprs))
	for i, expr := range exprs {
		simplifiedExprs[i] = c.simplifyExpr(stmts, expr)
	}
	return simplifiedExprs
}

func (c *simplifyContext) newVar(stmts *[]ast.Stmt, x ast.Expr) ast.Expr {
	id := c.newIdent(c.info.TypeOf(x))
	*stmts = append(*stmts, simpleAssign(id, token.DEFINE, x))
	return id
}

func (c *simplifyContext) newIdent(t types.Type) *ast.Ident {
	c.varCounter++
	id := ast.NewIdent(fmt.Sprintf("_%d", c.varCounter))
	c.info.Types[id] = types.TypeAndValue{Type: t} // TODO remove?
	c.info.Uses[id] = types.NewVar(token.NoPos, nil, id.Name, t)
	return id
}

func (c *simplifyContext) setType(x ast.Expr, t types.Type) ast.Expr {
	c.info.Types[x] = types.TypeAndValue{Type: t}
	return x
}

func simpleAssign(lhs ast.Expr, tok token.Token, rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{lhs},
		Tok: tok,
		Rhs: []ast.Expr{rhs},
	}
}

func ContainsCall(x ast.Expr) bool {
	switch x := x.(type) {
	case *ast.CallExpr:
		return true
	case *ast.CompositeLit:
		for _, elt := range x.Elts {
			if ContainsCall(elt) {
				return true
			}
		}
		return false
	case *ast.KeyValueExpr:
		return ContainsCall(x.Key) || ContainsCall(x.Value)
	case *ast.ParenExpr:
		return ContainsCall(x.X)
	case *ast.SelectorExpr:
		return ContainsCall(x.X)
	case *ast.IndexExpr:
		return ContainsCall(x.X) || ContainsCall(x.Index)
	case *ast.SliceExpr:
		return ContainsCall(x.X) || ContainsCall(x.Low) || ContainsCall(x.High) || ContainsCall(x.Max)
	case *ast.TypeAssertExpr:
		return ContainsCall(x.X)
	case *ast.StarExpr:
		return ContainsCall(x.X)
	case *ast.UnaryExpr:
		return ContainsCall(x.X)
	case *ast.BinaryExpr:
		return ContainsCall(x.X) || ContainsCall(x.Y)
	default:
		return false
	}
}
