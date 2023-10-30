package field

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFieldIndex(t *testing.T) {
	t.Run("returns index from field name", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, 1, indexTag.Id)
	})

	t.Run("returns index from field tag instead of field name when both match", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string `index:"2"`
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, 2, indexTag.Id)
	})

	t.Run("returns index from field tag", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name   string `index:"abcd"`
			F      string `index:"02"`
			Amount string `index:"3"`
		}{}).Elem()

		// get index from field Name
		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, -1, indexTag.Id)

		// get index from field F
		indexTag = NewIndexTag(st.Type().Field(1))
		require.Equal(t, 2, indexTag.Id)

		// get index from field Amount
		indexTag = NewIndexTag(st.Type().Field(2))
		require.Equal(t, 3, indexTag.Id)
	})

	t.Run("returns empty string when no tag and field name does not match the pattern", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name string
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, -1, indexTag.Id)
		require.Empty(t, indexTag.Tag)

		// single letter field without tag is ignored
		st = reflect.ValueOf(&struct {
			F string
		}{}).Elem()

		indexTag = NewIndexTag(st.Type().Field(0))
		require.Equal(t, -1, indexTag.Id)
		require.Empty(t, indexTag.Tag)
	})
}

func TestFieldIndexTag(t *testing.T) {
	t.Run("returns index from field name", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, "1", indexTag.Tag)
		require.Equal(t, 1, indexTag.Id)
	})

	t.Run("returns index from field tag instead of field name when both match", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			F1 string `index:"AB"`
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, "AB", indexTag.Tag)
	})

	t.Run("returns index from field tag", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name string `index:"abcd"`
			F    string `index:"02"`
		}{}).Elem()

		// get index from field Name
		indexTag := NewIndexTag(st.Type().Field(0))
		require.Equal(t, "abcd", indexTag.Tag)

		// get index from field F
		indexTag = NewIndexTag(st.Type().Field(1))
		require.Equal(t, "02", indexTag.Tag)
		require.Equal(t, 2, indexTag.Id)
	})

	t.Run("returns empty string when no tag and field name does not match the pattern", func(t *testing.T) {
		st := reflect.ValueOf(&struct {
			Name string
		}{}).Elem()

		indexTag := NewIndexTag(st.Type().Field(0))
		require.Empty(t, indexTag.Tag)

		// single letter field without tag is ignored
		st = reflect.ValueOf(&struct {
			F string
		}{}).Elem()

		indexTag = NewIndexTag(st.Type().Field(0))
		require.Empty(t, indexTag.Tag)
	})
}
