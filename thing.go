package thing

import (
	"fmt"
	"strings"
)

type Thing struct {
	Table string
	Id    *Id
}

func New(tb, id string) *Thing {
	return &Thing{tb, ParseId(id)}
}

func Parse(s string) *Thing {
	p := strings.SplitN(s, ":", 2)
	return &Thing{p[0], ParseId(p[1])}
}

func (t *Thing) String() string {
	return fmt.Sprintf("%s:%v", t.Table, t.Id)
}

func (t *Thing) UnmarshalJSON(data []byte) error {
	rid := Parse(strings.Trim(string(data), "\""))
	t.Table = rid.Table
	t.Id = rid.Id

	return nil
}
