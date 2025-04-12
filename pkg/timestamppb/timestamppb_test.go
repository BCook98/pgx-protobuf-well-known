package pgxprotobuftimestamppb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	testutil "github.com/bcook98/pgx-protobuf-well-known/pkg/test_util"
	pgxtimestamppb "github.com/bcook98/pgx-protobuf-well-known/pkg/timestamppb"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxtest"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var defaultConnTestRunner pgxtest.ConnTestRunner

func init() {
	defaultConnTestRunner = testutil.ConnTestRunner()
	defaultConnTestRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxtimestamppb.RegisterTimestamp(conn.TypeMap())
		pgxtimestamppb.RegisterTimestamptz(conn.TypeMap())
	}
}

var postgresTypes = []string{
	"timestamp",
	"timestamptz",
}

func isExpectedEqTimestamp(a *timestamppb.Timestamp) func(interface{}) bool {
	return func(v interface{}) bool {
		return a.AsTime().Equal((v.(*timestamppb.Timestamp)).AsTime())
	}
}

func timeMust(t time.Time, err error) time.Time {
	if err != nil {
		panic(err)
	}
	return t
}

func TestCodecDecodeValue(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			defaultConnTestRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
				original := timestamppb.Now()

				rows, err := conn.Query(context.Background(), fmt.Sprintf(`select $1::%s`, typ), original)
				require.NoError(t, err)

				for rows.Next() {
					values, err := rows.Values()
					require.NoError(t, err)

					require.Len(t, values, 1)
					v0, ok := values[0].(timestamppb.Timestamp)
					require.True(t, ok)

					require.True(t, isExpectedEqTimestamp(original)(&v0))
				}

				require.NoError(t, rows.Err())

				rows, err = conn.Query(context.Background(), fmt.Sprintf(`select $1::%s`, typ), nil)
				require.NoError(t, err)

				for rows.Next() {
					values, err := rows.Values()
					require.NoError(t, err)

					require.Len(t, values, 1)
					require.Nil(t, values[0])
				}

				require.NoError(t, rows.Err())
			})
		})
	}
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
					isExpectedEqTimestamp(inputSlice[i])(outputSlice[i])
				}
			})
		})
	}
}

func TestValueRoundTrip(t *testing.T) {
	for _, typ := range postgresTypes {
		t.Run(typ, func(t *testing.T) {
			pgxtest.RunValueRoundTripTests(context.Background(), t, defaultConnTestRunner, nil, typ, []pgxtest.ValueRoundTripTest{
				{
					Param:  timestamppb.New(timeMust(time.Parse(time.RFC3339Nano, "2023-10-01T12:34:56.789Z"))),
					Result: new(*timestamppb.Timestamp),
					Test:   isExpectedEqTimestamp(timestamppb.New(timeMust(time.Parse(time.RFC3339Nano, "2023-10-01T12:34:56.789Z")))),
				},
			})
		})
	}
}
