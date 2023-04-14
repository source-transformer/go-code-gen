// Copyright 2021-present Drop Fake Inc. All rights reserved.

package csharp

import (
	"fmt"
	"reflect"
	"strconv"
)

const (
	verbose = false
)

/*
// types to skip/ignore
var (
	dataSchemeTypeBaseType = reflect.TypeOf(definition.DataSchemaTypeBase{})
)
*/

func resolveReflectType(t reflect.Type) (reflect.Type, error) {
	// don't need to handle reflect.Map: our source data is not allowed to use a Go Map because it is non deterministically ordered
	switch t.Kind() {
	case reflect.Ptr:
		return resolveReflectType(t.Elem())
	case reflect.Slice: // not sure if we need to support Array (, reflect.Array) as well - need a test if we do
		return resolveReflectType(t.Elem())
	case reflect.Interface:
		return t, nil
	case reflect.Struct:
		return t, nil
	default:
		if verbose {
			fmt.Println("resolveReflectType: kind:", t.Kind(), "type:", t, "value:", t.String())
		}
		return t, nil
	}
}

func convertGoToCSharpType(t reflect.Type) (string, error) {
	if verbose {
		fmt.Println("convertGoToCSharpType: kind:", t.Kind(), "type:", t, "value:", t.String())
	}
	switch t.Kind() {
	case reflect.String:
		return "string", nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return "int", nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "uint", nil
	case reflect.Float32, reflect.Float64:
		return "float", nil
	case reflect.Bool:
		return "bool", nil
	case reflect.Interface:
		return "object", nil
	case reflect.Struct:
		return t.Name(), nil
	default:
		return "", fmt.Errorf("unsupported type: %v", t)
	}
}

func typeAlreadyDefined(reflectType reflect.Type, definedTypesPtr *[]reflect.Type) bool {
	definedTypes := *definedTypesPtr
	fmt.Println("typeAlreadyDefined: len(definedTypes):", len(definedTypes))
	for _, definedType := range definedTypes {
		fmt.Println("typeAlreadyDefined: checking for type:", reflectType.Name(), "against:", definedType.Name())
		if reflectType == definedType {
			return true
		}
	}
	return false
}

func codegenForStruct(reflectValue reflect.Value, result string, definedTypes *[]reflect.Type) (string, error) {
	if verbose {
		fmt.Println("Struct: reflectValue.Type:", reflectValue.Type(), "reflectValue.String():", reflectValue.String())
	}

	reflectType := reflectValue.Type()
	typeAlreadyDefined := typeAlreadyDefined(reflectType, definedTypes)
	if typeAlreadyDefined {
		fmt.Println("traverseReflectValue: type:", reflectType.Name(), "already defined")
		return result, nil
	}
	classDefine := fmt.Sprintf("\npublic class %v\n{\n", reflectType.Name())
	*definedTypes = append(*definedTypes, reflectType)
	fmt.Println("added:", reflectType.Name(), "to definedTypes (", len(*definedTypes), ")")
	// Iterate over struct fields and traverse them recursively
	for i := 0; i < reflectValue.NumField(); i++ {
		fieldValue := reflectValue.Field(i)
		fieldType := reflectType.Field(i)
		fieldPrefix := fieldType.Name

		if verbose {
			fmt.Println("fieldValue.Kind():", fieldValue.Kind(), "CanAddr:", fieldValue.CanAddr(), "fieldValue.Type():", fieldValue.Type(), "fieldValue.String():", fieldValue.String())
			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.CanAddr() {
					fmt.Println("fieldValue.Addr().Type():", fieldValue.Addr().Type(), "fieldValue.Addr().Kind():", fieldValue.Addr().Kind(), "fieldValue.Addr().String():", fieldValue.Addr().String())
				}
			}
		}

		resolvedType, err := resolveReflectType(fieldValue.Type())
		if verbose {
			fmt.Println("Struct: field(2):", i, "resolvedType.Name():", resolvedType.Name(), "resolvedType:", resolvedType, "err:", err)
		}
		/*if resolvedType == dataSchemeTypeBaseType {
			fmt.Println("skipping DataSchemaTypeBase field: ", fieldPrefix)
			continue
		}
		*/

		// Traverse nested field recursively
		result, err = traverseReflectValue(fieldPrefix, fieldValue, result, definedTypes)
		if err != nil {
			return result, err
		}

		arraySubscript := ""
		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			arraySubscript = "[]"
		}
		if verbose {
			fmt.Println("reflectValue.Type():", reflectValue.Type(), "fieldType:", fieldType, "fieldValue.Type():", fieldValue.Type().Name())
		}

		csharpTypeStr, err := convertGoToCSharpType(resolvedType)
		if err != nil {
			return result, err
		}
		fieldStr := fmt.Sprintf("\tpublic %v%v %v { get; set; }", csharpTypeStr, arraySubscript, fieldPrefix)
		classDefine += fieldStr
		classDefine += "\n"
	}

	classDefine += "}\n"
	result += classDefine

	return result, nil
}

