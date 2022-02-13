package iso8583

func Marshal(message *Message, v interface{}) error {
	return message.Marshal(v)
}

// Unmarshal traverses the message fields recursively and for each field, sets
// the field value into the corresponding struct field value pointed by v. If v
// is nil or not a pointer it returns error.
func Unmarshal(message *Message, v interface{}) error {
	return message.Unmarshal(v)
}
