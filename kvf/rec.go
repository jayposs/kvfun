package kvf

import (
	"log"
	"strings"

	"github.com/valyala/fastjson"
)

// FindCondition Ops
const (
	Contains int = iota
	Matches
	StartsWith
	LessThan
	GreaterThan
	EqualTo
)

type FindCondition struct {
	Fld    string // fld in Rec containing compare value
	Op     int    // see constants above
	ValStr string // for Ops Matches, StartsWith, Contains
	ValInt int    // for Ops EqualTo, LessThan, GreaterThan
}

func recGetStr(rec []byte, fld string) string {
	return fastjson.GetString(rec, fld)
}

func recGetInt(rec []byte, fld string) int {
	return fastjson.GetInt(rec, fld)
}

func recFind(rec []byte, conditions []FindCondition) bool {
	for _, condition := range conditions {
		ok := false
		switch condition.Op {
		case Contains, Matches, StartsWith: // String Comparison
			compareVal := strings.ToLower(condition.ValStr)
			recVal := recGetStr(rec, condition.Fld)
			recVal = strings.ToLower(recVal)
			if condition.Op == Matches && recVal == compareVal {
				ok = true
			} else if condition.Op == StartsWith && strings.HasPrefix(recVal, compareVal) {
				ok = true
			} else if condition.Op == Contains && strings.Index(recVal, compareVal) > -1 {
				ok = true
			}
		case LessThan, GreaterThan, EqualTo:
			recVal := recGetInt(rec, condition.Fld)
			if condition.Op == LessThan && recVal < condition.ValInt {
				ok = true
			} else if condition.Op == GreaterThan && recVal > condition.ValInt {
				ok = true
			} else if condition.Op == EqualTo && recVal == condition.ValInt {
				ok = true
			}
		default:
			log.Println("invalid find op", condition.Op)
			return false
		}
		if !ok {
			return false // condition was not met
		}
	}
	return true
}
