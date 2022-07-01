module github.com/8run0/otel/backend

go 1.18

require go.opentelemetry.io/otel v1.7.0

require (
	github.com/felixge/httpsnoop v1.0.2 // indirect
	github.com/mxk/go-sqlite v0.0.0-20140611214908-167da9432e1f // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0 // indirect
	go.opentelemetry.io/otel/sdk v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
)

require (
	github.com/go-chi/chi v1.5.4
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-chi/render v1.0.1
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/jmoiron/sqlx v1.3.5
	github.com/mattn/go-sqlite3 v1.14.14
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.7.0
	go.opentelemetry.io/otel/metric v0.30.0
	go.opentelemetry.io/otel/trace v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
)
