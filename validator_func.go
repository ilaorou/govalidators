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

type Func func(ft reflect.Type, fv reflect.Value, params ...string) (err error)

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
	//"unique":   isUnique,
}

func isEq(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func isLt(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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
	if flag {
		err = fmt.Errorf("不等于")
	}
	return

}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLte(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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
	if flag {
		err = fmt.Errorf("不等于")
	}
	return
}

// IsGt is the validation function for validating if the current field's value is greater than the param's value.
func isGt(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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
	if flag {
		err = fmt.Errorf("不等于")
	}
	return
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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
	if flag {
		err = fmt.Errorf("不等于")
	}
	return
}

// HasLengthOf is the validation function for validating if the current field's value is equal to the param's value.
func hasLengthOf(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	var flag bool
	if len(params) != 1 {
		err = fmt.Errorf("参数个数有误")
	}
	param := params[0]
	switch ft.Kind() {
	case reflect.String:
		p := asInt(param)
		flag = int64(utf8.RuneCountInString(fv.String())) == p

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(param)
		flag = int64(fv.Len()) == p

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

// HasMinOf
func hasMinOf(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	return isGte(ft, fv, params...)
}

// HasMaxOf
func hasMaxOf(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	return isLte(ft, fv, params...)
}

// HasValue
func hasValue(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	if isZeroValue(fv) {
		err = fmt.Errorf("不能为空")
	}
	return

}

// isIn
func isIn(ft reflect.Type, fv reflect.Value, params ...string) (err error) {

	var valsI []reflect.Value
	var argsI []interface{}
	kind := ft.Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		kind = ft.Elem().Kind()
		arrLen := fv.Len()
		for i := 0; i < arrLen; i++ {
			valsI = append(valsI, fv.Index(i))
		}
	case reflect.Map:
		kind = ft.Elem().Kind()
		keys := fv.MapKeys()
		for _, key := range keys {
			valsI = append(valsI, fv.MapIndex(key))
		}

	default:
		valsI = append(valsI, fv)
		//err = fmt.Errorf("校验类型不对")
		//return
	}
	if !checkBool(kind) && !checkNumber(kind) && !checkString(kind) {
		err = fmt.Errorf("校验类型不对")
		return
	}
	if len(valsI) == 0 {
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
	for _, valI := range valsI {
		if !InArray(parseReflectV(valI, kind), argsI) {
			err = fmt.Errorf("%v不在指定范围11:%v", valI, params)
		}
	}
	//if flag {
	//	err = fmt.Errorf("%v不在指定范围11:%v", v, params)
	//}
	return
}

// IsEmail is the validation function for validating if the current field's value is a valid email address.
func isEmail(ft reflect.Type, fv reflect.Value, params ...string) (err error) {

	if !emailRegex.MatchString(fv.String()) {
		err = fmt.Errorf("非Email:%s", fv.String())
	}
	return
}

// isPhone
func isPhone(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	if !phoneRegex.MatchString(fv.String()) {
		err = fmt.Errorf(trans(lang.ValidIsPhone), fv.String())
	}
	return
}

// isNumber
func isNumber(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
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

// IsIPv4
func isIPv4(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To4() != nil {
		err = fmt.Errorf("非IPv4")
	}
	return
}

// IsIPv6
func isIPv6(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil || ip.To16() != nil {
		err = fmt.Errorf("非IPv6")
	}
	return
}

// IsIP
func isIP(ft reflect.Type, fv reflect.Value, params ...string) (err error) {
	ip := net.ParseIP(fv.String())
	if ip == nil {
		err = fmt.Errorf("非IP")
	}
	return
}
