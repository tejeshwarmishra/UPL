package lexer

import (
	"fmt"
	"unicode"
)

// TokenType defines all possible token kinds in UPL
type TokenType int

const (
	TokenEOF TokenType = iota
	TokenError

	// Literals
	TokenString
	TokenNumber
	TokenIdent

	// ============= TIER 1: CORE KEYWORDS (60 total) =============

	// Declarations (6)
	TokenConst        // स्थिर बनाओ / const
	TokenLet          // बदलाव योग्य / let
	TokenVar          // अपरिवर्तनीय / var
	TokenFunction     // कार्य बनाओ / function
	TokenReturn       // वापस करो / return
	TokenArguments    // तर्क / arguments

	// Control Flow (12)
	TokenIf        // यदि / if
	TokenElse      // वरना / else
	TokenSwitch    // विकल्प / switch
	TokenCase      // स्थिति / case
	TokenDefault   // अन्यथा / default
	TokenFor       // के लिए / for
	TokenWhile     // जब तक / while
	TokenDo        // करो / do
	TokenIn        // में / in
	TokenOf        // का / of
	TokenBreak     // तोड़ो / break
	TokenContinue  // जारी रखो / continue

	// Object-Oriented (9)
	TokenClass       // श्रेणी / class
	TokenNew         // नया / new
	TokenThis        // यह / this
	TokenExtends     // बढ़ाओ / extends
	TokenSuper       // मुख्य / super
	TokenConstructor // बनावट / constructor
	TokenGet         // पाओ / get
	TokenSet         // बदलो / set
	TokenStatic      // स्थिरता / static

	// Asynchronous (3)
	TokenAsync   // अतुल्यकालिक / async
	TokenAwait   // इंतजार करो / await
	TokenPromise // वचन / Promise

	// Type & Reflection (6)
	TokenTypeof     // प्रकार / typeof
	TokenInstanceof // उदाहरण / instanceof
	TokenVoid       // शून्य / void
	TokenDelete     // हटाओ / delete
	TokenDebugger   // त्रुटि_खोजो / debugger
	TokenEval       // मूल्यांकन / eval

	// Exception Handling (4)
	TokenTry     // प्रयास करो / try
	TokenCatch   // पकड़ो / catch
	TokenThrow   // फेंको / throw
	TokenFinally // अंततः / finally

	// Value Literals (6)
	TokenTrue      // सत्य / true
	TokenFalse     // असत्य / false
	TokenNull      // रिक्त / null
	TokenUndefined // अपरिभाषित / undefined
	TokenNaN       // संख्या_नहीं / NaN
	TokenInfinity  // अनंत / Infinity

	// Core Global Objects (8)
	TokenLog      // दिखाओ / console.log
	TokenWindow   // खिड़की / window
	TokenDocument // दस्तावेज़ / document
	TokenMath     // गणित / Math
	TokenJSON     // जेसन / JSON
	TokenDate     // दिनांक / Date
	TokenArray    // व्यूह / Array
	TokenObject   // वस्तु / Object

	// UPL-Specific (3)
	TokenStyle    // शैली / style
	TokenColorHub // रंग-मंच / colorhub
	TokenJoin     // जोड़ो / import

	// ============= OPERATORS & DELIMITERS =============
	TokenColon        // :
	TokenSemicolon    // ;
	TokenComma        // ,
	TokenDot          // .
	TokenEquals       // =
	TokenDoubleEquals // ==
	TokenTripleEquals // ===
	TokenNotEquals    // !=
	TokenPlus         // +
	TokenMinus        // -
	TokenStar         // *
	TokenSlash        // /
	TokenPercent      // %
	TokenAnd          // &&
	TokenOr           // ||
	TokenNot          // !
	TokenLeftParen    // (
	TokenRightParen   // )
	TokenLeftBrace    // {
	TokenRightBrace   // }
	TokenLeftBracket  // [
	TokenRightBracket // ]
	TokenArrow        // =>
	TokenQuestionMark // ?

	// Indentation
	TokenIndent
	TokenDedent
	TokenNewline
)

