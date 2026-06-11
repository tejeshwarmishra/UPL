package parser

import (
	"fmt"
	"strings"

	"github.com/tejeshwarmishra/upl/pkg/lexer"
)

// ASTNode is the interface all AST nodes must implement
type ASTNode interface {
	ToJavaScript() string
	Line() int
}

// Block represents a sequence of statements
type Block struct {
	Statements []ASTNode
	LineNum    int
}

func (b *Block) ToJavaScript() string {
	var out strings.Builder
	for _, stmt := range b.Statements {
		out.WriteString(stmt.ToJavaScript())
		out.WriteString("\n")
	}
	return out.String()
}

func (b *Block) Line() int { return b.LineNum }

// ConstStmt represents: स्थिर बनाओ x = value:
type ConstStmt struct {
	Name     string
	Value    ASTNode
	LineNum  int
}

func (c *ConstStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("const ")
	out.WriteString(c.Name)
	out.WriteString(" = ")
	out.WriteString(c.Value.ToJavaScript())
	out.WriteString(";")
	return out.String()
}

func (c *ConstStmt) Line() int { return c.LineNum }

// LetStmt represents: बदलाव योग्य x = value:
type LetStmt struct {
	Name     string
	Value    ASTNode
	LineNum  int
}

func (l *LetStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("let ")
	out.WriteString(l.Name)
	out.WriteString(" = ")
	out.WriteString(l.Value.ToJavaScript())
	out.WriteString(";")
	return out.String()
}

func (l *LetStmt) Line() int { return l.LineNum }

// VarStmt represents: अपरिवर्तनीय x = value:
type VarStmt struct {
	Name     string
	Value    ASTNode
	LineNum  int
}

func (v *VarStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("var ")
	out.WriteString(v.Name)
	out.WriteString(" = ")
	out.WriteString(v.Value.ToJavaScript())
	out.WriteString(";")
	return out.String()
}

func (v *VarStmt) Line() int { return v.LineNum }

// FuncDecl represents: कार्य बनाओ funcName(params):
type FuncDecl struct {
	Name       string
	Params     []string
	Body       *Block
	IsAsync    bool
	LineNum    int
}

func (f *FuncDecl) ToJavaScript() string {
	var out strings.Builder
	if f.IsAsync {
		out.WriteString("async ")
	}
	out.WriteString("function ")
	out.WriteString(f.Name)
	out.WriteString("(")
	out.WriteString(strings.Join(f.Params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.ToJavaScript())
	out.WriteString("}")
	return out.String()
}

func (f *FuncDecl) Line() int { return f.LineNum }

// ReturnStmt represents: वापस करो value:
type ReturnStmt struct {
	Value   ASTNode
	LineNum int
}

func (r *ReturnStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("return ")
	if r.Value != nil {
		out.WriteString(r.Value.ToJavaScript())
	}
	out.WriteString(";")
	return out.String()
}

func (r *ReturnStmt) Line() int { return r.LineNum }

// LogStmt represents: दिखाओ(value):
type LogStmt struct {
	Args    []ASTNode
	LineNum int
}

func (l *LogStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("console.log(")
	for i, arg := range l.Args {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.ToJavaScript())
	}
	out.WriteString(");")
	return out.String()
}

func (l *LogStmt) Line() int { return l.LineNum }

// IfStmt represents: यदि condition:
type IfStmt struct {
	Condition ASTNode
	Then      *Block
	Else      *Block
	LineNum   int
}

func (i *IfStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("if (")
	out.WriteString(i.Condition.ToJavaScript())
	out.WriteString(") {\n")
	out.WriteString(i.Then.ToJavaScript())
	out.WriteString("}")
	if i.Else != nil {
		out.WriteString(" else {\n")
		out.WriteString(i.Else.ToJavaScript())
		out.WriteString("}")
	}
	return out.String()
}

func (i *IfStmt) Line() int { return i.LineNum }

// WhileStmt represents: जब तक condition:
type WhileStmt struct {
	Condition ASTNode
	Body      *Block
	LineNum   int
}

