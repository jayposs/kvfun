package kvf

import (
	"cmp"
	"log"
	"strings"

	"github.com/valyala/fastjson"
)

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

// NOTE - The Op code determines if ValStr or ValInt is used for comparison
type FindCondition struct {
	Fld    string // fld in Rec containing compare value
	Op     int    // see constants above
	ValStr string // for Ops Matches, StartsWith, Contains, LessThanStr, GreaterThanStr
	ValInt int    // for Ops EqualTo, LessThan, GreaterThan
}

func recGetStr(rec []byte, fld string) string {
	return fastjson.GetString(rec, fld)
}

func recGetInt(rec []byte, fld string) int {
	return fastjson.GetInt(rec, fld)
}

// recFind determines if rec values meet all find conditions
func recFind(rec []byte, conditions []FindCondition) bool {
	var ok bool
	var n int                        // compare result  1:greater, -1:less, 0:equal
	var compareVal, recValStr string // only used for strings
	for _, condition := range conditions {
		ok = false
		switch condition.Op {
		case Contains, Matches, StartsWith, LessThanStr, GreaterThanStr: // string comparison
			compareVal = strings.ToLower(condition.ValStr)
			recValStr = recGetStr(rec, condition.Fld)
			recValStr = strings.ToLower(recValStr)
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
				ok = true
			}
		case LessThan, LessThanStr:
			if n == -1 {
				ok = true
			}
		case GreaterThan, GreaterThanStr:
			if n == 1 {
				ok = true
			}
		case StartsWith:
			if strings.HasPrefix(recValStr, compareVal) {
				ok = true
			}
		case Contains:
			if strings.Index(recValStr, compareVal) > -1 {
				ok = true
			}
		}
		if !ok {
			return false // condition was not met
		}
	}
	return true // no condition check returned false
}

/*

			switch condition.Op {
			case Matches, LessThanStr, GreaterThanStr:
				n := cmp.Compare(recVal, compareVal)
				if condition.Op == Matches && n == 0 {
					ok = true
				} else if condition.Op == GreaterThanStr && n == 1 {
					ok = true
				} else if condition.Op == LessThanStr && n == -1 {
					ok = true
				}
			case StartsWith:
				if strings.HasPrefix(recVal, compareVal) {
					ok = true
				}
			case Contains:
				if strings.Index(recVal, compareVal) > -1 {
					ok = true
				}
			}
			if !ok {
				return false // condition was not met
			}

		}
		default: // for all int ops
			recVal := recGetInt(rec, condition.Fld)
			n := cmp.Compare(recVal, condition.ValInt)
			switch condition.Op {
			case EqualTo:
				if n == 0 {
					ok = true
				}
			case LessThan:
				if n == -1 {
					ok = true
				}
			case GreaterThan:
				if n == 1 {
					ok = true
				}
			default:
				log.Println("invalid find op", condition.Op)
				return false
			}
		}
		if !ok {
			return false // condition was not met
		}
	}
	return true
*/

/*
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
*/
