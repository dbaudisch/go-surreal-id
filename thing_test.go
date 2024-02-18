package thing

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
)

func TestNew(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    []string
		expected *Thing
	}{
		{
			"Text ID",
			[]string{"person", "tobie"},
			&Thing{"person", &Id{"tobie"}},
		},
		{
			"Complex Text ID",
			[]string{"entry", "⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩"},
			&Thing{"entry", &Id{_uuid}},
		},
		{
			"Complex Numeric ID",
			[]string{"entry", "⟨42⟩"},
			&Thing{"entry", &Id{"42"}},
		},
		{
			"Numeric ID",
			[]string{"entry", "1337"},
			&Thing{"entry", &Id{int64(1337)}},
		},
		{
			"Array ID",
			[]string{"entry", fmt.Sprintf("['London', '%s']", now)},
			&Thing{"entry", &Id{ArrayId{"London", date}}},
		},
		{
			"Object ID",
			[]string{"entry", fmt.Sprintf("{ location: 'London', date: '%s' }", now)},
			&Thing{"entry", &Id{ObjectId{"location": "London", "date": date}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.input[0], tt.input[1]); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestParse(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    string
		expected *Thing
	}{
		{
			"Text ID",
			"person:tobie",
			&Thing{"person", &Id{"tobie"}},
		},
		{
			"Complex Text ID",
			"entry:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
			&Thing{"entry", &Id{_uuid}},
		},
		{
			"Complex Numeric ID",
			"entry:⟨42⟩",
			&Thing{"entry", &Id{"42"}},
		},
		{
			"Numeric ID",
			"entry:1337",
			&Thing{"entry", &Id{int64(1337)}},
		},
		{
			"Array ID",
			fmt.Sprintf("entry:['London', '%s']", now),
			&Thing{"entry", &Id{ArrayId{"London", date}}},
		},
		{
			"Object ID",
			fmt.Sprintf("entry:{ location: 'London', date: '%s' }", now),
			&Thing{"entry", &Id{ObjectId{"location": "London", "date": date}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Parse(tt.input); !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("Parse() = %#v, want %#v", got, tt.expected)
			}
		})
	}
}

func TestString(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    *Thing
		expected string
	}{
		{
			"Text ID",
			&Thing{"person", &Id{"tobie"}},
			"person:tobie",
		},
		{
			"Complex Text ID",
			&Thing{"entry", &Id{_uuid}},
			"entry:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
		},
		{
			"Complex Numeric ID",
			&Thing{"entry", &Id{"42"}},
			"entry:⟨42⟩",
		},
		{
			"Numeric ID",
			&Thing{"entry", &Id{int64(1337)}},
			"entry:1337",
		},
		{
			"Array ID",
			&Thing{"entry", &Id{ArrayId{"London", date}}},
			fmt.Sprintf("entry:['London', '%s']", now),
		},
		{
			"Object ID",
			&Thing{"entry", &Id{ObjectId{"location": "London", "date": date}}},
			fmt.Sprintf("entry:{location: 'London', date: '%s'}", now),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.input.String(); got != tt.expected {
				t.Errorf("String(\"%#v\") = %s; expected %s", tt.input, got, tt.expected)
			}
		})
	}
}
