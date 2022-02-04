package interpreter

func BoolFromGoBoolean(x bool) *BeccaValue {
	if x {
		return True()
	}
	return False()
}

func False() *BeccaValue {
	return &BeccaValue{
		Type:  TBool,
		Value: false,
	}
}

func True() *BeccaValue {
	return &BeccaValue{
		Type:  TBool,
		Value: true,
	}
}

func Truthiness(value *BeccaValue) *BeccaValue {
	switch value.Type {
	case TBool:
		return value
	case TString:
		if value.Value.(string) == "" {
			return False()
		} else {
			return True()
		}
	case TInt:
		if value.Value.(int64) == 0 {
			return False()
		} else {
			return True()
		}
	case TFloat:
		if value.Value.(float64) == 0.0 {
			return False()
		} else {
			return True()
		}
	case TNull:
		return False()
	case TFunction:
		return True()
	}
	// TODO - handle error correctly
	panic("How did we get here")
}
