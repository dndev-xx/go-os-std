package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/require" //nolint:depguard
)

func TestCopy(t *testing.T) {
	tests := []struct {
		name      string
		from      string
		to        string
		offset    int64
		limit     int64
		expectErr error
	}{
		{
			name:      "Error when copying to unsupported file",
			from:      "testdata/input2.txt",
			to:        "/dev/urandom",
			offset:    0,
			limit:     0,
			expectErr: ErrUnsupportedFile,
		},

		{
			name:      "Error for non-existing source file",
			from:      "filename",
			to:        "out.txt",
			offset:    0,
			limit:     0,
			expectErr: nil,
		},
		{
			name:      "Error for negative offset",
			from:      "testdata/input.txt",
			to:        "out.txt",
			offset:    -10,
			limit:     0,
			expectErr: nil,
		},
		{
			name:      "Error for no permission on source file",
			from:      "fromFile",
			to:        "out.txt",
			offset:    0,
			limit:     0,
			expectErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.name == "Error for no permission on source file" {
				toFile, _ := os.OpenFile(tt.from, os.O_CREATE|os.O_WRONLY, 0o000)
				_ = toFile.Close()
				defer os.Remove(tt.from)
			}

			err := Copy(tt.from, tt.to, tt.offset, tt.limit)

			if tt.expectErr != nil {
				require.Error(t, err)
				require.False(t, errors.Is(err, tt.expectErr))
			} else {
				require.Error(t, err)
			}
		})
	}
}