var tokenNames = map[TokenType]string{
	TokenEOF:   "EOF",
	TokenError: "ERROR",

	TokenString:  "STRING",
	TokenNumber:  "NUMBER",
	TokenIdent:   "IDENT",
	TokenColon:   ":",
	TokenNewline: "NEWLINE",
	TokenIndent:  "INDENT",
	TokenDedent:  "DEDENT",
}

// Token represents a single token in the source code
type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func (t Token) String() string {
	if name, ok := tokenNames[t.Type]; ok {
		return fmt.Sprintf("%s(%q) @ %d:%d", name, t.Value, t.Line, t.Column)
	}
	return fmt.Sprintf("Token(%d) @ %d:%d", t.Type, t.Line, t.Column)
}

// KeywordMap maps Hindi and English keywords to token types
// TIER 1: 60 CORE KEYWORDS (all explicit, no shortcuts)
var KeywordMap = map[string]TokenType{
	// ============= DECLARATIONS (6 keywords × 2 languages = 12 mappings) =============
	"स्थिर बनाओ":   TokenConst,
	"const":        TokenConst,
	"बदलाव योग्य": TokenLet,
	"let":          TokenLet,
	"अपरिवर्तनीय": TokenVar,
	"var":          TokenVar,
	"कार्य बनाओ":  TokenFunction,
	"function":     TokenFunction,
	"वापस करो":    TokenReturn,
	"return":       TokenReturn,
	"तर्क":        TokenArguments,
	"arguments":    TokenArguments,

	// ============= CONTROL FLOW (12 keywords × 2 languages = 24 mappings) =============
	"यदि":      TokenIf,
	"if":       TokenIf,
	"वरना":     TokenElse,
	"else":     TokenElse,
	"विकल्प":   TokenSwitch,
	"switch":   TokenSwitch,
	"स्थिति":  TokenCase,
	"case":     TokenCase,
	"अन्यथा":   TokenDefault,
	"default":  TokenDefault,
	"के लिए":   TokenFor,
	"for":      TokenFor,
	"जब तक":    TokenWhile,
	"while":    TokenWhile,
	"करो":     TokenDo,
	"do":       TokenDo,
	"में":      TokenIn,
	"in":       TokenIn,
	"का":      TokenOf,
	"of":       TokenOf,
	"तोड़ो":    TokenBreak,
	"break":    TokenBreak,
	"जारी रखो": TokenContinue,
	"continue": TokenContinue,

	// ============= OBJECT-ORIENTED (9 keywords × 2 languages = 18 mappings) =============
	"श्रेणी":     TokenClass,
	"class":      TokenClass,
	"नया":       TokenNew,
	"new":        TokenNew,
	"यह":        TokenThis,
	"this":       TokenThis,
	"बढ़ाओ":     TokenExtends,
	"extends":    TokenExtends,
	"मुख्य":     TokenSuper,
	"super":      TokenSuper,
	"बनावट":    TokenConstructor,
	"constructor": TokenConstructor,
	"पाओ":       TokenGet,
	"get":        TokenGet,
	"बदलो":     TokenSet,
	"set":        TokenSet,
	"स्थिरता":  TokenStatic,
	"static":     TokenStatic,

	// ============= ASYNCHRONOUS (3 keywords × 2 languages = 6 mappings) =============
	"अतुल्यकालिक": TokenAsync,
	"async":       TokenAsync,
	"इंतजार करो": TokenAwait,
	"await":       TokenAwait,
	"वचन":       TokenPromise,
	"Promise":     TokenPromise,

	// ============= TYPE & REFLECTION (6 keywords × 2 languages = 12 mappings) =============
	"प्रकार":       TokenTypeof,
	"typeof":      TokenTypeof,
	"उदाहरण":     TokenInstanceof,
	"instanceof":  TokenInstanceof,
	"शून्य":      TokenVoid,
	"void":        TokenVoid,
	"हटाओ":      TokenDelete,
	"delete":      TokenDelete,
	"त्रुटि_खोजो": TokenDebugger,
	"debugger":    TokenDebugger,
	"मूल्यांकन":   TokenEval,
	"eval":        TokenEval,

	// ============= EXCEPTION HANDLING (4 keywords × 2 languages = 8 mappings) =============
	"प्रयास करो": TokenTry,
	"try":        TokenTry,
	"पकड़ो":     TokenCatch,
	"catch":      TokenCatch,
	"फेंको":     TokenThrow,
	"throw":      TokenThrow,
	"अंततः":    TokenFinally,
	"finally":    TokenFinally,

	// ============= VALUE LITERALS (6 keywords × 2 languages = 12 mappings) =============
	"सत्य":         TokenTrue,
	"true":         TokenTrue,
	"असत्य":       TokenFalse,
	"false":        TokenFalse,
	"रिक्त":       TokenNull,
	"null":         TokenNull,
	"अपरिभाषित":  TokenUndefined,
	"undefined":    TokenUndefined,
	"संख्या_नहीं": TokenNaN,
	"NaN":          TokenNaN,
	"अनंत":       TokenInfinity,
	"Infinity":     TokenInfinity,

	// ============= CORE GLOBAL OBJECTS (8 keywords × 2 languages = 16 mappings) =============
	"दिखाओ":      TokenLog,
	"console.log": TokenLog,
	"खिड़की":     TokenWindow,
	"window":      TokenWindow,
	"दस्तावेज़":  TokenDocument,
	"document":    TokenDocument,
	"गणित":      TokenMath,
	"Math":        TokenMath,
	"जेसन":      TokenJSON,
	"JSON":        TokenJSON,
	"दिनांक":    TokenDate,
	"Date":        TokenDate,
	"व्यूह":     TokenArray,
	"Array":       TokenArray,
	"वस्तु":     TokenObject,
	"Object":      TokenObject,

	// ============= UPL-SPECIFIC (3 keywords × 2 languages = 6 mappings) =============
	"शैली":      TokenStyle,
	"style":       TokenStyle,
	"रंग-मंच":   TokenColorHub,
	"colorhub":    TokenColorHub,
	"जोड़ो":     TokenJoin,
	"import":      TokenJoin,

	// TOTAL: 60 keywords × 2 = 120 explicit mappings
}

