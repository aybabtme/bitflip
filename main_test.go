package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseByteOffset(t *testing.T) {
	tests := []struct {
		name           string
		in             string
		wantByteOffset uint64
		wantBitOffset  uint8
		wantErr        string
	}{
		{
			name:           "valid string",
			in:             "122@0",
			wantByteOffset: 122,
			wantBitOffset:  0,
			wantErr:        "",
		},
		{
			name:           "valid string",
			in:             "0@0",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "",
		},
		{
			name:           "invalid bit offset",
			in:             "122@8",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "bit offset must be between 0 and 7, was 8",
		},
		{
			name:           "invalid bit offset",
			in:             "122@-1",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid bit offset: strconv.ParseUint: parsing \"-1\": invalid syntax",
		},
		{
			name:           "invalid byte offset",
			in:             "lol@lol",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid byte offset: strconv.ParseUint: parsing \"lol\": invalid syntax",
		},
		{
			name:           "invalid bit offset",
			in:             "1@lol",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid bit offset: strconv.ParseUint: parsing \"lol\": invalid syntax",
		},
		{
			name:           "invalid string",
			in:             "@",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid byte offset: strconv.ParseUint: parsing \"\": invalid syntax",
		},
		{
			name:           "invalid string",
			in:             "@@@",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid offset string, too many '@' symbols found",
		},
		{
			name:           "invalid string",
			in:             "122:0",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        "invalid offset string, no '@' symbols found",
		},
		{
			name:           "invalid string",
			in:             "       122@0        ",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        `invalid byte offset: strconv.ParseUint: parsing "       122": invalid syntax`,
		},
		{
			name:           "invalid byte offset string",
			in:             "-1@0",
			wantByteOffset: 0,
			wantBitOffset:  0,
			wantErr:        `invalid byte offset: strconv.ParseUint: parsing "-1": invalid syntax`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotByteOffset, gotBitOffset, gotErr := parseByteOffset(tt.in)
			require.Equal(t, tt.wantByteOffset, gotByteOffset)
			require.Equal(t, tt.wantBitOffset, gotBitOffset)
			if tt.wantErr != "" {
				require.Contains(t, tt.wantErr, gotErr.Error())
			} else {
				require.NoError(t, gotErr)
			}
		})
	}
}

func TestToggleNthBit(t *testing.T) {
	tests := []struct {
		name     string
		inByte   byte
		inNthBit uint8
		want     byte
	}{
		{
			name:     "0th bit",
			inByte:   0b00000000,
			inNthBit: 0,
			want:     0b00000001,
		},
		{
			name:     "1st bit",
			inByte:   0b00000000,
			inNthBit: 1,
			want:     0b00000010,
		},
		{
			name:     "2nd bit",
			inByte:   0b00000000,
			inNthBit: 2,
			want:     0b00000100,
		},
		{
			name:     "3rd bit",
			inByte:   0b00000000,
			inNthBit: 3,
			want:     0b00001000,
		},
		{
			name:     "4th bit",
			inByte:   0b00000000,
			inNthBit: 4,
			want:     0b00010000,
		},
		{
			name:     "5th bit",
			inByte:   0b00000000,
			inNthBit: 5,
			want:     0b00100000,
		},
		{
			name:     "6th bit",
			inByte:   0b00000000,
			inNthBit: 6,
			want:     0b01000000,
		},
		{
			name:     "7th bit",
			inByte:   0b00000000,
			inNthBit: 7,
			want:     0b10000000,
		},
		{
			name:     "8th bit",
			inByte:   0b00000000,
			inNthBit: 8,
			want:     0b00000000,
		},
		{
			name:     "9th bit",
			inByte:   0b00000000,
			inNthBit: 9,
			want:     0b00000000,
		},
		{
			name:     "32nd bit",
			inByte:   0b00000000,
			inNthBit: 32,
			want:     0b00000000,
		},
		{
			name:     "64th bit",
			inByte:   0b00000000,
			inNthBit: 64,
			want:     0b00000000,
		},
		{
			name:     "65th bit",
			inByte:   0b00000000,
			inNthBit: 65,
			want:     0b00000000,
		},
		{
			name:     "66th bit",
			inByte:   0b00000000,
			inNthBit: 66,
			want:     0b00000000,
		},
		{
			name:     "4th bit (from 0 to 1) in non-zero data",
			inByte:   0b01010101,
			inNthBit: 4,
			want:     0b01000101,
		},
		{
			name:     "4th bit (from 1 to 0) in non-zero data",
			inByte:   0b10101010,
			inNthBit: 4,
			want:     0b10111010,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := toggleNthBit(tt.inByte, tt.inNthBit)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestFlipBitAtOffset(t *testing.T) {
	tests := []struct {
		name         string
		inByteOffset uint64
		inBitOffset  uint8
		in           []byte
		want         []byte
		wantErr      string
	}{
		{
			name:         "regular flip",
			inByteOffset: 0,
			inBitOffset:  0,
			in:           []byte{0b00000001, 0b00000000, 0b00000000, 0b00000000},
			want:         []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000},
			wantErr:      "",
		},
		{
			name:         "mid file flip",
			inByteOffset: 2,
			inBitOffset:  4,
			in:           []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000},
			want:         []byte{0b00000000, 0b00000000, 0b00010000, 0b00000000},
			wantErr:      "",
		},
		{
			name:         "end of file flip",
			inByteOffset: 3,
			inBitOffset:  7,
			in:           []byte{0b00000000, 0b00000000, 0b00000000, 0b00000000},
			want:         []byte{0b00000000, 0b00000000, 0b00000000, 0b10000000},
			wantErr:      "",
		},
	}

	dname, err := ioutil.TempDir(os.TempDir(), "")
	require.NoError(t, err)
	t.Cleanup(func() { os.RemoveAll(dname) })

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// prepare the target file with the data
			f, err := ioutil.TempFile(dname, "")
			require.NoError(t, err)

			fname := f.Name()

			t.Cleanup(func() { os.Remove(fname) })

			_, err = f.Write(tt.in)
			require.NoError(t, err)
			require.NoError(t, f.Close())

			// reopen the target file and do the test
			f, err = os.OpenFile(fname, os.O_RDWR, 0)
			require.NoError(t, err)

			gotErr := flipBitAtOffset(f, int64(tt.inByteOffset), tt.inBitOffset)
			if tt.wantErr != "" {
				require.Contains(t, tt.wantErr, gotErr.Error())
				return
			}
			require.NoError(t, gotErr, "should have flipped bit")

			require.NoError(t, f.Close(), "should be able to close file")

			got, err := ioutil.ReadFile(fname)
			require.NoError(t, err, "should be able to read back the file")
			require.Equal(t, tt.want, got)
		})
	}
}
