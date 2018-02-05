package core

import (
    "fmt"
    "errors"
    "reflect"
    "strings"
)

const (
    PrefixSet       = "Set"
    PrefixOnUpdate  = "OnUpdate"
)

var (
    ResUpdateIgnoreKey      = errors.New("ignore this key")
    ResUpdateInvSetReturn   = errors.New("invalid set return")
)

var (
    processors  map[string]*confProcessor
    errorInterface = reflect.TypeOf((*error)(nil)).Elem()
)

func init() {
    processors = make(map[string]*confProcessor)
}

func RegisterEntry(tag string, goName string, v reflect.Value) error {
    _, ok := processors[tag]
    if ok {
        return fmt.Errorf("duplicated pconf tag: %v", tag)
    }
    processors[tag] = &confProcessor{
        NameKey:    tag,
        Name:       goName,
    }
    return processors[tag].BuildFunc(v)
}

//获取所有配置项的key
func EntryKeys() []string {
    names := make([]string, len(processors))
    cnt := 0
    for _, p := range processors {
        names[cnt] = p.NameKey
        cnt++
    }
    return names
}

//初始化特定配置值
func InitEntry(key, val string) error {
    processor, ok := processors[key]
    if !ok {
        return ResUpdateIgnoreKey
    }

    //call setter
    in := []reflect.Value{reflect.ValueOf(val)}
    returns := processor.FuncSet.Call(in)
    if len(returns) != 1 {
        return ResUpdateInvSetReturn
    }
    err := returns[0].Interface()
    if err != nil {
        return fmt.Errorf("set conf failed: key %v, val %v, err %v", key, val, err)
    }
    Log("key %v init local ok.", key)
    return nil
}

//更新特定配置值
func UpdateEntry(key, oldVal, val string) error {
    processor, ok := processors[key]
    if !ok {return ResUpdateIgnoreKey}

    //call setter
    in := []reflect.Value{reflect.ValueOf(val)}
    returns := processor.FuncSet.Call(in)
    if len(returns) != 1 {
        return ResUpdateInvSetReturn
    }
    err := returns[0].Interface()
    if err != nil {
        return fmt.Errorf("set conf failed: key %v, val %v, err %v", key, val, err)
    }
    Log("key %v update set local value %v ok.", key, val)

    //call on updater
    in = []reflect.Value{reflect.ValueOf(oldVal), reflect.ValueOf(val)}
    processor.FuncOnUpdate.Call(in)

    Log("key %v on-update local ok.", key)
    return nil
}

type confProcessor struct {
    NameKey         string          //小写下划线配置名字
    Name            string
    FuncGet         reflect.Value   //获取函数
    FuncSet         reflect.Value   //设置函数
    FuncOnUpdate    reflect.Value   //更新回调
}

func (cp *confProcessor) BuildFunc(v reflect.Value) error {
    cp.FuncGet = v.MethodByName(cp.Name)
    if cp.FuncGet.IsValid() == false {
        return fmt.Errorf("error conf<key %v, name %v>: get method not found!", cp.NameKey, cp.Name)
    } else {
        if cp.FuncGet.Type().NumIn() > 0 {
            return fmt.Errorf("error conf<key %v, name %v>: get method must NOT have parameters", cp.NameKey, cp.Name)
        }
        if cp.FuncGet.Type().NumOut() == 0 {
            return fmt.Errorf("error conf<key %v, name %v>: get method without output parameters!", cp.NameKey, cp.Name)
        }
    }

    cp.FuncSet = v.MethodByName(PrefixSet + cp.Name)
    if cp.FuncSet.IsValid() == false {
        return fmt.Errorf("error conf<key %v, name %v>: set method not found!", cp.NameKey, cp.Name)
    } else {
        if cp.FuncSet.Type().NumIn() != 1 {
            return fmt.Errorf("error conf<key %v, name %v>: Set method with input parameter count not equal to 1",
                cp.NameKey, cp.Name)
        } else {
            param := cp.FuncSet.Type().In(0)
            if param.Kind() != reflect.String {
                return fmt.Errorf("error conf<key %v, name %v>: Set method parameter must be string",
                    cp.NameKey, cp.Name)
            }
        }
        if cp.FuncSet.Type().NumOut() != 1 {
            return fmt.Errorf("error conf<key %v, name %v>: set method must has 1 erro return value!",
                cp.NameKey, cp.Name)
        } else {
            param := cp.FuncSet.Type().Out(0)
            if param.Kind() != reflect.Interface || param.Implements(errorInterface) == false {
                return fmt.Errorf("error conf<key %v, name %v>: Set method return must be type error, kind %v",
                    cp.NameKey, cp.Name, param.Kind())
            }
        }
    }

    cp.FuncOnUpdate = v.MethodByName(PrefixOnUpdate + cp.Name)
    if cp.FuncOnUpdate.IsValid() == false {
        return fmt.Errorf("error conf<key %v, name %v>: on-update method not found!", cp.Name, cp.Name)
    } else {
        if cp.FuncOnUpdate.Type().NumIn() != 2 {
            return fmt.Errorf("error conf<key %v, name %v>: on-update method with input parameter count not equal to 2",
                cp.Name, cp.Name)
        } else {
            param := cp.FuncOnUpdate.Type().In(0)
            if param.Kind() != reflect.String {
                return fmt.Errorf("error conf<key %v, name %v>: on-update method parameter 1 must be string",
                    cp.Name, cp.Name)
            }
            param = cp.FuncOnUpdate.Type().In(1)
            if param.Kind() != reflect.String {
                return fmt.Errorf("error conf<key %v, name %v>: on-update method parameter 2 must be string",
                    cp.Name, cp.Name)
            }
        }
        if cp.FuncOnUpdate.Type().NumOut() > 0 {
            return fmt.Errorf("error conf<key %v, name %v>: on-update method with return value.", cp.Name, cp.Name)
        }
    }
    return nil
}

func gonicCasedName(name string) string {
    newstr := make([]rune, 0, len(name)+3)
    for idx, chr := range name {
        if isASCIIUpper(chr) && idx > 0 {
            if !isASCIIUpper(newstr[len(newstr)-1]) {
                newstr = append(newstr, '_')
            }
        }

        if !isASCIIUpper(chr) && idx > 1 {
            l := len(newstr)
            if isASCIIUpper(newstr[l-1]) && isASCIIUpper(newstr[l-2]) {
                newstr = append(newstr, newstr[l-1])
                newstr[l-1] = '_'
            }
        }

        newstr = append(newstr, chr)
    }
    return strings.ToLower(string(newstr))
}

func isASCIIUpper(r rune) bool {
    return 'A' <= r && r <= 'Z'
}
