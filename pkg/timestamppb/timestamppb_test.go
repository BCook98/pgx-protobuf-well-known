package pgxprotobuftimestamppb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	testutil "github.com/bcook98/pgx-protobuf-well-known/pkg/test_util"
	pgxprotobuftimestamppb "github.com/bcook98/pgx-protobuf-well-known/pkg/timestamppb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var defaultConnTestRunner pgxtest.ConnTestRunner

func init() {
	defaultConnTestRunner = testutil.ConnTestRunner()
	defaultConnTestRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxprotobuftimestamppb.Register(conn.TypeMap())
	}
}

var postgresTypes = []string{
	"timestamp",
	"timestamptz",
}

func TestNull(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
				var ts timestamppb.Timestamp
				err := conn.QueryRow(context.Background(), fmt.Sprintf(`select null::%s`, typ)).Scan(&ts)
				require.EqualError(t, err, `can't scan into dest[0]: cannot scan NULL into *timestamppb.Timestamp`)
			})
		})
	}
}

func TestArray(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
				inputSlice := []*timestamppb.Timestamp{}

				now := time.Now()
				for i := 0; i < 10; i++ {
					inputSlice = append(inputSlice, timestamppb.New(now.Add(time.Duration(i)*time.Hour*24)))
				}

				var outputSlice []*timestamppb.Timestamp
				err := conn.QueryRow(context.Background(), fmt.Sprintf(`select $1::%s[]`, typ), inputSlice).Scan(&outputSlice)
				require.NoError(t, err)

				require.Equal(t, len(inputSlice), len(outputSlice))

				for i := 0; i < len(inputSlice); i++ {
					inputSlice[i].AsTime().Equal(outputSlice[i].AsTime())
				}
			})
		})
	}
}

func isExpectedEqTimestamp(a *timestamppb.Timestamp) func(interface{}) bool {
	return func(v interface{}) bool {
		return a.AsTime().Equal((v.(*timestamppb.Timestamp)).AsTime())
	}
}

func TestValueRoundTrip(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			pgxtest.RunValueRoundTripTests(context.Background(), t, defaultConnTestRunner, nil, typ, []pgxtest.ValueRoundTripTest{
				{
					Param:  timestamppb.New(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1000, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1600, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1700, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1700, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1800, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(1999, 12, 31, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2001, 1, 2, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2001, 1, 2, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2004, 2, 29, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2013, 7, 4, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2013, 7, 4, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2013, 12, 25, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2013, 12, 25, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2029, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2081, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2081, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2096, 2, 29, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2096, 2, 29, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(2550, 1, 1, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(2550, 1, 1, 0, 0, 0, 0, time.UTC))),
				},
				{
					Param:  timestamppb.New(time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC)),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(time.Date(9999, 12, 31, 0, 0, 0, 0, time.UTC))),
				},
			})
		})
	}
}
