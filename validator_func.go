package validators

import (
	"context"
	"fmt"
	"net"
	"reflect"
	"time"
	"unicode/utf8"
	"validators/lang"
)

var (
	timeType = reflect.TypeOf(time.Time{})
	//defaultCField = &cField{namesEqual: true}
)

type FuncCtx func(ctx context.Context, fv reflect.Value) bool

type Func func(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error)

var defaultValidator = map[string]Func{
	"required": hasValue,
	"len":      hasLengthOf,
	"min":      hasMinOf,
	"max":      hasMaxOf,
	"eq":       isEq,
	"lt":       isLt,
	"lte":      isLte,
	"gt":       isGt,
	"gte":      isGte,
	"email":    isEmail,
	"number":   isNumber,
	"phone":    isPhone,
	"ipv4":     isIPv4,
	"ipv6":     isIPv6,
	"ip":       isIP,
	"in":       isIn,
	//"datetime": isDatetie,
	//"url":      isUrl,
	//"between":  isBetween,
	"unique": isUnique,
}

// isEq
func isEq(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]
	switch ft.Kind() {
	case reflect.String:
		flag = fv.Interface() != param

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) != p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)
		flag = fv.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		flag = fv.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		flag = fv.Float() == p

	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if flag {
		err = fmt.Errorf("不等于")
	}
	return
}

// IisLt
func isLt(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]

	switch ft.Kind() {
	case reflect.String:
		p := asInt(param)
		flag = int64(utf8.RuneCountInString(fv.String())) < p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) < p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)
		flag = fv.Int() < p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		flag = fv.Uint() < p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		flag = fv.Float() < p

	case reflect.Struct:
		if fv.Type() == timeType {
			flag = fv.Interface().(time.Time).Before(time.Now().UTC())
		}

	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if !flag {
		err = fmt.Errorf("%s不小于%s", title, param)
	}
	return

}

// isLte
func isLte(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]

	switch ft.Kind() {

	case reflect.String:
		p := asInt(param)

		flag = int64(utf8.RuneCountInString(fv.String())) <= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) <= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)
		flag = fv.Int() <= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		flag = fv.Uint() <= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		flag = fv.Float() <= p

	case reflect.Struct:
		if fv.Type() == timeType {
			now := time.Now().UTC()
			t := fv.Interface().(time.Time)
			flag = t.Before(now) || t.Equal(now)
		}
	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if !flag {
		err = fmt.Errorf("%s大于%s", title, param)
	}
	return
}

// isGt
func isGt(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]

	switch ft.Kind() {
	case reflect.String:
		p := asInt(param)
		flag = int64(utf8.RuneCountInString(fv.String())) > p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) > p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)
		flag = fv.Int() > p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		flag = fv.Uint() > p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		flag = fv.Float() > p

	case reflect.Struct:
		if fv.Type() == timeType {
			flag = fv.Interface().(time.Time).After(time.Now().UTC())
		}

	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if !flag {
		err = fmt.Errorf("%s不大于%s", title, param)
	}
	return
}

// isGte
func isGte(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]

	switch ft.Kind() {

	case reflect.String:
		p := asInt(param)
		flag = int64(utf8.RuneCountInString(fv.String())) >= p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) >= p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(param)
		flag = fv.Int() >= p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(param)
		flag = fv.Uint() >= p

	case reflect.Float32, reflect.Float64:
		p := asFloat(param)
		flag = fv.Float() >= p

	case reflect.Struct:
		if fv.Type() == timeType {
			now := time.Now().UTC()
			t := fv.Interface().(time.Time)
			flag = t.After(now) || t.Equal(now)
		}
	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if !flag {
		err = fmt.Errorf("%s小于%s", title, param)
	}
	return
}