func (w *WhileStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("while (")
	out.WriteString(w.Condition.ToJavaScript())
	out.WriteString(") {\n")
	out.WriteString(w.Body.ToJavaScript())
	out.WriteString("}")
	return out.String()
}

func (w *WhileStmt) Line() int { return w.LineNum }

// ForStmt represents: के लिए (init; condition; update):
type ForStmt struct {
	Init      ASTNode
	Condition ASTNode
	Update    ASTNode
	Body      *Block
	LineNum   int
}

func (f *ForStmt) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("for (")
	if f.Init != nil {
		out.WriteString(f.Init.ToJavaScript())
	}
	out.WriteString("; ")
	if f.Condition != nil {
		out.WriteString(f.Condition.ToJavaScript())
	}
	out.WriteString("; ")
	if f.Update != nil {
		out.WriteString(f.Update.ToJavaScript())
	}
	out.WriteString(") {\n")
	out.WriteString(f.Body.ToJavaScript())
	out.WriteString("}")
	return out.String()
}

func (f *ForStmt) Line() int { return f.LineNum }

// ClassDecl represents: श्रेणी ClassName:
type ClassDecl struct {
	Name       string
	Methods    map[string]*FuncDecl
	Properties []string
	LineNum    int
}

func (c *ClassDecl) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("class ")
	out.WriteString(c.Name)
	out.WriteString(" {\n")
	for _, prop := range c.Properties {
		out.WriteString("  ")
		out.WriteString(prop)
		out.WriteString(";\n")
	}
	for _, method := range c.Methods {
		out.WriteString("  ")
		out.WriteString(method.ToJavaScript())
		out.WriteString("\n")
	}
	out.WriteString("}")
	return out.String()
}

func (c *ClassDecl) Line() int { return c.LineNum }

// BinaryOp represents: left operator right
type BinaryOp struct {
	Left    ASTNode
	Op      string
	Right   ASTNode
	LineNum int
}

func (b *BinaryOp) ToJavaScript() string {
	var out strings.Builder
	out.WriteString("(")
	out.WriteString(b.Left.ToJavaScript())
	out.WriteString(" ")
	out.WriteString(b.Op)
	out.WriteString(" ")
	out.WriteString(b.Right.ToJavaScript())
	out.WriteString(")")
	return out.String()
}

func (b *BinaryOp) Line() int { return b.LineNum }

// CallExpr represents: funcName(args)
type CallExpr struct {
	Function ASTNode
	Args     []ASTNode
	LineNum  int
}

func (c *CallExpr) ToJavaScript() string {
	var out strings.Builder
	out.WriteString(c.Function.ToJavaScript())
	out.WriteString("(")
	for i, arg := range c.Args {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.ToJavaScript())
	}
	out.WriteString(")")
	return out.String()
}

func (c *CallExpr) Line() int { return c.LineNum }

// Identifier represents a variable name
type Identifier struct {
	Name    string
	LineNum int
}

func (i *Identifier) ToJavaScript() string {
	return i.Name
}

func (i *Identifier) Line() int { return i.LineNum }

// StringLiteral represents a string value
type StringLiteral struct {
	Value   string
	LineNum int
}

func (s *StringLiteral) ToJavaScript() string {
	return fmt.Sprintf(`"%s"`, strings.ReplaceAll(s.Value, `"`, `\"`))
}

func (s *StringLiteral) Line() int { return s.LineNum }

// NumberLiteral represents a numeric value
type NumberLiteral struct {
	Value   string
	LineNum int
}

func (n *NumberLiteral) ToJavaScript() string {
	return n.Value
}

func (n *NumberLiteral) Line() int { return n.LineNum }

// StyleBlock represents a CSS style block
type StyleBlock struct {
	Selector   string
	Properties map[string]string
	LineNum    int
}

func (s *StyleBlock) ToCSS() string {
	var out strings.Builder
	out.WriteString(s.Selector)
	out.WriteString(" {\n")
	for key, value := range s.Properties {
		out.WriteString("  ")
		out.WriteString(key)
		out.WriteString(": ")
		out.WriteString(value)
		out.WriteString(";\n")
	}
	out.WriteString("}")
	return out.String()
}

