package simc

// TODO:MJR A simc document is simply a slice of statements
// Where each statement is represented as a key-value pair with optional attributes and comments
// The key is the value before the first =
// The value is the value after the first =
// The value may be empty, i.e. "foo="
// Attrs are the comma-separated key-value pairs after the initial key-value pair
// A comment may trail the statement, indicated by a # character

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

type Document []Statement
type Statement struct {
	Key     string
	Value   string
	Attrs   []Attr
	Comment string
}
type Attr struct {
	Key   string
	Value string
}
