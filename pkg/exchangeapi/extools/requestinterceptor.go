package extools

import (
	"exterior-interactor/pkg/httptools"
	"net/http"
)

type RequestInterceptor interface {
	BeforeRequest(meta Meta, request httptools.Request) error
	AfterRequest(meta Meta, request httptools.Request, rsp *http.Response) error
}
