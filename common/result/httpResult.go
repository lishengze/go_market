package result

import (
	"bcts/common/crypto"
	"bcts/common/xerror"
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/status"
	"net/http"
)

//http返回
func HttpResult(r *http.Request, w http.ResponseWriter, resp interface{}, err error, aesKey string, isEncrypted bool) {
	//加密
	if isEncrypted {
		var responsePlainMsg []byte
		if err == nil {
			r := Success(resp)
			responsePlainMsg, err = json.Marshal(r)
			if err != nil {
				logx.Errorf("AddTest Marshal response  err:%+v", err)

				errRsp := Error(xerror.ErrorTryAgain.Code, xerror.ErrorTryAgain.ErrorMsg(), xerror.ErrorTryAgain.ErrorMsgCHS())

				responsePlainMsg, _ = json.Marshal(errRsp)
			}

		} else {
			//错误返回
			errCode := xerror.ErrorTryAgain.ErrCode()
			errMsg := xerror.ErrorTryAgain.ErrorMsg()
			errMsgCHS := xerror.ErrorTryAgain.ErrorMsgCHS()

			causeErr := errors.Cause(err)                  // err类型
			if e, ok := causeErr.(*xerror.CodeError); ok { //自定义错误类型
				errCode = e.ErrCode()
				errMsg = e.ErrorMsg()
				errMsgCHS = e.ErrorMsgCHS()
			} else {
				if grpcStatus, ok := status.FromError(causeErr); ok { // grpc err错误
					grpcCode := int(grpcStatus.Code())
					if res := xerror.GetCodeError(grpcCode); res != nil {
						errCode = res.ErrCode()
						errMsg = res.ErrorMsg()
						errMsgCHS = res.ErrorMsgCHS()
					}
				}
			}

			//
			errRsp := Error(errCode, errMsg, errMsgCHS)
			responsePlainMsg, _ = json.Marshal(errRsp)
		}

		encryptedRes := crypto.EncryptedResponse(string(responsePlainMsg), aesKey)
		w.Header().Set("X-Encrypted", "true")
		httpx.WriteJson(w, http.StatusOK, encryptedRes)
	} else {
		if err == nil { //成功返回
			r := Success(resp)
			httpx.WriteJson(w, http.StatusOK, r)
		} else {
			//错误返回
			errCode := xerror.ErrorTryAgain.ErrCode()
			errMsg := xerror.ErrorTryAgain.ErrorMsg()
			errMsgCHS := xerror.ErrorTryAgain.ErrorMsgCHS()

			causeErr := errors.Cause(err)                  // err类型
			if e, ok := causeErr.(*xerror.CodeError); ok { //自定义错误类型
				errCode = e.ErrCode()
				errMsg = e.ErrorMsg()
				errMsgCHS = e.ErrorMsgCHS()
			} else {
				if grpcStatus, ok := status.FromError(causeErr); ok { // grpc err错误
					grpcCode := int(grpcStatus.Code())
					if res := xerror.GetCodeError(grpcCode); res != nil {
						errCode = res.ErrCode()
						errMsg = res.ErrorMsg()
						errMsgCHS = res.ErrorMsgCHS()
					}
				}
			}
			httpx.WriteJson(w, http.StatusOK, Error(errCode, errMsg, errMsgCHS))
		}
	}
}

//http返回
func HttpResultBak(r *http.Request, w http.ResponseWriter, resp interface{}, err error) {
	if err == nil { //成功返回
		r := Success(resp)
		httpx.WriteJson(w, http.StatusOK, r)
	} else {
		//错误返回
		errCode := xerror.ErrorTryAgain.ErrCode()
		errMsg := xerror.ErrorTryAgain.ErrorMsg()
		errMsgCHS := xerror.ErrorTryAgain.ErrorMsgCHS()

		causeErr := errors.Cause(err)                  // err类型
		if e, ok := causeErr.(*xerror.CodeError); ok { //自定义错误类型
			errCode = e.ErrCode()
			errMsg = e.ErrorMsg()
			errMsgCHS = e.ErrorMsgCHS()
		} else {
			if grpcStatus, ok := status.FromError(causeErr); ok { // grpc err错误
				grpcCode := int(grpcStatus.Code())
				if res := xerror.GetCodeError(grpcCode); res != nil {
					errCode = res.ErrCode()
					errMsg = res.ErrorMsg()
					errMsgCHS = res.ErrorMsgCHS()
				}
			}
		}
		httpx.WriteJson(w, http.StatusOK, Error(errCode, errMsg, errMsgCHS))
	}
}

//http 参数错误返回
func ParamErrorResultBak(r *http.Request, w http.ResponseWriter, err error) {
	httpx.WriteJson(w, http.StatusOK, Error(xerror.ErrorParamError.ErrCode(),
		xerror.ErrorParamError.ErrorMsg(),
		xerror.ErrorParamError.ErrorMsgCHS()))
}

//http 参数错误返回
func ParamErrorResult(r *http.Request, w http.ResponseWriter, err error, aesKey string, isEncrypted bool) {
	if isEncrypted {
		errRsp := Error(xerror.ErrorTryAgain.Code, xerror.ErrorTryAgain.ErrorMsg(), xerror.ErrorTryAgain.ErrorMsgCHS())
		responsePlainMsg, _ := json.Marshal(errRsp)
		encryptedRes := crypto.EncryptedResponse(string(responsePlainMsg), aesKey)
		w.Header().Set("X-Encrypted", "true")
		httpx.WriteJson(w, http.StatusOK, encryptedRes)
	} else {
		httpx.WriteJson(w, http.StatusOK, Error(xerror.ErrorParamError.ErrCode(),
			xerror.ErrorParamError.ErrorMsg(),
			xerror.ErrorParamError.ErrorMsgCHS()))
	}
}
