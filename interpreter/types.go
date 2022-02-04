package interpreter

type BeccaType string

const (
	TString   BeccaType = "string"
	TInt      BeccaType = "int"
	TBool     BeccaType = "bool"
	TFloat    BeccaType = "float"
	TNull     BeccaType = "null"
	TFunction BeccaType = "function"
	TType     BeccaType = "type"
)

func Type(v *BeccaValue) *BeccaValue {
	return &BeccaValue{
		Type:  TType,
		Value: v.Value,
		// TODO - Support a function invocation
		// option here
	}
}