func (s *StyleBlock) ToJavaScript() string {
	return ""
}

func (s *StyleBlock) Line() int { return s.LineNum }

// ParseError represents a parsing error with location info
type ParseError struct {
	Line       int
	Column     int
	Message    string
	Context    string
	Suggestion string
}

func (e *ParseError) Error() string {
	var out strings.Builder
	out.WriteString(fmt.Sprintf("Parse Error at %d:%d\n", e.Line, e.Column))
	out.WriteString(fmt.Sprintf("Message: %s\n", e.Message))
	if e.Context != "" {
		out.WriteString(fmt.Sprintf("Context: %s\n", e.Context))
	}
	if e.Suggestion != "" {
		out.WriteString(fmt.Sprintf("Suggestion: %s\n", e.Suggestion))
	}
	return out.String()
}

// Parser tokenizes and builds an AST
type Parser struct {
	tokens    []lexer.Token
	current   int
	errors    []*ParseError
}

// New creates a new parser from a lexer
func New(input string) *Parser {
	lex := lexer.New(input)
	tokens := lex.TokenizeAll()
	return &Parser{
		tokens: tokens,
		current: 0,
		errors: []*ParseError{},
	}
}

// Errors returns any parsing errors encountered
func (p *Parser) Errors() []*ParseError {
	return p.errors
}

