package iso8583

import (
	"fmt"

	"github.com/moov-io/iso8583/field"
)

const (
	emvFirstIndex = 4
	emvLastIndex  = 4
	emvPattern    = " ... "
	pinFirstIndex = 2
	pinLastIndex  = 2
	pinPattern    = "****"
	panFistIndex  = 4
	panLastIndex  = 4
	panPattern    = "****"
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
	if len(in) < emvFirstIndex+emvLastIndex {
		return in
	}

	return in[0:emvFirstIndex] + emvPattern + in[len(in)-emvLastIndex:]
}

var PINFilter = func(in string, data field.Field) string {
	if len(in) < pinFirstIndex+pinLastIndex {
		return in
	}
	return in[0:pinFirstIndex] + pinPattern + in[len(in)-pinLastIndex:]
}

var PANFilter = func(in string, data field.Field) string {
	if len(in) < panFistIndex+panLastIndex {
		return in
	}
	return in[0:panFistIndex] + panPattern + in[len(in)-panLastIndex:]
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
