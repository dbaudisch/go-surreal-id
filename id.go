package thing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid/v5"
)

/*
SurrealDB supports these types for ID

- Random Generated ID

- Text ID
	Example:
	-	person:tobie
	-	article:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩
	-	article:⟨10⟩ (CREATE article SET id = "10";)

- Numeric ID
	Consists of a 64-bit int as ID

- Object ID
	Example: temperature:{ location: 'London', date: $now }

- Array ID
	Example: temperature:['London', $now]
*/

type Identifiable interface {
	string | uuid.UUID | int64 | ArrayId | ObjectId
}

type (
	ArrayId  []any
	ObjectId map[string]any
)

type Id struct {
	Val any
}

var ISO8601RegEx = regexp.MustCompile(`^\d{4}-(?:0[1-9]|1[0-2])-(?:[0-2][1-9]|[1-3]0|3[01])T(?:[0-1][0-9]|2[0-3])(?::[0-6]\d)(?::[0-6]\d)?(?:\.\d{3})?(?:[+-][0-2]\d:[0-5]\d|Z)?`)

func ParseId(id string) *Id {
	switch true {
	case isNumeric(id):
		n, _ := strconv.ParseInt(id, 10, 64)
		return &Id{n}
	case strings.ContainsAny(id, "{}"):
		return &Id{parseObjectId(id)}
	case strings.ContainsAny(id, "[]"):
		return &Id{parseArrayId(id)}
	case strings.ContainsAny(id, "⟨⟩"):
		trimmed := strings.Trim(id, "⟨⟩")

		if isNumeric(trimmed) {
			return &Id{trimmed}
		}

		_uuid, _ := uuid.FromString(trimmed)
		return &Id{_uuid}
	default:
		return &Id{id}
	}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func parseArrayId(s string) ArrayId {
	p := strings.Split(strings.Trim(s, "[ ]"), ",")

	res := make(ArrayId, len(p))
	for i, x := range p {
		var tmp any
		tmp = strings.Trim(strings.TrimSpace(x), "'")

		if ISO8601RegEx.MatchString(fmt.Sprint(tmp)) {
			tmp, _ = time.Parse(time.RFC3339, fmt.Sprint(tmp))
		}

		res[i] = tmp
	}

	return res
}

func parseObjectId(s string) ObjectId {
	res := make(ObjectId)
	props := strings.Split(strings.Trim(s, "{ }"), ",")

	for _, p := range props {
		prop := strings.SplitN(p, ":", 2)
		tmp := make([]any, len(prop))

		for i, x := range prop {
			tmp[i] = strings.TrimSpace(x)
		}

		tmp[1] = strings.Trim(fmt.Sprint(tmp[1]), "'")

		if ISO8601RegEx.MatchString(fmt.Sprint(tmp[1])) {
			tmp[1], _ = time.Parse(time.RFC3339, fmt.Sprint(tmp[1]))
		}

		res[fmt.Sprint(tmp[0])] = tmp[1]
	}

	return res
}

func (ai ArrayId) String() string {
	res := make([]string, len(ai))
	for i, v := range ai {
		switch x := v.(type) {
		case time.Time:
			res[i] = x.Format(time.RFC3339Nano)
		default:
			res[i] = x.(string)
		}
	}
	return fmt.Sprintf("['%s']", strings.Join(res, "', '"))
}

func (oi ObjectId) String() string {
	var res []string
	for k, v := range oi {
		switch x := v.(type) {
		case time.Time:
			res = append(res, fmt.Sprintf("%s: '%v'", k, x.Format(time.RFC3339Nano)))
		default:
			res = append(res, fmt.Sprintf("%s: '%v'", k, v))
		}
	}
	return fmt.Sprintf("{%s}", strings.Join(res, ", "))
}

func (id Id) String() string {
	switch val := id.Val.(type) {
	case uuid.UUID:
		return fmt.Sprintf("⟨%s⟩", val.String())
	case string:
		if isNumeric(val) {
			val = fmt.Sprintf("⟨%s⟩", val)
		}

		return fmt.Sprint(val)
	default:
		return fmt.Sprint(val)
	}
}