// hasLengthOf
func hasLengthOf(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var vInt int64
	var vFloat float64
	if len(params) < 1 {
		err = fmt.Errorf("参数个数有误")
	}
	kind := ft.Kind()
	switch ft.Kind() {
	case reflect.String:
		kind = reflect.Int32
		vInt = int64(utf8.RuneCountInString(fv.String()))
		//fmt.Println("LenString:", fv.String(), kind, vInt)
	case reflect.Slice, reflect.Map, reflect.Array:
		kind = reflect.Int32
		vInt = int64(fv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		kind = reflect.Int64
		vInt = int64(fv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		kind = reflect.Int64
		vInt = int64(fv.Uint())
	case reflect.Float32, reflect.Float64:
		kind = reflect.Float64
		vFloat = float64(fv.Float())
	default:
		panic(fmt.Sprintf("Bad field type %T", fv.Interface()))
	}
	if !checkNumber(kind) {
		err = fmt.Errorf("校验类型不对")
		return
	}

	if len(params) == 1 {
		if kind == reflect.Float64 {
			if asFloat(params[0]) != vFloat {
				err = fmt.Errorf("不等于%f", asFloat(params[0]))
			}
		} else if kind == reflect.Int64 {
			if asInt(params[0]) != vInt {
				err = fmt.Errorf("不等于%d", asInt(params[0]))
			}
		} else {
			if asInt(params[0]) != vInt {
				err = fmt.Errorf("长度不等于%d", asInt(params[0]))
			}
		}
	} else if len(params) >= 1 {

		if params[0] != "_" {
			//fmt.Println("INT32:", fv.String(), kind, vInt, asInt(params[0]))
			if kind == reflect.Float64 {
				if asFloat(params[0]) > vFloat {
					err = fmt.Errorf("小于%f", asFloat(params[0]))
				}
			} else if kind == reflect.Int64 {
				if asInt(params[0]) > vInt {
					err = fmt.Errorf("小于%d", asInt(params[0]))
				}
			} else {

				if asInt(params[0]) > vInt {
					err = fmt.Errorf("长度小于%d", asInt(params[0]))
				}
			}
		}
		if params[1] != "_" {

			if kind == reflect.Float64 {
				if asFloat(params[1]) < vFloat {
					err = fmt.Errorf("大于%f", asFloat(params[1]))
				}
			} else if kind == reflect.Int64 {
				if asInt(params[1]) < vInt {
					err = fmt.Errorf("大于%d", asInt(params[1]))
				}
			} else {
				if asInt(params[1]) < vInt {
					err = fmt.Errorf("长度大于%d", asInt(params[1]))
				}
			}
		}
	}
	return
}

// hasMinOf
func hasMinOf(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	return isGte(ft, fv, title, params...)
}

// hasMaxOf
func hasMaxOf(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	return isLte(ft, fv, title, params...)
}

// hasValue
func hasValue(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	if isZeroValue(fv) {
		err = fmt.Errorf("不能为空")
	}
	return

}

// isIn
func isIn(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var vals []reflect.Value
	var argsI []interface{}
	kind := ft.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		kind = ft.Elem().Kind()
		arrLen := fv.Len()
		for i := 0; i < arrLen; i++ {
			vals = append(vals, fv.Index(i))
		}
	case reflect.Map:
		kind = ft.Elem().Kind()
		keys := fv.MapKeys()
		for _, key := range keys {
			vals = append(vals, fv.MapIndex(key))
		}

	default:
		vals = append(vals, fv)
	}
	if !checkNumber(kind) && !checkBool(kind) && !checkString(kind) {
		err = fmt.Errorf("校验类型不对")
		return
	}
	if len(vals) == 0 {
		err = fmt.Errorf("校验数据不能为空")
	}
	//根据 val 类型将 args 转为对应格式
	for _, param := range params {
		tmpArg, err := parseStr(param, kind)
		if err != nil {
			err = fmt.Errorf("参数与类型不匹配")
		}
		argsI = append(argsI, tmpArg)
	}
	for _, valI := range vals {
		if !InArray(parseReflectV(valI, kind), argsI) {
			err = fmt.Errorf("%v不在指定范围11:%v", valI, params)
		}
	}
	return
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func isEmail(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	if !emailRegex.MatchString(fv.String()) {
		err = fmt.Errorf("非Email:%s", fv.String())
	}
	return
}

// isPhone
func isPhone(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	if !phoneRegex.MatchString(fv.String()) {
		err = fmt.Errorf(trans(lang.ValidIsPhone), fv.String())
	}
	return
}

// isNumber
func isNumber(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	switch ft.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64:
		return
	default:
		if !numberRegex.MatchString(fv.String()) {
			err = fmt.Errorf(trans(lang.ValidIsNumber), fv.String())
		}
		return
	}
}

// isIPv4
func isIPv4(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() != nil {
		err = fmt.Errorf("非IPv4")
	}
	return
}

// isIPv6
func isIPv6(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To16() != nil {
		err = fmt.Errorf("非IPv6")
	}
	return
}

// isIP
func isIP(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil {
		err = fmt.Errorf("非IP")
	}
	return
}

// isUnique
func isUnique(ft reflect.Type, fv reflect.Value, title string, params ...string) (err error) {
	var flag bool
	v := reflect.ValueOf(1)
	switch ft.Kind() {
	case reflect.Slice, reflect.Array:

		switch fv.Type().Elem().Kind() {
		case reflect.String:
			m := make(map[string]int)
			for i := 0; i < fv.Len(); i++ {
				m[fv.Index(i).String()] = 1
			}
			flag = fv.Len() != len(m)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			m := make(map[int64]int)
			for i := 0; i < fv.Len(); i++ {
				m[fv.Index(i).Int()] = 1
			}
			flag = fv.Len() != len(m)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			m := make(map[uint64]int)
			for i := 0; i < fv.Len(); i++ {
				m[fv.Index(i).Uint()] = 1
			}
			flag = fv.Len() != len(m)
		case reflect.Float32, reflect.Float64:
			m := make(map[float64]int)
			for i := 0; i < fv.Len(); i++ {
				m[fv.Index(i).Float()] = 1
			}
			flag = fv.Len() != len(m)
		}
		//m := reflect.MakeMap(reflect.MapOf(fv.Type().Elem(), v.Type()))
		//for i := 0; i < fv.Len(); i++ {
		//	m.SetMapIndex(fv.Index(i), v)
		//}
		//flag = fv.Len() != m.Len()
	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(fv.Type().Elem(), v.Type()))
		for _, k := range fv.MapKeys() {
			m.SetMapIndex(fv.MapIndex(k), v)
		}
		flag = fv.Len() != m.Len()
	default:
		panic(fmt.Sprintf("唯一值校验类型不支持 %T", fv.Interface()))
	}
	if flag {
		err = fmt.Errorf("非唯一值")
	}
	return

}
