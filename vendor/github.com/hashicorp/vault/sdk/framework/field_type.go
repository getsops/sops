package framework

// FieldType is the enum of types that a field can be.
type FieldType uint

const (
	TypeInvalid FieldType = 0
	TypeString  FieldType = iota
	TypeInt
	TypeBool
	TypeMap

	// TypeDurationSecond represent as seconds, this can be either an
	// integer or go duration format string (e.g. 24h)
	TypeDurationSecond

	// TypeSignedDurationSecond represents a positive or negative duration
	// as seconds, this can be either an integer or go duration format
	// string (e.g. 24h).
	TypeSignedDurationSecond

	// TypeSlice represents a slice of any type
	TypeSlice

	// TypeStringSlice is a helper for TypeSlice that returns a sanitized
	// slice of strings
	TypeStringSlice

	// TypeCommaStringSlice is a helper for TypeSlice that returns a sanitized
	// slice of strings and also supports parsing a comma-separated list in
	// a string field
	TypeCommaStringSlice

	// TypeLowerCaseString is a helper for TypeString that returns a lowercase
	// version of the provided string
	TypeLowerCaseString

	// TypeNameString represents a name that is URI safe and follows specific
	// rules.  These rules include start and end with an alphanumeric
	// character and characters in the middle can be alphanumeric or . or -.
	TypeNameString

	// TypeKVPairs allows you to represent the data as a map or a list of
	// equal sign delimited key pairs
	TypeKVPairs

	// TypeCommaIntSlice is a helper for TypeSlice that returns a sanitized
	// slice of Ints
	TypeCommaIntSlice

	// TypeHeader is a helper for sending request headers through to Vault.
	// For instance, the AWS and AliCloud credential plugins both act as a
	// benevolent MITM for a request, and the headers are sent through and
	// parsed.
	TypeHeader
)

func (t FieldType) String() string {
	switch t {
	case TypeString:
		return "string"
	case TypeLowerCaseString:
		return "lowercase string"
	case TypeNameString:
		return "name string"
	case TypeInt:
		return "int"
	case TypeBool:
		return "bool"
	case TypeMap:
		return "map"
	case TypeKVPairs:
		return "keypair"
	case TypeDurationSecond, TypeSignedDurationSecond:
		return "duration (sec)"
	case TypeSlice, TypeStringSlice, TypeCommaStringSlice, TypeCommaIntSlice:
		return "slice"
	case TypeHeader:
		return "header"
	default:
		return "unknown type"
	}
}
