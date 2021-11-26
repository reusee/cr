package main

import "unicode"

type Parser func(
	input rune,
) (
	next Parser,
	err error,
)

func ParseTokens(
	res *[][]rune,
	cont Parser,
) Parser {
	return func(r rune) (Parser, error) {
		if unicode.IsSpace(r) {
			return ParseTokens(res, cont), nil
		}
		if r == '`' || r == '"' || r == '\'' {
			// string
			*res = append(*res, []rune{})
			return ParseQuoted(r, res, ParseTokens(res, cont)), nil
		}
		*res = append(*res, []rune{r})
		return ParseToken(res, ParseTokens(res, cont)), nil
	}
}

func ParseToken(
	res *[][]rune,
	cont Parser,
) Parser {
	return func(r rune) (Parser, error) {
		if unicode.IsSpace(r) {
			return ParseTokens(res, cont), nil
		}
		(*res)[len(*res)-1] = append((*res)[len(*res)-1], r)
		return ParseToken(res, cont), nil
	}
}

func ParseQuoted(
	quote rune,
	res *[][]rune,
	cont Parser,
) Parser {
	return func(r rune) (Parser, error) {
		if r == quote {
			return cont, nil
		} else if r == '\\' {
			// escape
			return ParseEscape(res, cont), nil
		}
		(*res)[len(*res)-1] = append((*res)[len(*res)-1], r)
		return ParseQuoted(quote, res, cont), nil
	}
}

func ParseEscape(
	res *[][]rune,
	cont Parser,
) Parser {
	return func(r rune) (Parser, error) {
		(*res)[len(*res)-1] = append((*res)[len(*res)-1], r)
		return cont, nil
	}
}
