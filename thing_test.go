package thing

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/gofrs/uuid/v5"
)

func TestParse(t *testing.T) {
	now := "2024-02-16T00:18:48.084Z"
	date, _ := time.Parse(time.RFC3339, now)
	_uuid, _ := uuid.FromString("8424486b-85b3-4448-ac8d-5d51083391c7")

	tests := []struct {
		name     string
		input    string
		expected *Thing
	}{
		{"Text ID", "person:tobie", &Thing{"person", &Id{"tobie", false}}},
		{
			"Complex Text ID",
			"entry:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
			&Thing{"entry", &Id{_uuid, true}},
		},
		{"Numeric ID", "entry:1337", &Thing{"entry", &Id{int64(1337), false}}},
		{
			"Array ID",
			fmt.Sprintf("entry:['London', '%s']", now),
			&Thing{"entry", &Id{ArrayId{"London", date}, true}},
		},
		{
			"Object ID",
			fmt.Sprintf("entry:{ location: 'London', date: '%s' }", now),
			&Thing{"entry", &Id{ObjectId{"location": "London", "date": date}, true}},
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
		{"Text ID", &Thing{"person", &Id{"tobie", false}}, "person:tobie"},
		{
			"Complex Text ID",
			&Thing{"entry", &Id{_uuid, true}},
			"entry:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩",
		},
		{"Numeric ID", &Thing{"entry", &Id{int64(1337), false}}, "entry:1337"},
		{
			"Array ID",
			&Thing{"entry", &Id{ArrayId{"London", date}, true}},
			fmt.Sprintf("entry:['London', '%s']", now),
		},
		{
			"Object ID",
			&Thing{"entry", &Id{ObjectId{"location": "London", "date": date}, true}},
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
