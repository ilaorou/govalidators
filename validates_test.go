package validators

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRequired(t *testing.T) {
	validator := New()

	testString := []struct {
		param    string `validate:"required" title:"学生ID++++++"`
		expected bool
	}{
		{"1", true},
		{"ssss", true},
		{"", false},
	}
	for _, test := range testString {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}

	testInt64 := []struct {
		param    int64 `validate:"required"`
		expected bool
	}{
		{1000, true},
		{0, false},
	}
	for _, test := range testInt64 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}

	testFloat64 := []struct {
		param    float64 `validate:"required"`
		expected bool
	}{
		{10000.23, true},
		{0, false},
		{0.00, false},
	}

	for _, test := range testFloat64 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}

	testSlice := []struct {
		param    []string `validate:"required"`
		expected bool
	}{
		{[]string{"ss", "aa", "ss"}, true},
		{[]string{}, false},
	}
	for _, test := range testSlice {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}

	testMap := []struct {
		param    map[string]string `validate:"required"`
		expected bool
	}{
		{map[string]string{"a": "aa", "b": "bb"}, true},
		{map[string]string{}, false},
	}
	for _, test := range testMap {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}
}

func TestIn2(t *testing.T) {
	validator := New()
	testBetween := []struct {
		param    int64 `validate:"in=-10,20"`
		expected bool
	}{
		{-10, true},
		{15, false},
		{20, true},
		{-11, false},
		{21, false},
	}
	for _, test := range testBetween {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}

	testEqual := []struct {
		param    string `validate:"in=-10"`
		expected bool
	}{
		{"-10", true},
		{"10", false},
		{"15", false},
		{"9", false},
	}
	for _, test := range testEqual {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}
	testGt := []struct {
		param    int `validate:"gt=10;lt=13" title:"abc"`
		expected bool
	}{
		{-10, false},
		{11, false},
		{15, false},
		{9, false},
	}
	for _, test := range testGt {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected required,value %v,err %v,expected  %v", test.param, err, test.expected)
		}
	}
}
func TestIn(t *testing.T) {
	validator := New()
	testIn1 := []struct {
		param    map[string]string `validate:"in=a,b,c"`
		expected bool
	}{
		{map[string]string{"k1": "a", "k2": "c"}, true},
		{map[string]string{"k1": "a", "k2": "d"}, false},
		{map[string]string{}, false},
	}
	for _, test := range testIn1 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected in,value %v,err %v", test.param, err)
		}
	}
	testIn2 := []struct {
		param    []int64 `validate:"in=1,20,01"`
		expected bool
	}{
		{[]int64{1, 20}, true},
		{[]int64{1, 10}, false},
		{[]int64{}, false},
	}

	for _, test := range testIn2 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected in,value %v,err %v", test.param, err)
		}
	}

	testIn3 := []struct {
		param    []float64 `validate:"in=1.11,20.22,01.10"`
		expected bool
	}{
		{[]float64{1.11, 20.22, 1.1}, true},
		{[]float64{1.12, 20.33}, false},
		{[]float64{}, false},
	}

	for _, test := range testIn3 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected in,value %v,err %v", test.param, err)
		}
	}

	testIn4 := []struct {
		param    string `validate:"in=a,b,c"`
		expected bool
	}{
		{"a", true},
		{"d", false},
		{"", false},
	}
	for _, test := range testIn4 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected in,value %v,err %v", test.param, err)
		}
	}
}

func TestLen(t *testing.T) {
	validator := New()
	testBetween := []struct {
		param    string `validate:"len=1,5"`
		expected bool
	}{
		{"s", true},
		{"中文", true},
		{"sssss", true},
		{"", false},
		{"ssssss", false},
	}
	for _, test := range testBetween {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testEqual := []struct {
		param    string `validate:"len=5"`
		expected bool
	}{
		{"sssss", true},
		{"s", false},
		{"", false},
		{"ssssss", false},
	}
	for _, test := range testEqual {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testGreatThen := []struct {
		param    string `validate:"len=2,_"`
		expected bool
	}{
		{"sss", true},
		{"sssssssss", true},
		{"s", false},
		{"", false},
	}
	for _, test := range testGreatThen {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testLessThan := []struct {
		param    string `validate:"len=_,3"`
		expected bool
	}{
		{"", true},
		{"sss", true},
		{"ssss", false},
		{"ssssssssss", false},
	}
	for _, test := range testLessThan {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testInt := []struct {
		param    int `validate:"len=1,500"`
		expected bool
	}{
		{1, true},
		{123, true},
		{12345, false},
		{1234567, false},
	}
	for _, test := range testInt {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testFloat := []struct {
		param    float64 `validate:"len=1,500"`
		expected bool
	}{
		{1, true},
		{123.23, true},
		{1234.21, false},
		{123456.12, false},
	}
	for _, test := range testFloat {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected string,value %v,err %v", test.param, err)
		}
	}

	testArray := []struct {
		param    []string `validate:"len=1,5"`
		expected bool
	}{
		{[]string{"s"}, true},
		{[]string{"s", "s", "s", "s", "s"}, true},
		{[]string{}, false},
		{[]string{"s", "s", "s", "s", "s", "s"}, false},
	}
	for _, test := range testArray {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected array,value %v,err %v", test.param, err)
		}
	}

	testMap := []struct {
		param    map[string]string `validate:"len=_,3"`
		expected bool
	}{
		{map[string]string{"k1": "s1"}, true},
		{map[string]string{}, true},
		{map[string]string{"k1": "s1", "k2": "s2", "k3": "s3", "k4": "s4"}, false},
	}
	for _, test := range testMap {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected array,value %v,err %v", test.param, err)
		}
	}
}

type uniqT struct {
	Name string `validate:"unique"`
}

func TestUnique(t *testing.T) {
	validator := New()

	//testUnique2 := []struct {
	//	Param    []uniqT
	//	expected bool
	//}{
	//	{[]uniqT{{"a"}, uniqT{"b"}, uniqT{"c"}, uniqT{"d"}, uniqT{"e"}}, true},
	//	{[]uniqT{{"a"}, uniqT{"b"}, uniqT{"c"}, uniqT{"d"}, uniqT{"a"}}, false},
	//}
	//for _, test := range testUnique2 {
	//	err := validator.Struct(test)
	//	if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
	//		t.Errorf("Expected unique,value %v,err %v", test.Param, err)
	//	}
	//}

	testUnique0 := []struct {
		param    []int `validate:"unique"`
		expected bool
	}{
		{[]int{1, 2, 3, 4, 5}, true},
	}
	for _, test := range testUnique0 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected unique,value %v,err %v", test.param, err)
		}
	}

	testUnique1 := []struct {
		param    []string `validate:"unique"`
		expected bool
	}{
		{[]string{"a", "b", "c", "d", "e"}, true},
		{[]string{"a", "b", "c", "d", "a"}, false},
	}
	for _, test := range testUnique1 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected unique,value %v,err %v", test.param, err)
		}
	}

}

