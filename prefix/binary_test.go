package prefix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFixedBinary(t *testing.T) {
	pref := binaryFixedPrefixer{}

	dataLen, read, err := pref.DecodeLength(16, []byte("some data"))

	require.NoError(t, err)
	require.Equal(t, 16, dataLen)
	require.Equal(t, 0, read)
}

func Test_binaryVarPrefixer_EncodeLength(t *testing.T) {
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
			name: "success(L)",
			fields: fields{
				Digits: 1,
			},

			args: args{
				maxLen:  32,
				dataLen: 24,
			},
			want:    []byte{0x18},
			wantErr: false,
		},
		{
			name: "success(LL)",
			fields: fields{
				Digits: 2,
			},

			args: args{
				maxLen:  512,
				dataLen: 256,
			},
			want:    []byte{0x01, 0x00},
			wantErr: false,
		},
		{
			name: "success(LLL)",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen:  19999999,
				dataLen: 11258879,
			},
			want:    []byte{0xab, 0xcb, 0xff},
			wantErr: false,
		},
		{
			name: "success(LLL) with prepended zeros",
			fields: fields{
				Digits: 3,
			},

			args: args{
				maxLen:  32,
				dataLen: 24,
			},
			want:    []byte{0x00, 0x00, 0x18},
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
				dataLen: 256,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &binaryVarPrefixer{
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

func Test_binaryVarPrefixer_DecodeLength(t *testing.T) {
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
			name: "success(L)",
			fields: fields{
				Digits: 1,
			},
			args: args{
				maxLen: 32,
				data:   []byte{0x18},
			},
			wantDataLen: 24,
			wantRead:    1,
			wantErr:     false,
		},
		{
			name: "success(L)_withMoreData",
			fields: fields{
				Digits: 1,
			},
			args: args{
				maxLen: 32,
				data:   []byte{0x18, 0x0, 0x0, 0x0},
			},
			wantDataLen: 24,
			wantRead:    1,
			wantErr:     false,
		},
		{
			name: "success(LLL)",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 19999999,
				data:   []byte{0xab, 0xcb, 0xff},
			},
			wantDataLen: 11258879,
			wantRead:    3,
			wantErr:     false,
		},
		{
			name: "success(LLL)_withMoreData",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 19999999,
				data:   []byte{0xab, 0xcb, 0xff, 0x0, 0x0, 0x0},
			},
			wantDataLen: 11258879,
			wantRead:    3,
			wantErr:     false,
		},
		{
			name: "not_enough_data",
			fields: fields{
				Digits: 3,
			},
			args: args{
				maxLen: 32,
				data:   []byte("0"),
			},
			wantDataLen: 0,
			wantRead:    0,
			wantErr:     true,
		},
		{
			name: "data_length_exceeds_max_len",
			fields: fields{
				Digits: 1,
			},
			args: args{
				maxLen: 8,
				data:   []byte{0x18},
			},
			wantDataLen: 0,
			wantRead:    0,
			wantErr:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &binaryVarPrefixer{
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
