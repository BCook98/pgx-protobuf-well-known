package pgxprotobufstringvalue_test

import (
	"context"
	"fmt"
	"testing"

	pgxprotobufstringvalue "github.com/bcook98/pgx-protobuf-well-known/pkg/stringvalue"
	testutil "github.com/bcook98/pgx-protobuf-well-known/pkg/test_util"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var defaultConnTestRunner pgxtest.ConnTestRunner

func init() {
	defaultConnTestRunner = testutil.ConnTestRunner()
	defaultConnTestRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxprotobufstringvalue.Register(conn.TypeMap())
	}
}

var postgresTypes = []string{
	"text",
	"varchar",
}

func TestNull(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
				var s wrapperspb.StringValue
				err := conn.QueryRow(context.Background(), fmt.Sprintf(`select null::%s`, typ)).Scan(&s)
				require.EqualError(t, err, `can't scan into dest[0]: cannot scan NULL into *wrapperspb.StringValue`)

				var sPtrIn *wrapperspb.StringValue
				var sPtrOut *wrapperspb.StringValue
				err = conn.QueryRow(context.Background(), fmt.Sprintf(`select $1::%s`, typ), sPtrIn).Scan(&sPtrOut)
				require.NoError(t, err)
				require.Nil(t, sPtrOut)
			})
		})
	}
}

func TestArray(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
				inputSlice := []*wrapperspb.StringValue{}

				for i := 0; i < 10; i++ {
					inputSlice = append(inputSlice, wrapperspb.String(fmt.Sprintf("test %d", i)))
				}

				var outputSlice []*wrapperspb.StringValue
				err := conn.QueryRow(context.Background(), fmt.Sprintf(`select $1::%s[]`, typ), inputSlice).Scan(&outputSlice)
				require.NoError(t, err)

				require.Equal(t, len(inputSlice), len(outputSlice))

				for i := 0; i < len(inputSlice); i++ {
					require.Equal(t, inputSlice[i].GetValue(), outputSlice[i].GetValue())
				}
			})
		})
	}
}

func isExpectedEqStringValue(a *wrapperspb.StringValue) func(interface{}) bool {
	return func(v interface{}) bool {
		return a.GetValue() == (v.(*wrapperspb.StringValue)).GetValue()
	}
}

func TestValueRoundTrip(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			pgxtest.RunValueRoundTripTests(context.Background(), t, defaultConnTestRunner, nil, typ, []pgxtest.ValueRoundTripTest{
				{
					Param:  wrapperspb.String("abc123"),
					Result: new(*wrapperspb.StringValue),
					Test:   isExpectedEqStringValue(wrapperspb.String("abc123")),
				},
			})
		})
	}
}
