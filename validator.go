package validators

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"validators/lang"
)

const (
	STRUCT_EMPTY            = "struct %v is empty"
	VALIDATOR_VALUE_SIGN    = "="
	VALIDATOR_RANGE_SPLIT   = ","
	VALIDATOR_IGNORE_SIGN   = "_"
	VALIDATOR_MUTIPLE_SPLIT = ";"

	//邮箱验证正则
	MAIL_REG = `\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`
	//url验证正则
	URL_REG = `^(http|ftp|https):\/\/[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?$`
	//是否为整数正则
	INTEGER_REG = `^(-)?[0-9]+$`
	//是否为float正则
	FLOAT_REG = `^(-)?[0-9]+(.[0-9]+)$`
	//年正则
	YEAR_REG = `(19|2[0-4])\d{2}`
	//月正则
	MONTH_REF = `(10|11|12|0[1-9])`
	//日正则
	DAY_REF = `(30|31|0[1-9]|[1-2][0-9])`
	//小时正则
	HOUR_REF = `(20|21|22|23|[0-1]\d)`
	//分钟正则
	MINUTE_REF = `([0-5]\d)`
	//秒正则
	SECOND_REF = `([0-5]\d)`
)

/****************************************************
 * range 验证错误提示 map
 ****************************************************/
var stringErrorMap = map[string]string{
	"lt":      "[name] should be less than [max] chars long",
	"eq":      "[name] should be eq [min] chars long",
	"gt":      "[name] should be great than [min] chars long",
	"between": "[name] should be betwween [min] and [max] chars long",
}

var numberErrorMap = map[string]string{
	"lt":      "[name] should be less than [max]",
	"eq":      "[name] should be eq [min]",
	"gt":      "[name] should be great than [min]",
	"between": "[name] should be betwween [min] and [max]",
}

var arrayErrorMap = map[string]string{
	"lt":      "array [name] length should be less than [max]",
	"eq":      "array [name] length should be eq [min]",
	"gt":      "array [name] length should be great than [min]",
	"between": "array [name] length should be betwween [min] and [max]",
}

/****************************************************
 * range 验证错误提示 map
 ****************************************************/

var errorMsg map[string][]string
var Lang = "zh"

type Validator struct {
	ValidTag   string
	TitleTag   string
	lazy       bool
	allowEmpty bool
	validator  map[string]Func
}

func New() *Validator {
	return &Validator{
		ValidTag:   "validate",
		TitleTag:   "title",
		lazy:       true,
		allowEmpty: true,
		validator:  defaultValidator,
	}
}

// SetValidTag 设置校验tag
func (v *Validator) SetValidTag(tag string) *Validator {
	v.ValidTag = tag
	return v
}

// SetTitleTag 设置字段标题tag
func (v *Validator) SetTitleTag(titleTag string) *Validator {
	v.TitleTag = titleTag
	return v
}

// SetAllowEmpty 允许空结构
func (v *Validator) SetAllowEmpty(skip bool) *Validator {
	v.allowEmpty = skip
	return v
}

// SetLang 设置语言
func (v *Validator) SetLang(lang string) *Validator {
	Lang = lang

	return v
}

// SetLazy 设置语言
func (v *Validator) SetLazy(lazy bool) *Validator {
	v.lazy = lazy
	return v
}

// RegisterValidator 注册新验证规则
func (v *Validator) RegisterValidator(validatorK string, validator Func) *Validator {
	v.validator[validatorK] = validator
	return v
}

// RegisterValidators 批量注册新验证规则
func (v *Validator) RegisterValidators(validatorMap map[string]Func) *Validator {
	for validatorK, validatorV := range validatorMap {
		v.validator[validatorK] = validatorV
	}
	return v
}

// LazyValidate 延迟校验输出
func (v *Validator) LazyValidate(s interface{}) (err error) {
	syncMap := &sync.Map{}
	parentKey := v.ValidTag
	errArr := v.validate(s, true, syncMap, parentKey)
	syncMap = nil
	if errArr != nil {
		err = errArr[0]
	}
	return
}

// Struct 校验结构体
func (v *Validator) Struct(s interface{}) (err []error) {
	syncMap := &sync.Map{}
	parentKey := v.ValidTag
	err = v.validate(s, false, syncMap, parentKey)
	syncMap = nil
	return
}

// Value 校验值
func (v *Validator) Value(s interface{}) (err []error) {
	syncMap := &sync.Map{}
	parentKey := v.ValidTag
	err = v.validate(s, false, syncMap, parentKey)
	syncMap = nil
	return
}

