package validators

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

const (
	STRUCT_EMPTY            = "struct %v is empty"
	VALIDATOR_VALUE_SIGN    = "="
	VALIDATOR_RANGE_SPLIT   = ","
	VALIDATOR_IGNORE_SIGN   = "_"
	VALIDATOR_MUTIPLE_SPLIT = ";"
)

var errorMsg map[string][]string
var lang = "zh"

type Validator struct {
	ValidTag   string
	TitleTag   string
	lazy       bool
	allowEmpty bool
	validator  map[string]FuncCtx
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
func (v *Validator) SetLang(l string) *Validator {
	lang = l
	return v
}

// SetLazy 设置语言
func (v *Validator) SetLazy(lazy bool) *Validator {
	v.lazy = lazy
	return v
}

// RegisterValidator 注册新验证规则
func (v *Validator) RegisterValidator(validatorK string, validator FuncCtx) *Validator {
	v.validator[validatorK] = validator
	return v
}

// RegisterValidators 批量注册新验证规则
func (v *Validator) RegisterValidators(validatorMap map[string]FuncCtx) *Validator {
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
			title := fieldTypeInfo.Tag.Get(v.TitleTag)
			if tag != "" {
				//fmt.Println("ffff:", fv, ft, fieldTypeInfo)
				//没有配置 required，并且 field 为 0 值的，直接跳过
				isZeroValue := isZeroValue(fv)
				if isZeroValue && !strings.Contains(tag, "required") && !v.allowEmpty {
					continue
				}
				if title == "" {
					title = fieldTypeInfo.Name
				}
				errArr = v.validateRule(ft, fv, title, tag)
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
					//fmt.Println("tmpParentKey:",fv.Index(i).Interface(),tmpParentKey)tmpParentKey
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

func (v *Validator) validateRule(typeObj reflect.Type, typeValue reflect.Value, title string, rulerString string) (errs []error) {
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
			errs = append(errs, fmt.Errorf(trans(ValidNotExist), ruler))
			if v.lazy == false {
				return
			}
			continue
		}

		// 验证规则
		//fmt.Println("typeValue", typeValue.String())
		err := v.validator[ruler](typeObj, typeValue, title, params...)
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
