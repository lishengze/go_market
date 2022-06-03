//KqMessage
package kqueue

// 回调状态通知
type ThirdPaymentUpdateKycStatusNotifyMessage struct {
	KycStatus int64  `json:"KycStatus"`
	Sn        string `json:"Sn"`
}
