package pgxprotobufwellknown_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	pgxprotobufwellknown "github.com/bcook98/pgx-protobuf-well-known"
	testutil "github.com/bcook98/pgx-protobuf-well-known/pkg/test_util"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type KitchenSink struct {
	ID int64 `json:"id"`

	TS      *timestamppb.Timestamp   `json:"ts"`
	TSSlice []*timestamppb.Timestamp `json:"ts_slice"`

	TSTZ      *timestamppb.Timestamp   `json:"ts_tz"`
	TSTZSlice []*timestamppb.Timestamp `json:"ts_tz_slice"`
}

func TestKitchenSink(t *testing.T) {
	testRunner := testutil.ConnTestRunner()
	testRunner.AfterConnect = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		pgxprotobufwellknown.Register(conn.TypeMap())
	}

	tableName := fmt.Sprintf("test_kitchen_sink_%d", time.Now().UnixNano())

	testRunner.AfterTest = func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		_, err := conn.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		require.NoError(t, err)
	}

	testRunner.RunTest(context.Background(), t, func(ctx context.Context, t testing.TB, conn *pgx.Conn) {
		_, err := conn.Exec(ctx, fmt.Sprintf(`
			CREATE TABLE %s (
				id BIGINT PRIMARY KEY,

				ts TIMESTAMP,
				ts_slice TIMESTAMP[],

				ts_tz TIMESTAMPTZ,
				ts_tz_slice TIMESTAMPTZ[])
			`, tableName))
		require.NoError(t, err)

		ts := timestamppb.New(time.Now())
		tsTz := timestamppb.New(time.Now().In(time.UTC))
		tsSlice := []*timestamppb.Timestamp{timestamppb.New(time.Now()), timestamppb.New(time.Now().Add(1 * time.Hour))}
		tsTzSlice := []*timestamppb.Timestamp{timestamppb.New(time.Now().In(time.UTC)), timestamppb.New(time.Now().In(time.UTC).Add(1 * time.Hour))}

		kitchenSink := KitchenSink{
			ID:        1,
			TS:        ts,
			TSSlice:   tsSlice,
			TSTZ:      tsTz,
			TSTZSlice: tsTzSlice,
		}

		_, err = conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s (id, ts, ts_slice, ts_tz, ts_tz_slice) VALUES ($1, $2, $3, $4, $5)", tableName),
			kitchenSink.ID,
			kitchenSink.TS,
			kitchenSink.TSSlice,
			kitchenSink.TSTZ,
			kitchenSink.TSTZSlice,
		)
		require.NoError(t, err)

		var result KitchenSink
		err = conn.QueryRow(ctx, fmt.Sprintf("SELECT id, ts, ts_slice, ts_tz, ts_tz_slice FROM %s WHERE id = $1", tableName), kitchenSink.ID).Scan(
			&result.ID,
			&result.TS,
			&result.TSSlice,
			&result.TSTZ,
			&result.TSTZSlice,
		)
		require.NoError(t, err)
		require.Equal(t, kitchenSink.ID, result.ID)
		require.Equal(t, kitchenSink.TS.AsTime().Unix(), result.TS.AsTime().Unix())
		require.Equal(t, kitchenSink.TSSlice[0].AsTime().Unix(), result.TSSlice[0].AsTime().Unix())
		require.Equal(t, kitchenSink.TSSlice[1].AsTime().Unix(), result.TSSlice[1].AsTime().Unix())
		require.Equal(t, kitchenSink.TSTZ.AsTime().Unix(), result.TSTZ.AsTime().Unix())
		require.Equal(t, kitchenSink.TSTZSlice[0].AsTime().Unix(), result.TSTZSlice[0].AsTime().Unix())
		require.Equal(t, kitchenSink.TSTZSlice[1].AsTime().Unix(), result.TSTZSlice[1].AsTime().Unix())
	})
}
