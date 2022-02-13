package iso8583

// Marshal populates message fields with v struct field values. It traverses
// through the message fields and calls Unmarshal(...) on them setting the v If
// v  is nil or not a pointer it returns error.
func Marshal(message *Message, v interface{}) error {
	return message.Marshal(v)
}

// Unmarshal populates v struct fields with message field values. It traverses
// through the message fields and calls Unmarshal(...) on them setting the v If
// v  is nil or not a pointer it returns error.
func Unmarshal(message *Message, v interface{}) error {
	return message.Unmarshal(v)
}