// Parse builds the AST from tokens
func (p *Parser) Parse() *Block {
	block := &Block{
		Statements: []ASTNode{},
		LineNum:    1,
	}

	for !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

// parseStatement parses a single statement
func (p *Parser) parseStatement() ASTNode {
	if p.isAtEnd() {
		return nil
	}

	tok := p.peek()

	switch tok.Type {
	case lexer.TokenNewline:
		p.advance()
		return nil
	case lexer.TokenConst:
		return p.parseConst()
	case lexer.TokenLet:
		return p.parseLet()
	case lexer.TokenVar:
		return p.parseVar()
	case lexer.TokenFunction:
		return p.parseFunction()
	case lexer.TokenReturn:
		return p.parseReturn()
	case lexer.TokenLog:
		return p.parseLog()
	case lexer.TokenIf:
		return p.parseIf()
	case lexer.TokenWhile:
		return p.parseWhile()
	case lexer.TokenFor:
		return p.parseFor()
	case lexer.TokenClass:
		return p.parseClass()
	case lexer.TokenStyle, lexer.TokenColorHub:
		return p.parseStyle()
	default:
		stmt := p.parseExpression()
		p.skipNewlines()
		return stmt
	}
}

func (p *Parser) parseConst() *ConstStmt {
	line := p.peek().Line
	p.advance() // consume 'const'
	
	name := p.expectIdent()
	if name == "" {
		p.addError(line, p.peek().Column, "Expected identifier after 'const'", p.peek().Value, "Use: स्थिर बनाओ x = value:")
		return nil
	}

	if !p.match(lexer.TokenEquals) {
		p.addError(line, p.peek().Column, "Expected '=' after identifier", p.peek().Value, "Use: स्थिर बनाओ x = value:")
		return nil
	}

	value := p.parseExpression()
	p.skipToNewlineOrEnd()

	return &ConstStmt{
		Name:    name,
		Value:   value,
		LineNum: line,
	}
}

func (p *Parser) parseLet() *LetStmt {
	line := p.peek().Line
	p.advance()
	
	name := p.expectIdent()
	if name == "" {
		p.addError(line, p.peek().Column, "Expected identifier after 'let'", p.peek().Value, "Use: बदलाव योग्य x = value:")
		return nil
	}

	if !p.match(lexer.TokenEquals) {
		p.addError(line, p.peek().Column, "Expected '=' after identifier", p.peek().Value, "Use: बदलाव योग्य x = value:")
		return nil
	}

	value := p.parseExpression()
	p.skipToNewlineOrEnd()

	return &LetStmt{
		Name:    name,
		Value:   value,
		LineNum: line,
	}
}

func (p *Parser) parseVar() *VarStmt {
	line := p.peek().Line
	p.advance()
	
	name := p.expectIdent()
	if name == "" {
		p.addError(line, p.peek().Column, "Expected identifier after 'var'", p.peek().Value, "Use: अपरिवर्तनीय x = value:")
		return nil
	}

	if !p.match(lexer.TokenEquals) {
		p.addError(line, p.peek().Column, "Expected '=' after identifier", p.peek().Value, "Use: अपरिवर्तनीय x = value:")
		return nil
	}

	value := p.parseExpression()
	p.skipToNewlineOrEnd()

	return &VarStmt{
		Name:    name,
		Value:   value,
		LineNum: line,
	}
}

func (p *Parser) parseFunction() *FuncDecl {
	line := p.peek().Line
	p.advance() // consume 'function'

	name := p.expectIdent()
	if name == "" {
		p.addError(line, p.peek().Column, "Expected function name", p.peek().Value, "Use: कार्य बनाओ funcName():")
		return nil
	}

	if !p.match(lexer.TokenLeftParen) {
		p.addError(line, p.peek().Column, "Expected '(' after function name", p.peek().Value, "Use: कार्य बनाओ funcName():")
		return nil
	}

	params := []string{}
	if !p.check(lexer.TokenRightParen) {
		for {
			param := p.expectIdent()
			if param != "" {
				params = append(params, param)
			}
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if !p.match(lexer.TokenRightParen) {
		p.addError(line, p.peek().Column, "Expected ')' after parameters", p.peek().Value, "Use: कार्य बनाओ funcName():")
		return nil
	}

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after function signature", p.peek().Value, "Use: कार्य बनाओ funcName():")
		return nil
	}

	p.skipNewlines()
	body := p.parseBlock()

	return &FuncDecl{
		Name:    name,
		Params:  params,
		Body:    body,
		LineNum: line,
	}
}

func (p *Parser) parseReturn() *ReturnStmt {
	line := p.peek().Line
	p.advance()

	var value ASTNode
	if !p.check(lexer.TokenNewline) && !p.isAtEnd() {
		value = p.parseExpression()
	}

	return &ReturnStmt{
		Value:   value,
		LineNum: line,
	}
}

func (p *Parser) parseLog() *LogStmt {
	line := p.peek().Line
	p.advance()

	if !p.match(lexer.TokenLeftParen) {
		p.addError(line, p.peek().Column, "Expected '(' after 'log'", p.peek().Value, "Use: दिखाओ(value)")
		return nil
	}

	args := []ASTNode{}
	if !p.check(lexer.TokenRightParen) {
		for {
			args = append(args, p.parseExpression())
			if !p.match(lexer.TokenComma) {
				break
			}
		}
	}

	if !p.match(lexer.TokenRightParen) {
		p.addError(line, p.peek().Column, "Expected ')' after log arguments", p.peek().Value, "Use: दिखाओ(value)")
		return nil
	}

	return &LogStmt{
		Args:    args,
		LineNum: line,
	}
}

func (p *Parser) parseIf() *IfStmt {
	line := p.peek().Line
	p.advance()

	if !p.match(lexer.TokenLeftParen) {
		p.addError(line, p.peek().Column, "Expected '(' after 'if'", p.peek().Value, "Use: यदि (condition):")
		return nil
	}

	condition := p.parseExpression()

	if !p.match(lexer.TokenRightParen) {
		p.addError(line, p.peek().Column, "Expected ')' after condition", p.peek().Value, "Use: यदि (condition):")
		return nil
	}

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after if condition", p.peek().Value, "Use: यदि (condition):")
		return nil
	}

	p.skipNewlines()
	thenBlock := p.parseBlock()

	var elseBlock *Block
	if p.peek().Type == lexer.TokenElse {
		p.advance()
		if !p.match(lexer.TokenColon) {
			p.addError(line, p.peek().Column, "Expected ':' after 'else'", p.peek().Value, "Use: वरना:")
		}
		p.skipNewlines()
		elseBlock = p.parseBlock()
	}

	return &IfStmt{
		Condition: condition,
		Then:      thenBlock,
		Else:      elseBlock,
		LineNum:   line,
	}
}

func (p *Parser) parseWhile() *WhileStmt {
	line := p.peek().Line
	p.advance()

	if !p.match(lexer.TokenLeftParen) {
		p.addError(line, p.peek().Column, "Expected '(' after 'while'", p.peek().Value, "Use: जब तक (condition):")
		return nil
	}

	condition := p.parseExpression()

	if !p.match(lexer.TokenRightParen) {
		p.addError(line, p.peek().Column, "Expected ')' after condition", p.peek().Value, "Use: जब तक (condition):")
		return nil
	}

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after while condition", p.peek().Value, "Use: जब तक (condition):")
		return nil
	}

	p.skipNewlines()
	body := p.parseBlock()

	return &WhileStmt{
		Condition: condition,
		Body:      body,
		LineNum:   line,
	}
}

