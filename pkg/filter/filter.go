package filter

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Matcher interface {
	Match(string) bool
}

func iterateArray(val reflect.Value, visited map[uintptr]bool, m string, glob Matcher) (any, error) {
	arr := make([]any, 0)
	if val.Kind() == reflect.Array || val.Kind() == reflect.Slice {
		for i := range val.Len() {
			elem := val.Index(i)

			if !elem.CanInterface() {
				continue
			}
			n := m
			if n == "" {
				n = fmt.Sprint(i)
			} else {
				n = fmt.Sprintf("%v.%v", n, i)
			}
			// fmt.Printf("%d :: \n", i)
			if val, err := FilterType(elem.Interface(), visited, n, glob); err != nil {
				return val, err
			} else {
				if val != nil {
					arr = append(arr, val)
				}
			}
		}
		if len(arr) > 0 {
			return arr, nil
		}
	}
	return nil, nil
}

func iterateStruct(val reflect.Value, visited map[uintptr]bool, m string, glob Matcher) (any, error) {
	in := make(map[string]any)

	typ := val.Type()
	for i := range val.NumField() {
		field := typ.Field(i)
		fieldValue := val.Field(i)

		if !fieldValue.CanInterface() {
			continue
		}

		n := m
		if n == "" {
			n = field.Name
		} else {
			n = fmt.Sprintf("%v.%v", n, field.Name)
		}
		if val, err := FilterType(fieldValue.Interface(), visited, n, glob); err != nil {
			return val, err
		} else {
			if val != nil {
				in[field.Name] = val
			}
		}

	}
	return in, nil
}

func canBeIgnored(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Chan:
	case reflect.Func:
	case reflect.Invalid:
	case reflect.Interface:
		return true
	}
	return false
}

func FilterType(v any, visited map[uintptr]bool, m string, glob Matcher) (any, error) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return nil, nil
		}
		ptr := val.Pointer()
		if visited[ptr] {
			return nil, errors.New("cyclic reference detected")
		}
		visited[ptr] = true
		val = val.Elem()
	}

	if canBeIgnored(val) {
		return nil, nil
	}

	switch val.Kind() {
	case reflect.Array:
	case reflect.Slice:
		return iterateArray(val, visited, m, glob)
	case reflect.Struct:
		return iterateStruct(val, visited, m, glob)
	case reflect.Map:
		panic("not implemented")
	default:
		// fmt.Printf(":: %v\n", val.Interface())
		if glob.Match(m) {
			return val.Interface(), nil
		} else {
			return nil, nil
		}
	}

	return nil, nil
}

// func mainx() {
// 	// 	type Inner struct {
// 	// 		C int
// 	// 	}
// 	// 	type Outer struct {
// 	// 		A int
// 	// 		B *Inner
// 	// 	}

// 	// 	inner := &Inner{C: 10}
// 	// 	outer := &Outer{A: 5, B: inner}
// 	// 	innerRefCycle := &Outer{A: 8, B: inner}

// 	// 	visited := make(map[uintptr]bool)
// 	// 	if err := iterateStruct(outer, visited); err != nil {
// 	// 		fmt.Println("Error:", err)
// 	// 	}

// 	// 	visited = make(map[uintptr]bool)
// 	// 	if err := iterateStruct(innerRefCycle, visited); err != nil {
// 	// 		fmt.Println("Error:", err)
// 	// 	}
// 	type Post struct {
// 		Title string
// 		Test  string
// 	}

// 	type User struct {
// 		Name string
// 		Post []Post
// 	}

// 	post := Post{"hello", "test"}
// 	user := User{
// 		Post: []Post{post, post},
// 		Name: "hero",
// 	}

// 	visited := make(map[uintptr]bool)
// 	g := glob.MustCompile(CreateGlobPattern([]string{"*", "!Name", "Post.*.Title"}), '.')
// 	if val, err := FilterType(user, visited, "", g); err != nil {
// 		fmt.Println("Error:", err, val)
// 	} else {
// 		fmt.Println(val)
// 		b, err := json.MarshalIndent(val, "", "  ")
// 		if err != nil {
// 			fmt.Println("error:", err)
// 		}
// 		fmt.Println(string(b))

// 	}

// 	/**
// 	*, !Name, !Post.Test
// 	 */

// 	// fmt.Println(globs([]string{"*", "!Name", "!Post.Test", "Post.1"}))
// 	fmt.Println(glob.MustCompile(CreateGlobPattern([]string{"*", "!Name", "!Post.2.*", "Post.*.1"}), '.').Match("Post.2.2"))

// }

func CreateGlobPattern(patterns []string) string {
	return fmt.Sprintf("{%v}", strings.Join(patterns, ","))
}
