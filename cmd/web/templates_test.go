package main

import (
	"testing"
	"time"

	"snippetbox.alberttseng.net/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		tm   time.Time
		want string
	}{
		{
			name: "UTC",
			tm:   time.Date(2024, 10, 27, 12, 35, 05, 0, time.UTC),
			want: "2024-10-27 at 12:35:05",
		},
		{
			name: "Empty",
			tm:   time.Time{},
			want: "",
		},
		{
			name: "CET",
			tm:   time.Date(2024, 10, 27, 12, 35, 05, 0, time.FixedZone("CET", 1*60*60)),
			want: "2024-10-27 at 11:35:05",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hd := humanDate(tt.tm)

			// if hd != tt.want {
			// 	t.Errorf("got %q; want %q", hd, tt.want)
			// }
			assert.Equal(t, hd, tt.want)
		})
	}
}
