package data

import (
	"errors"
	"reflect"
	"regexp"
	"strconv"
	"time"

	"github.com/varunamachi/libx/errx"
	"github.com/varunamachi/libx/str"
)

var (
	ErrInvalidWalkPath = errors.New("invalid walk path")
	ErrInvalidData     = errors.New("invalid data")
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
		if !config.Visitor(&state) {
			return
		}

	}

}

func Get(obj any, out any, path ...string) error {
	val, err := WalkPath(obj, path...)
	if err != nil {
		return errx.Wrap(err)
	}

	if !val.IsValid() {
		return errx.Errf(
			ErrInvalidWalkPath, "invalid value found at given path")
	}

	if val.Type().AssignableTo(reflect.TypeOf(out)) {
		reflect.ValueOf(out).Set(val)
		return nil
	}

	// TODO - what happens here??
	return nil
}

func Set(obj any, newVal any, path ...string) error {
	val, err := WalkPath(obj, path...)
	if err != nil {
		return errx.Wrap(err)
	}

	if !val.IsValid() {
		return errx.Errf(
			ErrInvalidWalkPath, "invalid value found at given path")
	}

	if reflect.TypeOf(newVal).AssignableTo(val.Type()) {
		val.Set(reflect.ValueOf(newVal))
		return nil
	}

	// TODO - handle cases of struct to map??

	return nil
}

// TODO - handle properly
var indexRegEx, _ = regexp.Compile(`^#\-?\d+`)

func WalkPath(obj any, path ...string) (reflect.Value, error) {
	var cur reflect.Value

	for pi := 0; pi < len(path); {
		pathElem := path[pi]
		cur = reflect.ValueOf(obj)

		switch cur.Kind() {
		case reflect.Ptr:
			cur = cur.Elem()
			if !cur.IsValid() {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath, "invalid pointer value")
			}
			continue

		case reflect.Interface:
			cur = cur.Elem()
			continue
		case reflect.Struct:
			if indexRegEx.MatchString(pathElem) {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath, "expected slice got struct")
			}

			for i := 0; i < cur.NumField(); i++ {
				field := cur.Field(i)
				//Not interested in unexported fields
				if !field.CanSet() {
					continue
				}
				structField := cur.Type().Field(i)

				if str.EqFold(structField.Name, pathElem) {
					cur = field
				}
			}

			return reflect.ValueOf(nil),
				errx.Errf(ErrInvalidWalkPath,
					"could not find path element '%s'", pathElem)

		case reflect.Slice:
			if !indexRegEx.MatchString(pathElem) {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath,
						"expected key or a name, go slice")
			}

			index, err := strconv.Atoi(pathElem[1:])
			if err != nil {
				return reflect.ValueOf(nil),
					errx.Errf(err,
						"failed to get index from path element: '%s'", pathElem)
			}

			index = Qop(index < 0, cur.Len()+index, index)
			if index < 0 || index >= cur.Len() {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath,
						"invalid index '%d', slice size: '%d'",
						index, cur.Len())
			}
			cur = cur.Index(index)
		case reflect.Map:
			if indexRegEx.MatchString(pathElem) {
				return reflect.ValueOf(nil),
					errx.Errf(
						ErrInvalidWalkPath, "expected slice got struct")
			}
			cur = cur.MapIndex(reflect.ValueOf(pathElem))
			if cur.IsZero() {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath,
						"map does not have key '%s'", pathElem)
			}
		default:
			if !IsBasicType(cur.Kind()) {
				return reflect.ValueOf(nil), errx.Errf(
					ErrInvalidData,
					"chans, funcs, unsafe-ptrs are not supported, found: '%s'",
					cur.Kind().String())
			}
			if pi != len(path)-1 {
				return reflect.ValueOf(nil),
					errx.Errf(ErrInvalidWalkPath,
						"found a primitive in the middle of the path, "+
							"cannot proceed. Index: %d, PathElement: %s",
						pi, pathElem)
			}
			// Here it should automatically brake, since pi == len(path) in
			// next iteration
		}
	}
	return cur, nil
}
