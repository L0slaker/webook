package time

import (
	"github.com/go-playground/assert/v2"
	"testing"
)

func TestIsValidDate(t *testing.T) {
	testCases := []struct {
		name       string
		date       string
		layout     string
		wantResult bool
	}{
		{
			name:       "非法的时间格式",
			date:       "00-12-14",
			layout:     YYYYMMDD,
			wantResult: false,
		},
		{
			name:       "合法的时间格式",
			date:       "2000-12-14",
			layout:     YYYYMMDD,
			wantResult: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res := IsValidDate(tc.layout, tc.date)
			assert.Equal(t, tc.wantResult, res)
		})
	}
}