func (p *Parser) parseFor() *ForStmt {
	line := p.peek().Line
	p.advance()

	if !p.match(lexer.TokenLeftParen) {
		p.addError(line, p.peek().Column, "Expected '(' after 'for'", p.peek().Value, "Use: के लिए (init; condition; update):")
		return nil
	}

	var init ASTNode
	if !p.check(lexer.TokenSemicolon) {
		init = p.parseExpression()
	}
	p.match(lexer.TokenSemicolon)

	var condition ASTNode
	if !p.check(lexer.TokenSemicolon) {
		condition = p.parseExpression()
	}
	p.match(lexer.TokenSemicolon)

	var update ASTNode
	if !p.check(lexer.TokenRightParen) {
		update = p.parseExpression()
	}
	p.match(lexer.TokenRightParen)

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after for declaration", p.peek().Value, "Use: के लिए (init; condition; update):")
		return nil
	}

	p.skipNewlines()
	body := p.parseBlock()

	return &ForStmt{
		Init:      init,
		Condition: condition,
		Update:    update,
		Body:      body,
		LineNum:   line,
	}
}

func (p *Parser) parseClass() *ClassDecl {
	line := p.peek().Line
	p.advance()

	name := p.expectIdent()
	if name == "" {
		p.addError(line, p.peek().Column, "Expected class name", p.peek().Value, "Use: श्रेणी ClassName:")
		return nil
	}

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after class name", p.peek().Value, "Use: श्रेणी ClassName:")
		return nil
	}

	p.skipNewlines()
	body := p.parseBlock()

	// Convert body statements into methods
	methods := make(map[string]*FuncDecl)
	for _, stmt := range body.Statements {
		if fn, ok := stmt.(*FuncDecl); ok {
			methods[fn.Name] = fn
		}
	}

	return &ClassDecl{
		Name:    name,
		Methods: methods,
		LineNum: line,
	}
}

func (p *Parser) parseStyle() *StyleBlock {
	line := p.peek().Line
	selector := ""

	if p.peek().Type == lexer.TokenColorHub {
		selector = "रंग-मंच"
	} else if p.peek().Type == lexer.TokenStyle {
		selector = "style"
	}

	p.advance()

	if !p.match(lexer.TokenColon) {
		p.addError(line, p.peek().Column, "Expected ':' after style declaration", p.peek().Value, "Use: शैली:")
		return nil
	}

	p.skipNewlines()

	properties := make(map[string]string)

	for !p.isAtEnd() && p.peek().Type != lexer.TokenDedent {
		key := p.expectIdent()
		if key == "" {
			break
		}

		if !p.match(lexer.TokenEquals) {
			break
		}

		var value string
		if p.peek().Type == lexer.TokenString {
			value = strings.Trim(p.advance().Value, `"'`)
		} else {
			value = p.advance().Value
		}

		properties[key] = value
		p.skipNewlines()
	}

	return &StyleBlock{
		Selector:   selector,
		Properties: properties,
		LineNum:    line,
	}
}

