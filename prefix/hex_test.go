package prefix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedHex(t *testing.T) {
	pref := hexFixedPrefixer{}

	dataLen, read, err := pref.DecodeLength(16, []byte("whatever"))

	require.NoError(t, err)
	require.Equal(t, 16, dataLen)
	require.Equal(t, 0, read)
}

func Test_hexVarPrefixer_EncodeLength(t *testing.T) {
	type fields struct {
		Digits int
	}
	type args struct {
		maxLen  int
		dataLen int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "success",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen:  32,
				dataLen: 24,
			},
			want:    []byte("000018"),
			wantErr: false,
		},
		{
			name: "data_length_exceeds_max_len",
			fields: fields{
				Digits: 1,
			},
			args: args{
				maxLen:  16,
				dataLen: 24,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "data_length_exceeds_max_possible_len",
			fields: fields{
				Digits: 1,
			},
			args: args{
				maxLen:  512,
				dataLen: 512,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &hexVarPrefixer{
				Digits: tt.fields.Digits,
			}
			got, err := p.EncodeLength(tt.args.maxLen, tt.args.dataLen)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got, "EncodeLength() mismatch")
		})
	}
}

func Test_hexVarPrefixer_DecodeLength(t *testing.T) {
	type fields struct {
		Digits int
	}
	type args struct {
		maxLen int
		data   []byte
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDataLen int
		wantRead    int
		wantErr     bool
	}{
		{
			name: "success",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 32,
				data:   []byte("000018whateverwhateverwhatever"),
			},
			wantDataLen: 0x18,
			wantRead:    6,
			wantErr:     false,
		},
		{
			name: "not_enough_data",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 32,
				data:   []byte("0000"),
			},
			wantDataLen: 0,
			wantRead:    0,
			wantErr:     true,
		},
		{
			name: "parse_error",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 32,
				data:   []byte("SSSSSSwhateverwhateverwhatever"),
			},
			wantDataLen: 0,
			wantRead:    0,
			wantErr:     true,
		},
		{
			name: "data_length_exceeds_max_len",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 8,
				data:   []byte("000018whateverwhateverwhatever"),
			},
			wantDataLen: 0,
			wantRead:    0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &hexVarPrefixer{
				Digits: tt.fields.Digits,
			}
			dataLen, read, err := p.DecodeLength(tt.args.maxLen, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("DecodeLength() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.wantDataLen, dataLen, "DataLen mismatch")
			assert.Equal(t, tt.wantRead, read, "Read mismatch")
		})
	}
}
