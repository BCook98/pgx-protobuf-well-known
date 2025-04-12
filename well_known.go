package pgxprotobufwellknown

import (
	pgxprotobufstringvalue "github.com/bcook98/pgx-protobuf-well-known/pkg/stringvalue"
	pgxprotobuftimestamppb "github.com/bcook98/pgx-protobuf-well-known/pkg/timestamppb"
	"github.com/jackc/pgx/v5/pgtype"
)

func Register(conn *pgtype.Map) {
	// Timestamp
	pgxprotobuftimestamppb.Register(conn)
	// StringValue
	pgxprotobufstringvalue.Register(conn)
}
