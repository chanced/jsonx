package jsonx

type Type uint8

const (
	TypeInvalid Type = iota
	TypeEmpty
	TypeNull
	TypeBool
	TypeNumber
	TypeString
	TypeArray
	TypeObject
)

func (t Type) String() string {
	switch t {
	case TypeNull:
		return "null"
	case TypeBool:
		return "bool"
	case TypeNumber:
		return "number"
	case TypeString:
		return "string"
	case TypeArray:
		return "array"
	case TypeObject:
		return "object"
	case TypeEmpty:
		return "empty"
	default:
		return "invalid"
	}
}

func TypeOf(d []byte) Type {
	if len(d) == 0 {
		return TypeEmpty
	}
	switch {
	case IsArray(d):
		return TypeArray
	case IsObject(d):
		return TypeObject
	case IsString(d):
		return TypeString
	case IsNumber(d):
		return TypeNumber
	case IsBool(d):
		return TypeBool
	case IsNull(d):
		return TypeNull
	default:
		return TypeInvalid
	}
}
