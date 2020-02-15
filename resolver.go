package gqlgen_casbin_RBAC_example

import (
	"context"
	"log"
	"reflect"
) // THIS CODE IS A STARTING POINT ONLY. IT WILL NOT BE UPDATED WITH SCHEMA CHANGES.

type Resolver struct{
	todos []*Todo
}

func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

func (r *Resolver) Todo() TodoResolver {
	return &todoResolver{r}
}

type mutationResolver struct{ *Resolver }

func accessControl(input interface{},sub,domain,table,act string) (res map[string]interface{})   {
	res = make(map[string]interface{})
	v:=reflect.TypeOf(input)
	rv:=reflect.ValueOf(input)
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		value:=rv.Field(i).Interface()
		ok,err:=enforcer.Enforce(sub,domain,table,field.Name,act)
		if err != nil {
			// handle err
		}
		if ok==true{
			res[field.Name]=value
		} else {

		}
	}
	return
}
func createStructByReflect(data map[string]interface{}, inStructPtr interface{}) {
	rType := reflect.TypeOf(inStructPtr)
	rVal := reflect.ValueOf(inStructPtr)
	if rType.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		rType = rType.Elem()
		rVal = rVal.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}
	// 遍历结构体
	for i := 0; i < rType.NumField(); i++ {
		t := rType.Field(i)
		f := rVal.Field(i)
		// 得到tag中的字段名
		//key := t.Tag.Get("key")
		key:=t.Name
		if v, ok := data[key]; ok {
			// 检查是否需要类型转换
			dataType := reflect.TypeOf(v)
			structType := f.Type()
			if structType == dataType {
				f.Set(reflect.ValueOf(v))
			} else {
				if dataType.ConvertibleTo(structType) {
					// 转换类型
					f.Set(reflect.ValueOf(v).Convert(structType))
				} else {
					panic(t.Name + " type mismatch")
				}
			}
		} else {
			log.Print(t.Name + " not found")
		}
	}
}
func (r *mutationResolver) CreateTodo(ctx context.Context, input NewTodo) (*Todo, error) {
	sub:="user3"
	domain:="domain1"
	table:="todo"
	act:="write"
	//reflect request input
	res:=accessControl(input,sub,domain,table,act)
	//structByReflect
	t:=&Todo{}
	createStructByReflect(res,t)

	r.todos = append(r.todos, t)
	return t, nil
}

type queryResolver struct{ *Resolver }

func (r *queryResolver) Todos(ctx context.Context) ([]*Todo, error) {
	return r.todos, nil
}

type todoResolver struct{ *Resolver }

func (r *todoResolver) User(ctx context.Context, obj *Todo) (*User, error) {
	return &User{ID: obj.UserID, Name: "user " + obj.UserID}, nil
}