func (v *Validator) validate(s interface{}, lazyFlag bool, syncMap *sync.Map, parentKey string) (errs []error) {
	var errArr []error
	rt := reflect.TypeOf(s)
	rv := reflect.ValueOf(s)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
		rv = rv.Elem()
	}
	switch rt.Kind() {
	case reflect.Slice, reflect.Array:
		//判断是否需要递归
		if ok, fieldNum := checkArrayValueIsMulti(rv); ok {
			for i := 0; i < fieldNum; i++ {
				tmpParentKey := fmt.Sprintf("%v_%v", parentKey, i)
				errArr = v.validate(rv.Index(i).Interface(), lazyFlag, syncMap, tmpParentKey)
				if len(errArr) > 0 {
					errs = append(errs, errArr...)
					if lazyFlag {
						return
					}
					continue
				}
			}
		} else {
			//不需要递归
			fmt.Println("======不递归=====>", rv)
		}
		break
	case reflect.Struct:
		numField := rv.NumField()
		if numField <= 0 {
			if v.allowEmpty {
				return
			}
			errs = append(errs, fmt.Errorf(STRUCT_EMPTY, rt.Name()))
			return
		}

		for i := 0; i < numField; i++ {
			fv := rv.Field(i)
			ft := rt.Field(i).Type
			fieldTypeInfo := rv.Type().Field(i)
			fieldType := fv.Type().Kind()
			tag := fieldTypeInfo.Tag.Get(v.ValidTag)
			if tag != "" {
				//fmt.Println("ffff:", fv, ft, fieldTypeInfo)
				//没有配置 required，并且 field 为 0 值的，直接跳过
				isZeroValue := isZeroValue(fv)
				if isZeroValue && !strings.Contains(tag, "required") && !v.allowEmpty {
					continue
				}
				errArr = v.validateRule(ft, fv, tag)
				if len(errArr) > 0 {
					errs = append(errs, errArr...)
					if lazyFlag {
						return
					}
					continue
				}
			}
			//判断是否需要递归
			if ok, fieldNum := checkArrayValueIsMulti(fv); ok {
				for i := 0; i < fieldNum; i++ {
					//tmpParentKey := fmt.Sprintf("%v_%v", parentKey, fieldTypeInfo.Name)
					//tmpParentKey := fmt.Sprintf("%v_%v", parentKey, fieldTypeInfo.Name)
					//fmt.Println("tmpParentKey:",fv.Index(i).Interface(),tmpParentKey)
					errArr = v.validate(fv.Index(i).Interface(), lazyFlag, syncMap, parentKey)
					if len(errArr) > 0 {
						errs = append(errs, errArr...)
						if lazyFlag {
							return
						}
						continue
					}
				}
			}

			if fieldType == reflect.Struct {
				tmpParentKey := fmt.Sprintf("%v_%v", parentKey, fieldTypeInfo.Name)
				errArr = v.validate(fv.Interface(), lazyFlag, syncMap, tmpParentKey)
				if len(errArr) > 0 {
					errs = append(errs, errArr...)
					if lazyFlag {
						return
					}
					continue
				}
			}
		}
	}
	return
}

func (v *Validator) validateRule(typeObj reflect.Type, typeValue reflect.Value, rulerString string) (errs []error) {
	//typeObj := reflect.TypeOf(s)
	//typeValue := reflect.ValueOf(s)
	rulers := strings.Split(rulerString, VALIDATOR_MUTIPLE_SPLIT)
	for _, ruler := range rulers {
		var params []string
		//查找是否含有赋值符号
		num := strings.Index(ruler, VALIDATOR_VALUE_SIGN)
		//不等于 -1, 表示含有"="
		if num != -1 {
			params = strings.Split(ruler[num+1:], VALIDATOR_RANGE_SPLIT)
			ruler = ruler[0:num]
		}
		// 判断验证规则是否存在
		if _, ok := v.validator[ruler]; !ok {
			errs = append(errs, fmt.Errorf(trans(lang.ValidNotExist), ruler))
			if v.lazy == false {
				return
			}
			continue
		}

		// 验证规则
		//fmt.Println("typeValue", typeValue.String())
		err := v.validator[ruler](typeObj, typeValue, params...)
		if err != nil {
			errs = append(errs, err)
			if v.lazy == false {
				return
			}
			continue
		}

	}
	return
}
