package main

import (
	"context"
	"log"
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
        log.Fatal("failed to load config:", err)
    }

    var dbCtx = context.Background()
    connPool, err = pgxpool.New(dbCtx, config.DBSource)
    if err != nil {
        log.Fatal(err)
    }

    var retryTimes = 0
CONNECT_DB:
    _, err = connPool.Query(dbCtx, "SELECT * FROM users LIMIT 1;")
    if err != nil {
        if strings.Contains(err.Error(), "failed to connect to") {
            if retryTimes < RETRY_LIMIT {
                retryTimes++
                log.Println("retry connecting to db ......")
                time.Sleep(RETRY_INTERVAL)
                goto CONNECT_DB
            } else {
                log.Fatal("failed to connect to db:", err)
            }
        }

        if strings.Contains(err.Error(), "does not exist") {
            data, err := os.ReadFile(config.DBInitSchemaFile)
            if err != nil {
                log.Fatalf("failed to read file: %s", config.DBInitSchemaFile)
            }
            connPool.Exec(dbCtx, string(data))
            if err != nil {
                log.Fatal("error occurred when init schema:", err)
            }
        } else {
            log.Fatal("error occurred:", err)
        }
    }

    store = db.NewStore(connPool)
}

func main() {
    go runGatewayServer()
    runGrpcServer()
}

func runGinServer() {
    server, err := api.NewServer(config, store)
    if err != nil {
        log.Fatal("failed to create api server:", err)
    }

    log.Println("Start api server")
    err = server.Start(config.HTTPServerAddress)
    if err != nil {
        log.Fatal("failed to start api server:", err)
    }
}

func runGatewayServer() {
    server, err := gapi.NewServer(config, store)
    if err != nil {
        log.Fatal("failed to create gapi server:", err)
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
        log.Fatal("failed to register handler server:", err)
    }

    mux := http.NewServeMux()
    mux.Handle("/", grpcMux)

    statikFS, err := fs.New()
    if err != nil {
        log.Fatal("failed to create statik fs:", err)
    }

    swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
    mux.Handle("/swagger/", swaggerHandler)

    listener, err := net.Listen("tcp", config.HTTPServerAddress)
    if err != nil {
        log.Fatal("failed to create listener:", err)
    }

    log.Printf("Start HTTP gateway server at %s\n", listener.Addr().String())
    err = http.Serve(listener, mux)
    if err != nil {
        log.Fatal("failed to start HTTP gateway server:", err)
    }
}

func runGrpcServer() {
    server, err := gapi.NewServer(config, store)
    if err != nil {
        log.Fatal("failed to create gapi server:", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterZimpleBankServer(grpcServer, server)

    // This step is optional, but is highly recommended.
    // It allows GRPC client easily explore what RPCs are available on the server and how to call them.
    reflection.Register(grpcServer)

    listener, err := net.Listen("tcp", config.GRPCServerAddress)
    if err != nil {
        log.Fatal("failed to create listener:", err)
    }

    log.Printf("Start GRPC server at %s\n", listener.Addr().String())
    err = grpcServer.Serve(listener)
    if err != nil {
        log.Fatal("failed to start GRPC server:", err)
    }
}
