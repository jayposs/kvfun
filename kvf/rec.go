// The rec.go file contains funcs that perform actions using a record - []byte.
// Funcs recGetStr and recGetInt return a field's value from the record.
// Func recFind determines if the record meets specified FindConditions.

package kvf

import (
	"cmp"
	"log"
	"strings"

	"github.com/valyala/fastjson"
)

const toLower = true // used by recFind/recGetStr below to always return lower case string value

// FindCondition Ops
const (
	Contains int = iota
	Matches
	StartsWith
	LessThanStr
	GreaterThanStr
	LessThan    // int
	GreaterThan // int
	EqualTo     // int
)

// Type FindCondition is used a parameter to recFind func.
// The Op code determines if ValStr or ValInt is used for comparison.
type FindCondition struct {
	Fld    string // field name in Rec containing compare value
	Op     int    // see constants above
	ValStr string // for Ops Matches, StartsWith, Contains, LessThanStr, GreaterThanStr
	ValInt int    // for Ops EqualTo, LessThan, GreaterThan
}

// NOTE - The compare string value (ValStr) is automatically converted to lower case, so caller doesn't need to convert.
// If this behaviour is not valid for your use case, code must be changed.

// Func recGetStr returns the string value associated with a field in the record.
// If optional parm toLower is true, then return lowerCase of value.
func recGetStr(rec []byte, fld string, toLower ...bool) string {
	val := fastjson.GetString(rec, fld)
	if len(toLower) > 0 && toLower[0] {
		return strings.ToLower(val)
	}
	return val
}

// Func recGetInt returns the int value associated with a field in the record.
func recGetInt(rec []byte, fld string) int {
	return fastjson.GetInt(rec, fld)
}

// Func recFind determines if rec values meet all find conditions.
func recFind(rec []byte, conditions []FindCondition) bool {
	var conditionMet bool
	var n int                        // compare result  1:greater, -1:less, 0:equal
	var compareVal, recValStr string // only used for strings, to support StartsWith and Contains ops
	for _, condition := range conditions {
		conditionMet = false
		switch condition.Op {
		case Contains, Matches, StartsWith, LessThanStr, GreaterThanStr: // string comparison
			compareVal = strings.ToLower(condition.ValStr)
			recValStr = recGetStr(rec, condition.Fld, toLower)
			n = cmp.Compare(recValStr, compareVal)
		case EqualTo, LessThan, GreaterThan: // int comparison
			recVal := recGetInt(rec, condition.Fld)
			n = cmp.Compare(recVal, condition.ValInt)
		default:
			log.Println("invalid find op", condition.Op)
			return false
		}
		switch condition.Op {
		case Matches, EqualTo:
			if n == 0 {
				conditionMet = true
			}
		case LessThan, LessThanStr:
			if n == -1 {
				conditionMet = true
			}
		case GreaterThan, GreaterThanStr:
			if n == 1 {
				conditionMet = true
			}
		case StartsWith:
			if strings.HasPrefix(recValStr, compareVal) {
				conditionMet = true
			}
		case Contains:
			if strings.Index(recValStr, compareVal) > -1 {
				conditionMet = true
			}
		}
		if !conditionMet {
			return false // condition was not met, end recFind
		}
	}
	return true // no condition check returned false
}