func traverseReflectValue(prefix string, reflectValue reflect.Value, generatedCode string, definedTypes *[]reflect.Type) (string, error) {
	if verbose {
		fmt.Println("traverseReflectValue: Type():", reflectValue.Type(), "kind:", reflectValue.Kind(), "value:", reflectValue.String(), "isValid:", strconv.FormatBool((reflectValue.IsValid())))
	}
	if !reflectValue.IsValid() {
		return generatedCode, nil
	}
	// no handling for reflect.Array or reflect.Map yet
	switch reflectValue.Kind() {
	case reflect.Ptr:
		// Handle nil values
		if reflectValue.IsNil() {
			fmt.Printf("%s nil\n", prefix)
			return generatedCode, nil
		}
		if verbose {
			fmt.Println("reflectValue.IsNil():", reflectValue.IsNil(), "String:", reflectValue.String(), "Type:", reflectValue.Type(), "Kind:", reflectValue.Kind())
			fmt.Println("Addr().Type():", reflectValue.Addr().Type(), "Addr().Kind():", reflectValue.Addr().Kind(), "reflectValue.Addr():", reflectValue.Addr().String())
		}
		result, err := traverseReflectValue(prefix, reflectValue.Elem(), generatedCode, definedTypes)
		return result, err
	case reflect.Slice:
		// can't call reflectValue.Elem() on Slice - will panic with: "reflect: call of reflect.Value.Elem on slice Value"
		sliceElemType := reflectValue.Type().Elem()
		sliceType := reflect.TypeOf(reflectValue)
		if sliceElemType.Kind() == reflect.Ptr {
			sliceElemType = sliceElemType.Elem()
		}
		resolvedType, err := resolveReflectType(reflectValue.Type())
		if verbose {
			fmt.Println("Slice: err:", err, "sliceType:", sliceType, "sliceElemType:", "resolvedType.Name():", resolvedType.Name(), "resolvedType:", resolvedType)
		}
		if err != nil {
			return generatedCode, err
		}

		reflectValue2 := reflect.New(sliceElemType)
		generatedCode, err = traverseReflectValue(prefix, reflectValue2, generatedCode, definedTypes)
		return generatedCode, err
	case reflect.Interface:
		result, err := traverseReflectValue(prefix, reflectValue.Elem(), generatedCode, definedTypes)
		return result, err
	case reflect.Struct:
		result, err := codegenForStruct(reflectValue, generatedCode, definedTypes)
		return result, err
	default:
		if verbose {
			fmt.Println("default: 1:", reflectValue.Type(), "1.5:", reflectValue.String(), "2:", reflectValue.Addr().Type(), "3:", reflectValue.Addr().Kind(), "4:", reflectValue.Addr().String())
		}
	}
	return generatedCode, nil
}

func GenerateCSharpFromInstance(obj interface{}) (string, error) {
	reflectValue := reflect.ValueOf(obj)
	definedTypesPtr := &[]reflect.Type{}
	result, err := traverseReflectValue("", reflectValue, "", definedTypesPtr)
	if verbose {
		fmt.Println("GenerateCSharpFromInstance: err:", err, "result:", result)
	}

	return result, err
}