func (p *Parser) parseBlock() *Block {
	block := &Block{
		Statements: []ASTNode{},
		LineNum:    p.peek().Line,
	}

	for !p.isAtEnd() && p.peek().Type != lexer.TokenDedent {
		if p.peek().Type == lexer.TokenNewline {
			p.advance()
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}

	return block
}

func (p *Parser) parseExpression() ASTNode {
	return p.parseComparison()
}

func (p *Parser) parseComparison() ASTNode {
	expr := p.parseAddition()

	for p.match(lexer.TokenDoubleEquals, lexer.TokenTripleEquals, lexer.TokenNotEquals) {
		op := p.previous().Value
		right := p.parseAddition()
		expr = &BinaryOp{
			Left:    expr,
			Op:      op,
			Right:   right,
			LineNum: expr.Line(),
		}
	}

	return expr
}

func (p *Parser) parseAddition() ASTNode {
	expr := p.parseMultiplication()

	for p.match(lexer.TokenPlus, lexer.TokenMinus) {
		op := p.previous().Value
		right := p.parseMultiplication()
		expr = &BinaryOp{
			Left:    expr,
			Op:      op,
			Right:   right,
			LineNum: expr.Line(),
		}
	}

	return expr
}

func (p *Parser) parseMultiplication() ASTNode {
	expr := p.parseCall()

	for p.match(lexer.TokenStar, lexer.TokenSlash, lexer.TokenPercent) {
		op := p.previous().Value
		right := p.parseCall()
		expr = &BinaryOp{
			Left:    expr,
			Op:      op,
			Right:   right,
			LineNum: expr.Line(),
		}
	}

	return expr
}

func (p *Parser) parseCall() ASTNode {
	expr := p.parsePrimary()

	for {
		if p.match(lexer.TokenLeftParen) {
			args := []ASTNode{}
			if !p.check(lexer.TokenRightParen) {
				for {
					args = append(args, p.parseExpression())
					if !p.match(lexer.TokenComma) {
						break
					}
				}
			}
			p.match(lexer.TokenRightParen)
			expr = &CallExpr{
				Function: expr,
				Args:     args,
				LineNum:  expr.Line(),
			}
		} else if p.match(lexer.TokenDot) {
			name := p.expectIdent()
			expr = &BinaryOp{
				Left:    expr,
				Op:      ".",
				Right:   &Identifier{Name: name, LineNum: p.peek().Line},
				LineNum: expr.Line(),
			}
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) parsePrimary() ASTNode {
	line := p.peek().Line

	if p.peek().Type == lexer.TokenString {
		value := strings.Trim(p.advance().Value, `"'`)
		return &StringLiteral{Value: value, LineNum: line}
	}

	if p.peek().Type == lexer.TokenNumber {
		value := p.advance().Value
		return &NumberLiteral{Value: value, LineNum: line}
	}

	if p.peek().Type == lexer.TokenIdent {
		name := p.advance().Value
		return &Identifier{Name: name, LineNum: line}
	}

	if p.match(lexer.TokenLeftParen) {
		expr := p.parseExpression()
		p.match(lexer.TokenRightParen)
		return expr
	}

	if p.peek().Type == lexer.TokenNot {
		p.advance()
		expr := p.parsePrimary()
		return &BinaryOp{
			Left:    &Identifier{Name: "!", LineNum: line},
			Op:      "!",
			Right:   expr,
			LineNum: line,
		}
	}

	p.addError(line, p.peek().Column, "Unexpected token", p.peek().Value, "Expected expression")
	return &Identifier{Name: "undefined", LineNum: line}
}

// Helper methods

func (p *Parser) peek() lexer.Token {
	if p.current >= len(p.tokens) {
		return lexer.Token{Type: lexer.TokenEOF, Value: ""}
	}
	return p.tokens[p.current]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens) || p.tokens[p.current].Type == lexer.TokenEOF
}

func (p *Parser) expectIdent() string {
	if p.peek().Type == lexer.TokenIdent {
		return p.advance().Value
	}
	return ""
}

func (p *Parser) skipNewlines() {
	for p.match(lexer.TokenNewline) {
	}
}

func (p *Parser) skipToNewlineOrEnd() {
	for !p.isAtEnd() && !p.check(lexer.TokenNewline) {
		p.advance()
	}
	p.skipNewlines()
}

func (p *Parser) addError(line, column int, message, context, suggestion string) {
	p.errors = append(p.errors, &ParseError{
		Line:       line,
		Column:     column,
		Message:    message,
		Context:    context,
		Suggestion: suggestion,
	})
}