func TestDateTime(t *testing.T) {
	validator := New()
	testDateTime1 := []struct {
		param    string `validate:"datetime=Y m"`
		expected bool
	}{
		{"2012 12", true},
		{"2012 13", false},
		{"2012-12", false},
	}
	for _, test := range testDateTime1 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected unique,value %v,err %v", test.param, err)
		}
	}

	testDateTime2 := []struct {
		param    string `validate:"datetime=Y-m-d H-i"`
		expected bool
	}{
		{"2012-12-31 13-01", true},
		{"2012-01-01 11-12", true},
		{"2012 13", false},
		{"2012-12-32 11:12", false},
		{"2012-13-01 11-12", false},
		{"2012-13-01 24-12", false},
	}
	for _, test := range testDateTime2 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected unique,value %v,err %v", test.param, err)
		}
	}
}

func TestEmail(t *testing.T) {
	validator := New()
	testEmail := []struct {
		param    string `validate:"email"`
		expected bool
	}{
		{"zl111sdaaj@sina.com", true},
		{"1232920@qq.com", true},
		{"2012-12@qq.com.cn", true},
		{"abcde.com", false},
		{"@abcde.com", false},
	}
	for _, test := range testEmail {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected unique,value %v,err %v", test.param, err)
		}
	}
}

type Class struct {
	Cid       int64  `validate:"required;integer=1,1000000"`
	Cname     string `validate:"required;string=1,5;unique"`
	BeginTime string `validate:"required;datetime=H:i"`
}

type Student struct {
	Uid          int64    `validate:"required;len=1,1000000" title:"学生ID"`
	Name         string   `validate:"required;len=1,5" title:"姓名"`
	Age          int64    `validate:"required;in=10,30"`
	Sex          string   `validate:"required;in=male,female"`
	Email        string   `validate:"email"`
	Phone        string   `validate:"phone"`
	PersonalPage string   `validate:"url"`
	Hobby        []string `validate:"array=_,2;unique;in=swimming,running,drawing"`
	CreateTime   string   `validate:"datetime"`
	Class        []Class  `validate:"array=1,3"`
	expected     bool
}

type UserStringValidator struct {
	EMsg string
}

func (self *UserStringValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {
	fmt.Println("=====", "UserStringValidator")
	return true, nil
}

func userMethod(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {
	fmt.Println("=====", "userMethod")
	return true, nil
}

func TestMuti1(t *testing.T) {
	validator := New()
	//validator.RegisterValidators(map[string]interface{}{
	//	"string": &StringValidator{
	//		Range: Range{
	//			RangeEMsg: map[string]string{
	//				"between": "[name] 长度必须在 [min] 和 [max] 之间",
	//			},
	//		},
	//	},
	//	"um":  userMethod,
	//	"usv": &UserStringValidator{},
	//})
	testMuti1 := []*Student{
		&Student{
			Uid:          123456,
			Name:         "张三",
			Age:          12,
			Sex:          "male",
			Email:        "123456@qq.com",
			Phone:        "123456",
			PersonalPage: "http://www.abcd.com",
			Hobby:        []string{"swimming", "running"},
			CreateTime:   "2018-03-03 05:00:00",
			expected:     true,
			Class: []Class{
				Class{
					Cid:       12345,
					Cname:     "语文",
					BeginTime: "13:00",
				},
				Class{
					Cid:       22345,
					Cname:     "数学",
					BeginTime: "13:00",
				},
			},
		},
		&Student{
			Uid:          1234567,
			Name:         "张三1111",
			Age:          31,
			Sex:          "male1",
			Email:        "@qq.com",
			PersonalPage: "www.abcd.com",
			Hobby:        []string{"swimming", "singing"},
			CreateTime:   "2018-03-03 05:60:00",
			expected:     false,
			Class: []Class{
				Class{
					Cid:       12345678,
					Cname:     "语文",
					BeginTime: "13:00",
				},
				Class{
					Cid:       22345678,
					Cname:     "数学",
					BeginTime: "13:00",
				},
				Class{
					Cid:       32345678,
					Cname:     "数学",
					BeginTime: "13:60",
				},
			},
		},
	}
	for _, test := range testMuti1 {
		err := validator.Struct(test)
		if (err != nil && test.expected == true) || (err == nil && test.expected != true) {
			t.Errorf("Expected muti,value %v,err %v", test, err)
		}
	}
}
