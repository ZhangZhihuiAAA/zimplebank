package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ZhangZhihuiAAA/zimplebank/api"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	_ "github.com/ZhangZhihuiAAA/zimplebank/doc/statik"
	"github.com/ZhangZhihuiAAA/zimplebank/gapi"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
    RETRY_INTERVAL = 5 * time.Second
    RETRY_LIMIT    = 10
)

var config util.Config
var connPool *pgxpool.Pool
var store db.Store

func init() {
    var err error
    config, err = util.LoadConfig(".")
    if err != nil {
        log.Fatal().Err(err).Msg("failed to load config")
    }

    var dbCtx = context.Background()
    connPool, err = pgxpool.New(dbCtx, config.DBSource)
    if err != nil {
        log.Fatal().Err(err)
    }

    var retryTimes = 0
CONNECT_DB:
    _, err = connPool.Query(dbCtx, "SELECT * FROM users LIMIT 1;")
    if err != nil {
        if strings.Contains(err.Error(), "failed to connect to") {
            if retryTimes < RETRY_LIMIT {
                retryTimes++
                log.Info().Msg("retry connecting to db ......")
                time.Sleep(RETRY_INTERVAL)
                goto CONNECT_DB
            } else {
                log.Fatal().Err(err).Msg("failed to connect to db")
            }
        }

        if strings.Contains(err.Error(), "does not exist") {
            data, err := os.ReadFile(config.DBInitSchemaFile)
            if err != nil {
                log.Fatal().Err(err).Msgf("failed to read file: %s", config.DBInitSchemaFile)
            }
            connPool.Exec(dbCtx, string(data))
            if err != nil {
                log.Fatal().Err(err).Msg("error occurred when init schema")
            }
        } else {
            log.Fatal().Err(err).Msg("error occurred")
        }
    }

    store = db.NewStore(connPool)
}

func main() {
    if config.Environment == "DEV" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }

    go runGatewayServer()
    runGrpcServer()
}

func runGinServer() {
    server, err := api.NewServer(config, store)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create api server")
    }

    log.Info().Msg("Start api server")
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start api server:")
    }
}

func runGatewayServer() {
    server, err := gapi.NewServer(config, store)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create gapi server")
    }

    jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
        MarshalOptions: protojson.MarshalOptions{
            UseProtoNames: true,
        },
        UnmarshalOptions: protojson.UnmarshalOptions{
            DiscardUnknown: true,
        },
    })

    grpcMux := runtime.NewServeMux(jsonOption)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    err = pb.RegisterZimpleBankHandlerServer(ctx, grpcMux, server)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to register handler server")
    }

    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    statikFS, err := fs.New()
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create statik fs:")
    }

    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
    mux.Handle("/swagger/", swaggerHandler)

    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create listener")
    }

    log.Info().Msgf("Start HTTP gateway server at %s", listener.Addr().String())
    handler := gapi.HttpLogger(mux)
    err = http.Serve(listener, handler)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start HTTP gateway server")
    }
}

func runGrpcServer() {
    server, err := gapi.NewServer(config, store)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create gapi server")
    }

    grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
    grpcServer := grpc.NewServer(grpcLogger)
    pb.RegisterZimpleBankServer(grpcServer, server)

    // This step is optional, but is highly recommended.
    // It allows GRPC client easily explore what RPCs are available on the server and how to call them.
    reflection.Register(grpcServer)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to create listener")
    }

    log.Info().Msgf("Start GRPC server at %s", listener.Addr().String())
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal().Err(err).Msg("failed to start GRPC server")
    }
}
