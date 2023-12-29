package gapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"

	mockdb "github.com/ZhangZhihuiAAA/zimplebank/db/mock"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/pb"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/ZhangZhihuiAAA/zimplebank/worker"
	mockworker "github.com/ZhangZhihuiAAA/zimplebank/worker/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateUserAPI(t *testing.T) {
    user, password := randomUser(t)

    testCases := []struct {
        name          string
        req           *pb.CreateUserRequest
        buildStubs    func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor)
        checkResponse func(t *testing.T, res *pb.CreateUserResponse, err error)
    }{
        {
            name: "OK",
            req: &pb.CreateUserRequest{
                Username: user.Username,
                Password: password,
                FullName: user.FullName,
                Email:    user.Email,
            },
            buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
                arg := db.CreateUserTxParams{
                    CreateUserParams: db.CreateUserParams{
                        Username: user.Username,
                        FullName: user.FullName,
                        Email:    user.Email,
                    },
                }
                store.EXPECT().
                    CreateUserTx(gomock.Any(), eqCreateUserTxParams(arg, password, user)).
                    Times(1).
                    Return(db.CreateUserTxResult{User: user}, nil)

                taskPayload := &worker.PayloadSendVerificationEmail{
                    Username: user.Username,
                }
                taskDistributor.EXPECT().
                    DistributeTaskSendVerificationEmail(gomock.Any(), taskPayload, gomock.Any()).
                    Times(1).
                    Return(nil)
            },
            checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
                require.NoError(t, err)
                require.NotNil(t, res)
                createdUser := res.GetUser()
                require.Equal(t, user.Username, createdUser.Username)
                require.Equal(t, user.FullName, createdUser.FullName)
                require.Equal(t, user.Email, createdUser.Email)
            },
        },
        {
            name: "InternalError",
            req: &pb.CreateUserRequest{
                Username: user.Username,
                Password: password,
                FullName: user.FullName,
                Email:    user.Email,
            },
            buildStubs: func(store *mockdb.MockStore, taskDistributor *mockworker.MockTaskDistributor) {
                store.EXPECT().
                    CreateUserTx(gomock.Any(), gomock.Any()).
                    Times(1).
                    Return(db.CreateUserTxResult{}, errors.New("Internal Error"))

                taskDistributor.EXPECT().
                    DistributeTaskSendVerificationEmail(gomock.Any(), gomock.Any(), gomock.Any()).
                    Times(0)
            },
            checkResponse: func(t *testing.T, res *pb.CreateUserResponse, err error) {
                require.Error(t, err)
                st, ok := status.FromError(err)
                require.True(t, ok)
                require.Equal(t, codes.Internal, st.Code())
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl1 := gomock.NewController(t)
            defer ctrl1.Finish()
            store := mockdb.NewMockStore(ctrl1)

            ctrl2 := gomock.NewController(t)
            defer ctrl2.Finish()
            taskdDstributor := mockworker.NewMockTaskDistributor(ctrl2)

            tc.buildStubs(store, taskdDstributor)

            server := newTestServer(t, store, taskdDstributor)
            res, err := server.CreateUser(context.Background(), tc.req)
            tc.checkResponse(t, res, err)
        })
    }
}

func randomUser(t *testing.T) (user db.User, password string) {
    password = util.RandomString(6)
    hashedPassword, err := util.HashPassword(password)
    require.NoError(t, err)

    user = db.User{
        Username:       util.RandomOwner(),
        HashedPassword: hashedPassword,
        FullName:       util.RandomOwner(),
        Email:          util.RandomEmail(),
    }
    return
}

type eqCreateUserTxParamsMatcher struct {
    arg      db.CreateUserTxParams
    password string
    user     db.User
}

func (e eqCreateUserTxParamsMatcher) Matches(x any) bool {
    arg, ok := x.(db.CreateUserTxParams)
    if !ok {
        return false
    }

    err := util.CheckPassword(e.password, arg.HashedPassword)
    if err != nil {
        return false
    }

    e.arg.HashedPassword = arg.HashedPassword
    if !reflect.DeepEqual(e.arg.CreateUserParams, arg.CreateUserParams) {
        return false
    }

    err = arg.AfterCreate(e.user)
    return err == nil
}

func (e eqCreateUserTxParamsMatcher) String() string {
    return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

func eqCreateUserTxParams(arg db.CreateUserTxParams, password string, user db.User) gomock.Matcher {
    return eqCreateUserTxParamsMatcher{arg, password, user}
}
