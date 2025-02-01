package output

type Separator string

const (
	Comma   Separator = ","
	Newline Separator = "\n"
)

func GetSeparator(s string) Separator {
	switch s {
	case "comma":
		return Comma
	case "newline":
		return Newline
	default:
		return ""
	}
}

type Delimiter string

const (
	Colon       Delimiter = ":"
	SingleQuote Delimiter = "'"
	DoubleQuote Delimiter = "\""
)

func GetDelimiter(d string) Delimiter {
	switch d {
	case "colon":
		return Colon
	case "singlequote":
		return SingleQuote
	case "doublequote":
		return DoubleQuote
	default:
		return ""
	}
}
