package thing

import (
	"fmt"
	"strings"
)

type Thing struct {
	table string
	id    *Id
}

func New(tb, id string) *Thing {
	return &Thing{tb, ParseId(id)}
}

func Parse(s string) *Thing {
	p := strings.SplitN(s, ":", 2)
	return &Thing{p[0], ParseId(p[1])}
}

func (t Thing) String() string {
	return fmt.Sprintf("%s:%v", t.table, t.id)
}