// Lexer tokenizes UPL source code
type Lexer struct {
	input       []rune
	position    int
	readPosition int
	current     rune
	line        int
	column      int
	lineStart   int
	indentStack []int
}

// New creates a new lexer from input string
func New(input string) *Lexer {
	l := &Lexer{
		input:       []rune(input),
		position:    0,
		readPosition: 0,
		line:        1,
		column:      0,
		lineStart:   0,
		indentStack: []int{0},
	}
	l.readChar()
	return l
}

// readChar advances to the next rune
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.current = 0 // EOF
	} else {
		l.current = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	if l.current == '\n' {
		l.line++
		l.lineStart = l.position + 1
		l.column = 0
	} else {
		l.column = l.position - l.lineStart
	}
}

// peekChar returns the next rune without advancing
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// peekCharN returns the rune n positions ahead
func (l *Lexer) peekCharN(n int) rune {
	pos := l.readPosition + n - 1
	if pos >= len(l.input) {
		return 0
	}
	return l.input[pos]
}

// skipWhitespace skips spaces and tabs but not newlines
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' {
		l.readChar()
	}
}

// NextToken returns the next token in the sequence
func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	tok := Token{Line: l.line, Column: l.column}

	switch l.current {
	case '\n':
		tok.Type = TokenNewline
		tok.Value = "\\n"
		l.readChar()
	case ':':
		tok.Type = TokenColon
		tok.Value = ":"
		l.readChar()
	case ';':
		tok.Type = TokenSemicolon
		tok.Value = ";"
		l.readChar()
	case ',':
		tok.Type = TokenComma
		tok.Value = ","
		l.readChar()
	case '.':
		tok.Type = TokenDot
		tok.Value = "."
		l.readChar()
	case '(':
		tok.Type = TokenLeftParen
		tok.Value = "("
		l.readChar()
	case ')':
		tok.Type = TokenRightParen
		tok.Value = ")"
		l.readChar()
	case '{':
		tok.Type = TokenLeftBrace
		tok.Value = "{"
		l.readChar()
	case '}':
		tok.Type = TokenRightBrace
		tok.Value = "}"
		l.readChar()
	case '[':
		tok.Type = TokenLeftBracket
		tok.Value = "["
		l.readChar()
	case ']':
		tok.Type = TokenRightBracket
		tok.Value = "]"
		l.readChar()
	case '+':
		tok.Type = TokenPlus
		tok.Value = "+"
		l.readChar()
	case '-':
		if l.peekChar() == '>' {
			tok.Type = TokenArrow
			tok.Value = "=>"
			l.readChar()
			l.readChar()
		} else {
			tok.Type = TokenMinus
			tok.Value = "-"
			l.readChar()
		}
	case '*':
		tok.Type = TokenStar
		tok.Value = "*"
		l.readChar()
	case '/':
		tok.Type = TokenSlash
		tok.Value = "/"
		l.readChar()
	case '%':
		tok.Type = TokenPercent
		tok.Value = "%"
		l.readChar()
	case '=':
		if l.peekChar() == '=' {
			if l.peekCharN(2) == '=' {
				tok.Type = TokenTripleEquals
				tok.Value = "==="
				l.readChar()
				l.readChar()
				l.readChar()
			} else {
				tok.Type = TokenDoubleEquals
				tok.Value = "=="
				l.readChar()
				l.readChar()
			}
		} else if l.peekChar() == '>' {
			tok.Type = TokenArrow
			tok.Value = "=>"
			l.readChar()
			l.readChar()
		} else {
			tok.Type = TokenEquals
			tok.Value = "="
			l.readChar()
		}
	case '!':
		if l.peekChar() == '=' {
			tok.Type = TokenNotEquals
			tok.Value = "!="
			l.readChar()
			l.readChar()
		} else {
			tok.Type = TokenNot
			tok.Value = "!"
			l.readChar()
		}
	case '&':
		if l.peekChar() == '&' {
			tok.Type = TokenAnd
			tok.Value = "&&"
			l.readChar()
			l.readChar()
		} else {
			tok.Type = TokenError
			tok.Value = "&"
			l.readChar()
		}
	case '|':
		if l.peekChar() == '|' {
			tok.Type = TokenOr
			tok.Value = "||"
			l.readChar()
			l.readChar()
		} else {
			tok.Type = TokenError
			tok.Value = "|"
			l.readChar()
		}
	case '?':
		tok.Type = TokenQuestionMark
		tok.Value = "?"
		l.readChar()
	case '"', '\'', '`':
		quote := l.current
		l.readChar()
		start := l.position
		for l.current != 0 && l.current != quote {
			if l.current == '\\' {
				l.readChar()
				if l.current != 0 {
					l.readChar()
				}
			} else {
				l.readChar()
			}
		}
		tok.Type = TokenString
		tok.Value = string(l.input[start : l.position])
		if l.current == quote {
			l.readChar()
		}
	case 0:
		tok.Type = TokenEOF
		tok.Value = ""
	default:
		if isLetter(l.current) {
			return l.readIdentifier(tok)
		} else if isDigit(l.current) {
			return l.readNumber(tok)
		} else {
			tok.Type = TokenError
			tok.Value = string(l.current)
			l.readChar()
		}
	}

	return tok
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier(tok Token) Token {
	start := l.position
	for isLetter(l.current) || isDigit(l.current) || l.current == '_' || l.current == '-' {
		l.readChar()
	}
	value := string(l.input[start : l.position])

	if t, ok := KeywordMap[value]; ok {
		tok.Type = t
	} else {
		tok.Type = TokenIdent
	}
	tok.Value = value
	return tok
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber(tok Token) Token {
	start := l.position
	for isDigit(l.current) || l.current == '.' {
		l.readChar()
	}
	tok.Type = TokenNumber
	tok.Value = string(l.input[start : l.position])
	return tok
}

// isLetter checks if a rune is a letter (including Hindi)
func isLetter(ch rune) bool {
	return unicode.IsLetter(ch)
}

// isDigit checks if a rune is a digit
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}

// TokenizeAll tokenizes the entire input at once
func (l *Lexer) TokenizeAll() []Token {
	var tokens []Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == TokenEOF {
			break
		}
	}
	return tokens
}