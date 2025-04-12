package pgxprotobufwellknown

import (
	pgxprotobuftimestamppb "github.com/bcook98/pgx-protobuf-well-known/pkg/timestamppb"
	"github.com/jackc/pgx/v5/pgtype"
)

func Register(conn *pgtype.Map) {
	pgxprotobuftimestamppb.RegisterTimestamp(conn)
	pgxprotobuftimestamppb.RegisterTimestamptz(conn)
}
