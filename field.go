package iso8583

type Field struct {
	ID    int
	Value string
}

func NewField(id int, value string) *Field {
	return &Field{
		ID:    id,
		Value: value,
	}
}
