package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockdb "github.com/ZhangZhihuiAAA/zimplebank/db/mock"
	db "github.com/ZhangZhihuiAAA/zimplebank/db/sqlc"
	"github.com/ZhangZhihuiAAA/zimplebank/token"
	"github.com/ZhangZhihuiAAA/zimplebank/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestTransferAPI(t *testing.T) {
    amount := 10.00

    user1, _ := randomUser(t)
    user2, _ := randomUser(t)
    user3, _ := randomUser(t)

    account1 := randomAccount(user1.Username)
    account2 := randomAccount(user2.Username)
    account3 := randomAccount(user3.Username)

    account1.Currency = util.CNY
    account2.Currency = util.CNY
    account3.Currency = util.USD

    testCases := []struct {
        name          string
        body          gin.H
        setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
        buildStubs    func(store *mockdb.MockStore)
        checkResponse func(recorder *httptest.ResponseRecorder)
    }{
        {
            name: "OK",
            body: gin.H{
                "from_account_id": account1.ID,
                "to_account_id": account2.ID,
                "amount": amount,
                "currency": account1.Currency,
            },
            setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
                addAuthorization(t, request, tokenMaker, AUTHORIZATION_TYPE_BEARER, user1.Username, user1.Role, time.Minute)
            },
            buildStubs: func(store *mockdb.MockStore) {
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account1.ID)).
                    Times(1).
                    Return(account1, nil)
                store.EXPECT().
                    GetAccount(gomock.Any(), gomock.Eq(account2.ID)).
                    Times(1).
                    Return(account2, nil)
                
                arg := db.TransferTxParams{
                    FromAccountID: account1.ID,
                    ToAccountID: account2.ID,
                    Amount: amount,
                }
                store.EXPECT().
                    TransferTx(gomock.Any(), gomock.Eq(arg)).
                    Times(1)
            },
            checkResponse: func(recorder *httptest.ResponseRecorder) {
                require.Equal(t, http.StatusOK, recorder.Code)
            },
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            ctrl := gomock.NewController(t)
            defer ctrl.Finish()

            store := mockdb.NewMockStore(ctrl)
            tc.buildStubs(store)

            server := newTestServer(t, store)
            recorder := httptest.NewRecorder()

            // Marshal body data to JSON
            data, err := json.Marshal(tc.body)
            require.NoError(t, err)

            url := "/transfers"
            request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(data))
            require.NoError(t, err)

            tc.setupAuth(t, request, server.tokenMaker)
            server.router.ServeHTTP(recorder, request)
            tc.checkResponse(recorder)
        })
    }

}
