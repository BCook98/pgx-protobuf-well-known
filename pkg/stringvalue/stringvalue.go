package pgxprotobufstringvalue

import (
	"fmt"

	"github.com/bcook98/pgx-protobuf-well-known/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type StringValue wrapperspb.StringValue

func (s *StringValue) ScanText(v pgtype.Text) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into *wrapperspb.StringValue")
	}

	*s = StringValue(wrapperspb.StringValue{Value: v.String})

	return nil
}

func (s StringValue) TextValue() (pgtype.Text, error) {
	return pgtype.Text{String: s.Value, Valid: true}, nil
}

func TrywrapEncodePlan(value interface{}) (plan pgtype.WrappedEncodePlanNextSetter, nextValue interface{}, ok bool) {
	switch value := value.(type) {
	case wrapperspb.StringValue:
		return &wrapEncodePlan{}, StringValue(value), true
	}

	return nil, nil, false
}

type wrapEncodePlan struct {
	next pgtype.EncodePlan
}

func (plan *wrapEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

func (plan *wrapEncodePlan) Encode(value interface{}, buf []byte) (newBuf []byte, err error) {
	return plan.next.Encode(StringValue(value.(wrapperspb.StringValue)), buf)
}

func TryWrapScanPlan(target interface{}) (plan pgtype.WrappedScanPlanNextSetter, nextDst interface{}, ok bool) {
	switch target := target.(type) {
	case *wrapperspb.StringValue:
		return &wrapScanPlan{}, (*StringValue)(target), true
	}

	return nil, nil, false
}

type wrapScanPlan struct {
	next pgtype.ScanPlan
}

func (plan *wrapScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

func (plan *wrapScanPlan) Scan(src []byte, dst interface{}) error {
	return plan.next.Scan(src, (*StringValue)(dst.(*wrapperspb.StringValue)))
}

func Register(m *pgtype.Map) {
	m.TryWrapEncodePlanFuncs = append([]pgtype.TryWrapEncodePlanFunc{TrywrapEncodePlan}, m.TryWrapEncodePlanFuncs...)
	m.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapScanPlan}, m.TryWrapScanPlanFuncs...)

	util.RegisterDefaultPgTypeVariants(m, "text", "_text", wrapperspb.StringValue{})
}
