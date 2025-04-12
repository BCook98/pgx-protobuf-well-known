package pgxprotobuftimestamppb

import (
	"fmt"

	"github.com/bcook98/pgx-protobuf-well-known/pkg/util"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Timestamppb timestamppb.Timestamp

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

// Register registers the protobuf/types/known/timestamppb integration with a pgtype.ConnInfo.
func Register(m *pgtype.Map) {
	m.TryWrapEncodePlanFuncs = append([]pgtype.TryWrapEncodePlanFunc{TryWrapTimestampEncodePlan}, m.TryWrapEncodePlanFuncs...)
	m.TryWrapScanPlanFuncs = append([]pgtype.TryWrapScanPlanFunc{TryWrapTimestampScanPlan}, m.TryWrapScanPlanFuncs...)

	util.RegisterDefaultPgTypeVariants(m, "timestamp", "_timestamp", timestamppb.Timestamp{})
	util.RegisterDefaultPgTypeVariants(m, "timestamptz", "_timestamptz", timestamppb.Timestamp{})
}
