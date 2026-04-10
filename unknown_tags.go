package iso8583

import "github.com/moov-io/iso8583/field"

// UnknownTags returns unknown TLV tags found in the message after unpacking,
// keyed by their dot-separated paths (e.g., "55.9F36" for unknown tag 9F36
// inside field 55). This requires StoreUnknownTLVTags to be enabled in the
// composite field specs that may contain unknown tags.
func UnknownTags(message *Message) map[string]field.Field {
	if message == nil {
		return make(map[string]field.Field)
	}

	result := make(map[string]field.Field)
	collectUnknownTags(&MessageWrapper{message}, nil, "", result)
	return result
}

// UnknownCompositeTags returns unknown TLV tags found in the composite field after unpacking,
// keyed by their dot separated paths (e.g., "9F36" for unknown tag 9F36 inside the composite).
// This requires StoreUnknownTLVTags to be enabled in the composite field specs that may contain
// unknown tags.
func UnknownCompositeTags(composite *field.Composite) map[string]field.Field {
	if composite == nil {
		return make(map[string]field.Field)
	}

	result := make(map[string]field.Field)
	collectUnknownTags(composite, composite.Spec().Subfields, "", result)
	return result
}

// collectUnknownTags recursively walks a field container, comparing
// subfield keys against the spec's subfield definitions. Any key present
// in GetSubfields() but missing from specSubfields is an unknown tag.
// When specSubfields is nil (at the message level), all fields are
// considered known and the function only recurses into sub-containers.
func collectUnknownTags(container FieldContainer, specSubfields map[string]field.Field, prefix string, result map[string]field.Field) {
	for tag, f := range container.GetSubfields() {
		path := tag
		if prefix != "" {
			path = prefix + "." + tag
		}

		if specSubfields != nil {
			if _, known := specSubfields[tag]; !known {
				result[path] = f
				continue
			}
		}

		// recurse into known sub-containers to find unknown tags deeper
		if subContainer, ok := f.(FieldContainer); ok {
			collectUnknownTags(subContainer, f.Spec().Subfields, path, result)
		}
	}
}
