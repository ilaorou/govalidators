validators
===========
golang初学者，在项目开发过程中造了一个简单的验证器轮子，欢迎大大们提宝贵建议和指导

### 安装
  go get github.com/ilaorou/validators

### 导入
```go
import "github.com/ilaorou/validators"
```

### 基本使用方式
为 struct 指定验证器
```go
package main

import "github.com/ilaorou/validators"

type Class struct {
  Cid       int64  `validate:"required;len=1,1000000"`
  Cname     string `validate:"required;len=1,5;unique"`
  BeginTime string `validate:"required;datetime=H:i"`
}

type Student struct {
  Uid          int64    `validate:"required;len=1,1000000"`
  Name         string   `validate:"required;len=1,5"`
  Age          int64    `validate:"required;len=10,30"`
  Sex          string   `validate:"required;in=male,female"`
  Email        string   `validate:"email"`
  PersonalPage string   `validate:"url"`
  Hobby        []string `validate:"array=_,2;unique;in=swimming,running,drawing"`
  CreateTime   string   `validate:"datetime"`
  Class        []Class  `validate:"array=1,3"`
}
```
验证
```go
validator := validators.New()
if err := validator.Validate(student); err != nil {
  fmt.Println(err)
}
```

### 自定义验证器

##### 1.支持自定义函数，必须是 ValidatorF 类型，ValidatorF 类型如下
```go
type ValidatorF func(params map[string]interface{}, val reflect.Value, args ...string) (bool, error)
```
自定义函数
```go
func validationMethod(params map[string]interface{}, val reflect.Value, args ...string) (bool, error){
  fmt.Println("validationMethod")
  ...
  return true, nil
}
```
##### 2.支持自定义struct，必须实现 Validator 接口，Validator 接口如下
```go
type Validator interface {
  Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error)
}
```
自定义struct
```go
type UserValidator struct {
  EMsg string
}

func (self *UserValidator) Validate(params map[string]interface{}, val reflect.Value, args ...string) (bool, error) {
  fmt.Println("UserValidator")
  return true, nil
}
```
##### 3.定义好验证器后，初始化验证器
```go
validator := validators.New()
validator.SetValidators(map[string]interface{}{
  "user" : &UserValidator{},
  "vm" : validationMethod,
})
```
##### 4.在需要验证的字段中，增加自定义验证器
```go
Email        string   `validate:"email;user;vm"`
```
##### 5.验证
```go
if err := validator.Validate(student); err != nil {
  fmt.Println(err)
}
```
##### 6.也可以对现有的验证器进行参数设置
```go
validator := validators.New()
validator.SetValidators(map[string]interface{}{
  "string": &validators.StringValidator{
      Range: validators.Range{
        RangeEMsg: map[string]string{
          "between": "[name] 长度必须在 [min] 和 [max] 之间",
        },
      },
    },
  "datetime": &validators.DateTimeValidator{
    FmtStr: "Y-m-d",
  },
  "Email": &validators.EmailValidator{
    Reg: `^(\d)+$`,
  },
})
if err := validator.Validate(student); err != nil {
  fmt.Println(err)
}
```

MIT licence.
