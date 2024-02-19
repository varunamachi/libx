package data

import (
	"reflect"
	"strconv"
	"time"
)

// IsBasicType - tells if the kind of data type is basic or composite
func IsBasicType(rt reflect.Kind) bool {
	switch rt {
	case reflect.Bool:
		return true
	case reflect.Int:
		return true
	case reflect.Int8:
		return true
	case reflect.Int16:
		return true
	case reflect.Int32:
		return true
	case reflect.Int64:
		return true
	case reflect.Uint:
		return true
	case reflect.Uint8:
		return true
	case reflect.Uint16:
		return true
	case reflect.Uint32:
		return true
	case reflect.Uint64:
		return true
	case reflect.Uintptr:
		return true
	case reflect.Float32:
		return true
	case reflect.Float64:
		return true
	case reflect.Complex64:
		return true
	case reflect.Complex128:
		return true
	case reflect.Array:
		return false
	case reflect.Chan:
		return false
	case reflect.Func:
		return false
	case reflect.Interface:
		return false
	case reflect.Map:
		return false
	case reflect.Ptr:
		return false
	case reflect.Slice:
		return false
	case reflect.String:
		return true
	case reflect.Struct:
		return false
	case reflect.UnsafePointer:
		return false
	}
	return false
}

// IsTime - tells if a reflected value is time
func IsTime(val *reflect.Value) bool {
	return val.IsValid() &&
		val.Kind() == reflect.Struct &&
		val.Type() == reflect.TypeOf(time.Time{})
}

// ToFlatMap - converts given composite data structure into a map of string to
// interfaces. The heirarchy of types are flattened into single level. The
// keys of the map indicate the original heirarchy
func ToFlatMap(obj interface{}, tagName string) (out map[string]interface{}) {
	out = make(map[string]interface{})
	Walk(obj, &WalkConfig{
		MaxDepth:         InfiniteDepth,
		IgnoreContainers: false,
		VisitPrivate:     false,
		VisitRootStruct:  false,
		FieldNameRetriever: func(field *reflect.StructField) string {
			jt := field.Tag.Get(tagName)
			if jt != "" {
				return jt
			}
			return field.Name
		},
		Visitor: func(state *WalkerState) bool {
			if IsBasicType(state.Current.Kind()) || IsTime(state.Current) {
				out[state.Path] = state.Current.Interface()
			}
			return true
		},
	})
	return out
}

// VisitorFunc - function that will be called on each value of reflected type.
// The return value decides whether to continue with depth search in current
// branch
type VisitorFunc func(state *WalkerState) (cont bool)

// FieldNameRetriever - retrieves name for the field from given
type FieldNameRetriever func(field *reflect.StructField) (name string)

// WalkConfig - determines how Walk is carried out
type WalkConfig struct {
	Visitor            VisitorFunc        //visitor function
	FieldNameRetriever FieldNameRetriever //func to get name from struct field
	MaxDepth           int                //Stop walk at this depth
	IgnoreContainers   bool               //Ignore slice and map parent objects
	VisitPrivate       bool               //Visit private fields
	VisitRootStruct    bool               //Visit the root struct thats passed
}

// WalkerState - current state of the walk
type WalkerState struct {
	Depth   int
	Field   *reflect.StructField
	Path    string
	Parent  *reflect.Value
	Current *reflect.Value
}

// InfiniteDepth - used to indicate that Walk should continue till all the nodes
// in the heirarchy are visited
const InfiniteDepth int = -1

// Walk - walk a given instance of struct/slice/map/basic type
func Walk(
	obj interface{},
	config *WalkConfig) {
	// Wrap the original in a reflect.Value
	original := reflect.ValueOf(obj)
	if config.Visitor == nil {
		return
	}
	if config.FieldNameRetriever == nil {
		config.FieldNameRetriever = func(field *reflect.StructField) string {
			return field.Name
		}
	}
	walkRecursive(
		config,
		WalkerState{
			Depth:   0,
			Field:   nil,
			Path:    "",
			Parent:  nil,
			Current: &original,
		})
}

func walkRecursive(config *WalkConfig, state WalkerState) {
	if config.MaxDepth > 0 && state.Depth == config.MaxDepth+1 {
		return
	}
	//We copy any field from state which is used inside the loops, so that
	//state is not cumulatevily modified in a loop
	cur := state.Current
	path := state.Path
	switch state.Current.Kind() {
	case reflect.Ptr:
		originalValue := state.Current.Elem()
		if !originalValue.IsValid() {
			return
		}
		state.Parent = state.Current
		state.Current = &originalValue
		walkRecursive(config, state)

	case reflect.Interface:
		originalValue := state.Current.Elem()
		state.Parent = state.Current
		state.Current = &originalValue
		walkRecursive(config, state)

	case reflect.Struct:
		state.Depth++
		if state.Depth == 1 {
			if config.VisitRootStruct && !config.Visitor(&state) {
				return
			}
		} else if !config.Visitor(&state) {
			return
		}
		for i := 0; i < cur.NumField(); i++ {
			field := cur.Field(i)
			//Dont want to walk unexported fields if VisitPrivate is false
			if !(config.VisitPrivate || field.CanSet()) {
				continue
			}
			structField := cur.Type().Field(i)
			state.Field = &structField
			if path != "" {
				state.Path = path + "." +
					config.FieldNameRetriever(&structField)
			} else {
				state.Path = config.FieldNameRetriever(&structField)
			}
			state.Parent = state.Current
			state.Current = &field
			walkRecursive(config, state)
		}

	case reflect.Slice:
		state.Depth++
		if config.IgnoreContainers {
			return
		}

		for i := 0; i < cur.Len(); i++ {
			state.Field = nil
			state.Path = path + "." + strconv.Itoa(i)
			value := cur.Index(i)
			// state.Parent = state.Current
			state.Current = &value
			walkRecursive(config, state)
		}
	case reflect.Map:
		state.Depth++
		if config.IgnoreContainers {
			return
		}
		for _, key := range cur.MapKeys() {
			originalValue := cur.MapIndex(key)
			state.Field = nil
			state.Path = path + "." + key.String()
			state.Parent = state.Current
			state.Current = &originalValue
			walkRecursive(config, state)
		}
	// And everything else will simply be taken from the original
	default:
		if cont := config.Visitor(&state); !cont {
			return
		}

	}

}

func WalkPath(obj any, path string, out any) error {
	// TODO - implement
	return nil
}

func walkPath(index int, path []string, state *WalkerState) error {
	// TODO - implement
	return nil
}
