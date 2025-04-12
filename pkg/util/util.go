package util

import (
	"reflect"

	"github.com/jackc/pgx/v5/pgtype"
)

func RegisterDefaultPgTypeVariants(m *pgtype.Map, name, arrayName string, value interface{}) {
	// T
	m.RegisterDefaultPgType(value, name)

	// *T
	valueType := reflect.TypeOf(value)
	m.RegisterDefaultPgType(reflect.New(valueType).Interface(), name)

	// []T
	sliceType := reflect.SliceOf(valueType)
	m.RegisterDefaultPgType(reflect.MakeSlice(sliceType, 0, 0).Interface(), arrayName)

	// *[]T
	m.RegisterDefaultPgType(reflect.New(sliceType).Interface(), arrayName)

	// []*T
	sliceOfPointerType := reflect.SliceOf(reflect.TypeOf(reflect.New(valueType).Interface()))
	m.RegisterDefaultPgType(reflect.MakeSlice(sliceOfPointerType, 0, 0).Interface(), arrayName)

	// *[]*T
	m.RegisterDefaultPgType(reflect.New(sliceOfPointerType).Interface(), arrayName)
}
