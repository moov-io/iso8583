package iso8583

import (
	"fmt"

	"github.com/moov-io/iso8583/field"
)

type FilterFunc func(in string, data field.Field) string

type FieldFilter func(fieldFilters map[int]FilterFunc)

func FilterField(id int, filterFn FilterFunc) FieldFilter {
	return func(fieldFilters map[int]FilterFunc) {
		fieldFilters[id] = filterFn
	}
}

var DefaultFilter = func() []FieldFilter {
	filters := []FieldFilter{
		FilterField(2, PANFilter),
		FilterField(20, PANFilter),
		FilterField(35, Track2Filter),
		FilterField(36, Track3Filter),
		FilterField(45, Track1Filter),
		FilterField(52, PINFilter),
		FilterField(55, EMVFilter),
	}
	return filters
}

var EMVFilter = func(in string, data field.Field) string {
	if len(in) < 8 {
		return in
	}

	return in[0:4] + " ... " + in[len(in)-4:]
}

var PINFilter = func(in string, data field.Field) string {
	if len(in) < 4 {
		return in
	}
	return in[0:2] + "****" + in[len(in)-2:]
}

var PANFilter = func(in string, data field.Field) string {
	if len(in) < 8 {
		return in
	}
	return in[0:4] + "****" + in[len(in)-4:]
}

var Track1Filter = func(in string, data field.Field) string {
	track := field.Track1{}
	if err := newTrackData(data, &track); err != nil {
		return in
	}

	track.PrimaryAccountNumber = PANFilter(track.PrimaryAccountNumber, nil)
	return getTrackDataString(in, &track)
}

var Track2Filter = func(in string, data field.Field) string {
	track := field.Track2{}
	if err := newTrackData(data, &track); err != nil {
		return in
	}

	track.PrimaryAccountNumber = PANFilter(track.PrimaryAccountNumber, nil)
	return getTrackDataString(in, &track)
}

var Track3Filter = func(in string, data field.Field) string {
	track := field.Track3{}
	if err := newTrackData(data, &track); err != nil {
		return in
	}
	track.PrimaryAccountNumber = PANFilter(track.PrimaryAccountNumber, nil)
	return getTrackDataString(in, &track)
}

func newTrackData(data, track field.Field) error {
	if raw, err := data.Pack(); err == nil {
		track.SetSpec(data.Spec())
		if _, err := track.Unpack(raw); err != nil {
			return fmt.Errorf("creating new track data")
		}
	}

	return nil
}

func getTrackDataString(in string, track field.Field) string {
	if converted, packErr := track.String(); packErr != nil {
		return in
	} else {
		return converted
	}
}
