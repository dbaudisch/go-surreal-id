package thing

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
)

func TestParseId(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    string
		expected *Id
	}{
		{"Text ID", "tobie", &Id{"tobie"}},
		{
			"Complex Text ID",
			"⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
			&Id{_uuid},
		},
		{"Numeric ID", "1337", &Id{int64(1337)}},
		{
			"Array ID",
			fmt.Sprintf("[ 'London', '%s' ]", now),
			&Id{ArrayId{"London", date}},
		},
		{
			"Object ID",
			fmt.Sprintf("{ location: 'London', date: '%s' }", now),
			&Id{ObjectId{"location": "London", "date": date}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseId(tt.input); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("ParseId(\"%s\") = %#v; expected %#v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestStringId(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    *Id
		expected string
	}{
		{"Text ID", &Id{"tobie"}, "tobie"},
		{
			"Complex Text ID",
			&Id{_uuid},
			"⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
		},
		{"Numeric ID", &Id{int64(1337)}, "1337"},
		{
			"Array ID",
			&Id{ArrayId{"London", date}},
			fmt.Sprintf("['London', '%s']", now),
		},
		{
			"Object ID",
			&Id{ObjectId{"location": "London", "date": date}},
			fmt.Sprintf("{location: 'London', date: '%s'}", now),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("String(\"%#v\") = %#v; expected %s", tt.input, got, tt.expected)
			}
		})
	}
}
