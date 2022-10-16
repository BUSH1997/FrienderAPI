package models

type Bool struct {
	Value   bool
	Defined bool
}

func (b Bool) IsDefined() bool {
	return b.Defined
}

func (b Bool) GetVal() bool {
	return b.Value
}

func (b Bool) IsTrue() bool {
	return b.Value == true
}

func (b Bool) IsFalse() bool {
	return b.Value == false
}

func (b Bool) IsDefinedTrue() bool {
	return b.IsDefined() && b.IsTrue()
}

func (b Bool) IsDefinedFalse() bool {
	return b.IsDefined() && b.IsFalse()
}

func DefinedBool(val bool) Bool {
	return Bool{
		Value:   val,
		Defined: true,
	}
}
