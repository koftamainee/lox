package token

type Type string

const (
	// Single-character tokens.
	LeftParen  Type = "LEFT_PAREN"
	RightParen Type = "RIGHT_PAREN"
	LeftBrace  Type = "LEFT_BRACE"
	RightBrace Type = "RIGHT_BRACE"
	Comma      Type = "COMMA"
	Dot        Type = "DOT"
	Minus      Type = "MINUS"
	Plus       Type = "PLUS"
	Semicolon  Type = "SEMICOLON"
	Slash      Type = "SLASH"
	Star       Type = "STAR"

	// One or two character tokens.
	Bang         Type = "BANG"
	BangEqual    Type = "BANG_EQUAL"
	Equal        Type = "EQUAL"
	EqualEqual   Type = "EQUAL_EQUAL"
	Greater      Type = "GREATER"
	GreaterEqual Type = "GREATER_EQUAL"
	Less         Type = "LESS"
	LessEqual    Type = "LESS_EQUAL"

	// Literals.
	Identifier Type = "IDENTIFIER"
	String     Type = "STRING"
	Number     Type = "NUMBER"

	// Keywords.
	And    Type = "AND"
	Class  Type = "CLASS"
	Else   Type = "ELSE"
	False  Type = "FALSE"
	Fun    Type = "FUN"
	For    Type = "FOR"
	If     Type = "IF"
	Nil    Type = "NIL"
	Or     Type = "OR"
	Print  Type = "PRINT"
	Return Type = "RETURN"
	Super  Type = "SUPER"
	This   Type = "THIS"
	True   Type = "TRUE"
	Var    Type = "VAR"
	While  Type = "WHILE"

	EOF Type = "EOF"
)
