package internal

import (
	"bufio"
	"io"
	"strings"
)

// TODO:MJR A simc document is simply a slice of statements
// Where each statement is represented as a key-value pair with optional attributes and comments
// The key is the value before the first =
// The value is the value after the first =
// The value may be empty, i.e. "foo="

// Parsing Rules:
// A leading comment line followed by a statement line attaches that comment to the statement
// The simulationcraft addon output formats as:
// 	# Saved Loadout: M+
// 	# talents=CEkAG5bbocFKcv+yIq8fPd6ORZmZmZmZWmxMzMGzwYmBAAAAAAMLGz2MMzAzYb2mZmxYglB2mNzYYWYMmZGzYDAAAYAAAAMzgBAAAgB
// and
// 	# Spellsnap Shadowmask (289)
//	head=,id=251109,enchant_id=8017,bonus_id=13440/6652/12667/13577/12699/12806
// Similarly for deriving item/talent candidates, those are exported as comments:
// 	# Saved Loadout: Raid
// 	# talents=CEkAG5bbocFKcv+yIq8fPd6ORZGMzMzmxMzMmZGGzMAAAAAAgZ5BGz2MMzMbzMjtZbegZYMMWGYbWMjhZjxYmZMsBAAAAAAAwMDGAAAAG
// and
// 	# Devouring Reaver's Intake (276)
//	# head=,id=250033,enchant_id=8017,bonus_id=6652/12667/13440/13338/13575/12798

type operator string

type Statement struct {
	Key      string
	Operator operator
	Value    string
	RawSimc  string
	Comment  string
	Disabled bool
	Line     int
}

type StatementUnmarshaler interface {
	UnmarshalStatement(Statement) error
}

const (
	operatorAssign operator = "="
	operatorAppend operator = "+="
)

// ParseStatements extracts statements from a simc document
func ParseStatements(r io.Reader) ([]Statement, error) {
	res := []Statement{}

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	var detachedComment string
	for scanner.Scan() {
		if scanner.Err() != nil {
			return []Statement{}, scanner.Err()
		}

		line := strings.TrimSpace(scanner.Text())
		lineNumber++
		if line == "" {
			detachedComment = ""
			continue
		}
		if strings.HasPrefix(line, "#") {
			// Skip empty comments
			strippedComment := strings.TrimSpace(strings.TrimLeft(line, "#"))
			if strippedComment == "" {
				detachedComment = ""
				continue
			}

			// Commented line with value, but not a statement
			if !strings.Contains(line, "=") {
				detachedComment = strippedComment
				continue
			}
		}
		stmt := newStatement(line, lineNumber)
		if detachedComment != "" {
			stmt.Comment = detachedComment
		}
		detachedComment = ""
		res = append(res, stmt)
	}
	return res, nil
}

func newStatement(line string, lineNum int) (res Statement) {
	res.Line = lineNum
	res.RawSimc = line

	if strings.HasPrefix(line, "#") {
		res.Disabled = true
		line = strings.TrimLeft(line, "#")
	}
	line = strings.TrimSpace(line)

	if strings.Contains(line, string(operatorAppend)) {
		res.Operator = operatorAppend
	} else if strings.Contains(line, string(operatorAssign)) {
		res.Operator = operatorAssign
	}

	if res.Operator != "" {
		key, value, found := strings.Cut(line, string(res.Operator))
		if !found {
			res.Value = line
		} else {
			res.Key = strings.TrimSpace(key)
			res.Value = strings.TrimSpace(value)
		}
	} else {
		res.Value = line
	}
	return res
}
