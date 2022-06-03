package order

import (
	"encoding/json"
	"io/ioutil"
	"market_server/common/crypto"
	"market_server/common/result"
	"net/http"

	"market_server/app/client/api/internal/logic/order"
	"market_server/app/client/api/internal/svc"
	"market_server/app/client/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func OtcQuoteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Infof("header:%+v", r.Header)
		if r.Method == "POST" {
			iv := r.Header.Get("X-Encrypt-Iv")
			//1. 可以认为是直接调用
			if len(iv) == 0 {
				var req types.OtcQuoteReq
				if err := httpx.Parse(r, &req); err != nil {
					result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, false)
					return
				}

				l := order.NewOtcQuoteLogic(r, r.Context(), svcCtx)
				resp, err := l.OtcQuote(&req)
				result.HttpResult(r, w, resp, err, svcCtx.PemFileBase.AESKey, false)
				return
			}

			//2. 需要做加解密的操作
			requestBody, _ := ioutil.ReadAll(r.Body)

			var encrypted crypto.EncryptedReqBody
			if err := json.Unmarshal(requestBody, &encrypted); err != nil {
				result.HttpResult(r, w, nil, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
				return
			}
			logx.Infof("after unmarshal: %+v", encrypted)

			requestPlainMsg, err := crypto.DecryptReqBody(&encrypted, svcCtx.PemFileBase.AESKey, iv)
			if err != nil {
				logx.Errorf("DecryptReqBody err %+v", err)
				result.HttpResult(r, w, nil, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
				return
			}
			logx.Info(string(requestPlainMsg))

			var req types.OtcQuoteReq
			if err := json.Unmarshal(requestPlainMsg, &req); err != nil {
				logx.Errorf("OtcQuoteLogic Unmarshal err:%+v", err)
				result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
				return
			}

			logx.Infof("OtcQuoteLogic req:%+v", req)

			l := order.NewOtcQuoteLogic(r, r.Context(), svcCtx)
			resp, err := l.OtcQuote(&req)
			result.HttpResult(r, w, resp, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
		}

		if r.Method == "GET" {
			iv := r.Header.Get("X-Encrypt-Iv")
			logx.Infof("OtcQuoteLogic iv:%s", iv)

			var req types.OtcQuoteReq
			if err := httpx.Parse(r, &req); err != nil {
				result.ParamErrorResult(r, w, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
				return
			}

			l := order.NewOtcQuoteLogic(r, r.Context(), svcCtx)
			resp, err := l.OtcQuote(&req)
			result.HttpResult(r, w, resp, err, svcCtx.PemFileBase.AESKey, svcCtx.Config.ResponseEncrypted)
		}
	}
}
