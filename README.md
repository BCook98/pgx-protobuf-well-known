# pgx-protobuf-well-known

PGX Type support for [Protobuf Well-Known Types](https://protobuf.dev/reference/protobuf/google.protobuf/).

This package currently supports the following well-known types:

- [Timestamp](https://protobuf.dev/reference/protobuf/google.protobuf/#timestamp) as both PostgreSQL `TIMESTAMP` and `TIMESTAMPTZ`

## Usage

### All Types

You can register all of the supported well known types using `pgxprotobufwellknown.Register(conn.TypeMap())`

Example:

```go
import (
	"github.com/jackc/pgx/v5"
	pgxprotobufwellknown "github.com/bcook98/pgx-protobuf-well-known"
)

myPool.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	pgxprotobufwellknown.Register(conn.TypeMap())
}
```

### Specific Types

Alternatively you can register the specific types you want to utilise.

```go
import (
	"github.com/jackc/pgx/v5"

	pgxprotobuftimestamppb "github.com/bcook98/pgx-protobuf-well-known/pkg/timestamppb"
)

myPool.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
	// Timestamp - TIMESTAMP
	pgxprotobuftimestamppb.RegisterTimestamp(conn.TypeMap())
	// Timestamp - TIMESTAMPTZ
	pgxprotobuftimestamppb.RegisterTimestamptz(conn.TypeMap())
}
```

## Testing

To test this package locally you can either run PostgreSQL natively on your system, and run the `go test` command as you would normally.

Or you can run PostgreSQL in docker using the supplied Docker Compose file, then provide the `TEST_PG_CONN_STRING` environment variable.

Example:

`TEST_PG_CONN_STRING="host=localhost port=5432 user=pgx_protobuf_well_known password=pgx_protobuf_well_known dbname=pgx_protobuf_well_known sslmode=disable" go test ./...`
