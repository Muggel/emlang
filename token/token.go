package token

type Token int

const (
	ILLEGAL Token = iota
	EOF
	IDENT
	INT

	ADD
	SUB
	MUL
	DIV

	LPAREN
	LBRACE

	RPAREN
	RBRACE

	ASSIGN
	COMMA

	FN
	RETURN
)

var tokens = [...]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT: "IDENT",
	INT:   "INT",

	ADD: "+",
	SUB: "-",
	MUL: "*",
	DIV: "/",

	LPAREN: "(",
	LBRACE: "{",

	RPAREN: ")",
	RBRACE: "}",

	ASSIGN: "=",
	COMMA:  ",",

	FN:     "fn",
	RETURN: "return",
}

func (t Token) String() string {
	return tokens[t]
}

var keywords = map[string]Token{
	"fn":     FN,
	"return": RETURN,
}

func Lookup(ident string) Token {
	if kw, ok := keywords[ident]; ok {
		return kw
	}
	return IDENT
}
