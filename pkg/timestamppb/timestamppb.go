package pgxprotobuftimestamppb

import (
	"fmt"

	"github.com/bcook98/pgx-protobuf-well-known/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Timestamppb timestamppb.Timestamp

// Timestamp

func (t *Timestamppb) ScanTimestamp(v pgtype.Timestamp) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into *timestamppb.Timestamp")
	}

	*t = Timestamppb(*timestamppb.New(v.Time))

	return nil
}

func (t Timestamppb) TimestampValue() (pgtype.Timestamp, error) {
	tt := timestamppb.Timestamp(t)
	return pgtype.Timestamp{Time: tt.AsTime(), Valid: true}, nil
}

func TryWrapTimestampEncodePlan(value interface{}) (plan pgtype.WrappedEncodePlanNextSetter, nextValue interface{}, ok bool) {
	switch value := value.(type) {
	case timestamppb.Timestamp:
		return &wrapTimestampEncodePlan{}, Timestamppb(value), true
	}

	return nil, nil, false
}

type wrapTimestampEncodePlan struct {
	next pgtype.EncodePlan
}

func (plan *wrapTimestampEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

func (plan *wrapTimestampEncodePlan) Encode(value interface{}, buf []byte) (newBuf []byte, err error) {
	return plan.next.Encode(Timestamppb(value.(timestamppb.Timestamp)), buf)
}

func TryWrapTimestampScanPlan(target interface{}) (plan pgtype.WrappedScanPlanNextSetter, nextDst interface{}, ok bool) {
	switch target := target.(type) {
	case *timestamppb.Timestamp:
		return &wrapTimestampScanPlan{}, (*Timestamppb)(target), true
	}

	return nil, nil, false
}

type wrapTimestampScanPlan struct {
	next pgtype.ScanPlan
}

func (plan *wrapTimestampScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

func (plan *wrapTimestampScanPlan) Scan(src []byte, dst interface{}) error {
	return plan.next.Scan(src, (*Timestamppb)(dst.(*timestamppb.Timestamp)))
}

type TimestampCodec struct {
	pgtype.TimestampCodec
}

func (TimestampCodec) DecodeValue(tm *pgtype.Map, oid uint32, format int16, src []byte) (interface{}, error) {
	if src == nil {
		return nil, nil
	}

	var target timestamppb.Timestamp
	scanPlan := tm.PlanScan(oid, format, &target)
	if scanPlan == nil {
		return nil, fmt.Errorf("PlanScan did not find a plan")
	}

	err := scanPlan.Scan(src, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// Timestamptz

func (t *Timestamppb) ScanTimestamptz(v pgtype.Timestamptz) error {
	if !v.Valid {
		return fmt.Errorf("cannot scan NULL into *timestamppb.Timestamp")
	}

	*t = Timestamppb(*timestamppb.New(v.Time))

	return nil
}

func (t Timestamppb) TimestamptzValue() (pgtype.Timestamptz, error) {
	tt := timestamppb.Timestamp(t)
	return pgtype.Timestamptz{Time: tt.AsTime(), Valid: true}, nil
}

func TryWrapTimestamptzEncodePlan(value interface{}) (plan pgtype.WrappedEncodePlanNextSetter, nextValue interface{}, ok bool) {
	switch value := value.(type) {
	case timestamppb.Timestamp:
		return &wrapTimestamptzEncodePlan{}, Timestamppb(value), true
	}

	return nil, nil, false
}

type wrapTimestamptzEncodePlan struct {
	next pgtype.EncodePlan
}

func (plan *wrapTimestamptzEncodePlan) SetNext(next pgtype.EncodePlan) { plan.next = next }

func (plan *wrapTimestamptzEncodePlan) Encode(value interface{}, buf []byte) (newBuf []byte, err error) {
	return plan.next.Encode(Timestamppb(value.(timestamppb.Timestamp)), buf)
}

func TryWrapTimestamptzScanPlan(target interface{}) (plan pgtype.WrappedScanPlanNextSetter, nextDst interface{}, ok bool) {
	switch target := target.(type) {
	case *timestamppb.Timestamp:
		return &wrapTimestamptzScanPlan{}, (*Timestamppb)(target), true
	}

	return nil, nil, false
}

type wrapTimestamptzScanPlan struct {
	next pgtype.ScanPlan
}

func (plan *wrapTimestamptzScanPlan) SetNext(next pgtype.ScanPlan) { plan.next = next }

func (plan *wrapTimestamptzScanPlan) Scan(src []byte, dst interface{}) error {
	return plan.next.Scan(src, (*Timestamppb)(dst.(*timestamppb.Timestamp)))
}

type TimestamptzCodec struct {
	pgtype.TimestamptzCodec
}

func (TimestamptzCodec) DecodeValue(tm *pgtype.Map, oid uint32, format int16, src []byte) (interface{}, error) {
	if src == nil {
		return nil, nil
	}

	var target timestamppb.Timestamp
	scanPlan := tm.PlanScan(oid, format, &target)
	if scanPlan == nil {
		return nil, fmt.Errorf("PlanScan did not find a plan")
	}

	err := scanPlan.Scan(src, &target)
	if err != nil {
		return nil, err
	}

	return target, nil
}

// Register registers the protobuf/types/known/timestamppb integration with a pgtype.ConnInfo.
func RegisterTimestamp(m *pgtype.Map) {
	m.TryWrapEncodePlanFuncs = append([]pgtype.TryWrapEncodePlanFunc{TryWrapTimestampEncodePlan}, m.TryWrapEncodePlanFuncs...)
	m.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapTimestampScanPlan}, m.TryWrapScanPlanFuncs...)

	m.RegisterType(&pgtype.Type{
		Name:  "timestamp",
		OID:   pgtype.TimestampOID,
		Codec: TimestampCodec{},
	})

	util.RegisterDefaultPgTypeVariants(m, "timestamp", "_timestamp", timestamppb.Timestamp{})
}

// Register registers the protobuf/types/known/timestamppb integration with a pgtype.ConnInfo.
func RegisterTimestamptz(m *pgtype.Map) {
	m.TryWrapEncodePlanFuncs = append([]pgtype.TryWrapEncodePlanFunc{TryWrapTimestamptzEncodePlan}, m.TryWrapEncodePlanFuncs...)
	m.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapTimestamptzScanPlan}, m.TryWrapScanPlanFuncs...)

	m.RegisterType(&pgtype.Type{
		Name:  "timestamptz",
		OID:   pgtype.TimestamptzOID,
		Codec: TimestamptzCodec{},
	})

	util.RegisterDefaultPgTypeVariants(m, "timestamptz", "_timestamptz", timestamppb.Timestamp{})
}
