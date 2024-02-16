package thing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
SurrealDB supports these types for ID

- Random Generated ID

- Text ID
	Example:
	-	person:tobie
	-	article:⟨8424486b-85b3-4448-ac8d-5d51083391c7⟩

- Numeric ID
	Consists of a 64-bit int as ID

- Object ID
	Example: temperature:{ location: 'London', date: $now }

- Array ID
	Example: temperature:['London', $now]
*/

type Identifiable interface {
	string | int64 | ArrayId | ObjectId
}

type (
	ArrayId  []any
	ObjectId map[string]any
)

type Id struct {
	val   any
	cmplx bool
}

func ParseId(id string) *Id {
	switch true {
	case isNumeric(id):
		num, _ := strconv.ParseInt(id, 10, 64)
		return &Id{int64(num), false}
	case strings.HasPrefix(id, "{") && strings.HasSuffix(id, "}"):
		return &Id{parseObjectId(id), true}
	case strings.HasPrefix(id, "[") && strings.HasSuffix(id, "]"):
		return &Id{parseArrayId(id), true}
	default:
		switch true {
		case strings.HasPrefix(id, "⟨") && strings.HasSuffix(id, "⟩"):
			return &Id{strings.Trim(id, "⟨⟩"), true}
		default:
			return &Id{id, false}
		}
	}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func parseArrayId(s string) ArrayId {
	p := strings.Split(strings.Trim(s, "[ ]"), ",")
	ISO8601RegEx := regexp.MustCompile(`^\d{4}-(?:0[1-9]|1[0-2])-(?:[0-2][1-9]|[1-3]0|3[01])T(?:[0-1][0-9]|2[0-3])(?::[0-6]\d)(?::[0-6]\d)?(?:\.\d{3})?(?:[+-][0-2]\d:[0-5]\d|Z)?`)

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
	res := strings.Trim(s, "{ }")
	return map[string]any{"val": res}
}

func (id Id) String() string {
	if !id.cmplx {
		return fmt.Sprint(id.val)
	}

	switch val := id.val.(type) {
	case ArrayId:
		res := make([]string, len(val))

		for i, v := range val {
			switch x := v.(type) {
			case time.Time:
				res[i] = x.Format(time.RFC3339Nano)
			default:
				res[i] = fmt.Sprint(x)
			}
		}

		return fmt.Sprintf("['%s']", strings.Join(res, "', '"))
	case ObjectId:
		return fmt.Sprintf("{%s}", id.val)
	default:
		return fmt.Sprintf("⟨%s⟩", id.val)
	}
}
