package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"market_server/app/admin/api/internal/logic/rbac"
	"market_server/app/admin/model"
	"market_server/common/middleware"
	"market_server/common/result"
	"market_server/common/tool"
	"market_server/common/xerror"
	"net/http"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RbacMiddleware struct {
	ctx              context.Context
	RbacLogic        *rbac.RbacLogic
	MenuModel        model.MenuModel
	SecurityLogModel model.SecurityLogModel
}

func NewRbacMiddleware(conn sqlx.SqlConn) *RbacMiddleware {
	return &RbacMiddleware{
		ctx:              context.Background(),
		RbacLogic:        rbac.NewRbacLogic(context.Background(), conn),
		MenuModel:        model.NewMenuModel(conn),
		SecurityLogModel: model.NewSecurityLogModel(conn),
	}
}

var (
	whiteList = map[string]struct{}{
		"POST:/api/login":                      {}, //不需要登录，单独记录secure_log日志
		"POST:/api/qr_code":                    {}, //不需要登录，不记录secure_log日志
		"GET:/api/version":                     {}, //不需要登录
		"GET:/api/rbac/menu":                   {}, //需要登录，但是所有登录用户都可以访问
		"GET:/api/user/info":                   {}, //需要登录，但是所有登录用户都可以访问
		"PUT:/api/user/passwd":                 {}, //(修改自己的密码)需要登录，但是所有登录用户都可以访问
		"GET:/api/rbac/roles":                  {}, //(role_id&role_name列表)需要登录，但是所有登录用户都可以访问
		"GET:/api/parameter/currency/ids":      {},
		"GET:/api/parameter/hedge/ids":         {},
		"GET:/api/group/dropdown":              {},
		"GET:/api/parameter/currency/dropdown": {}, //公用接口，用于代替 GET:/api/parameter/currency(此接口受角色权限控制)

		"GET:/api/parameter/symbol":             {},
		"GET:/api/parameter/symbol/ids":         {},
		"GET:/api/member/organinfo":             {}, // 客户信息管理--认证
		"GET:/api/fiatsettle/twofiatsettleinfo": {}, // 法币结算信息管理
		"GET:/api/fiat/memberinfo":              {},
		"POST:/api/member/uploadnoidx":          {},
		"GET:/api/parameter/symbol/newids":      {}, //  获取品种id列表 分为结算币和基础币 symbol
		"GET:/api/hedge/risk/monitor/export":    {}, //    导出对冲交易风险监控
		"GET:/api/operate/settle/tpl/download":  {}, //  运营管理结算对账模板下载
		"POST:/api/operate/settle/upload":       {}, //       上传结算对账模板

		"GET:/api/report/trade/Export": {},
		//hedge
		"GET:/api/hedge/trade/export": {},
		//order
		"GET:/api/order/all/export":     {},
		"GET:/api/order/trade":          {},
		"GET:/api/order/trade/export":   {},
		"GET:/api/order/current/export": {},
		//operate
		"DELETE:/api/operate/offlinetrade": {},
		//member_trade
		"GET:/api/member/trade/parameter/groups": {},
		//member_deposit
		"GET:/api/deposit_export2": {},
		//member_auth
		"GET:/api/member/trade/authority/groups": {},
		//member
		"GET:/api/member/export":        {},
		"GET:/api/member/file/export":   {},
		"GET:/api/member/twoorganinfos": {},
		"POST:/api/member/authorgan":    {},
		"POST:/api/member/authperson":   {},
		"POST:/api/member/finalapprove": {},
		"POST:/api/member/upload":       {},
		//fiat_recharge
		"POST:/api/fiat/upload": {},
		//comm
		"GET:/api/comm/member": {},
		//auth
		"GET:/api/auth/check": {},
		//assetsplat_export
		//"GET:/api/member/assetsplat_export": {},//客户资产快照--导出 ————已经移动到数据库配置里
	}
)

func (m *RbacMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		currentPath := strings.ToLower(r.Method + ":" + r.URL.Path)
		//println(currentPath)
		//白名单列表
		if _, ok := whiteList[currentPath]; ok {
			return
		}

		userIdVal := r.Context().Value(middleware.JWT_USER_ID)
		if userIdVal == nil {
			result.HttpResult(r, w, nil, xerror.ErrorRbacUserOpIdNotFound)
			return
		}
		userId := userIdVal.(json.Number)
		user, err := userId.Int64()
		if err != nil {
			return
		}
		userOperations, err := m.RbacLogic.FindUserOperations(user)
		if err != nil {
			result.HttpResult(r, w, nil, xerror.ErrorRbacUserIdNotFound)
			return
		}

		if strings.Index(currentPath, "get:") == 0 { //不记录HTTP GET请求
			next(w, r)
			return
		}

		for _, userOperation := range userOperations {
			if strings.ToLower(userOperation.Path) == currentPath {
				var secureLogContent = userOperation.MenuNameEn
				if userOperation.ParentId > 0 {
					parentMenu, err := m.MenuModel.FindOne(m.ctx, userOperation.ParentId)
					if err == nil {
						secureLogContent = parentMenu.MenuNameEn + "-" + secureLogContent
					}
				}
				secureLog := &model.SecurityLog{
					Path:     currentPath,
					Operator: user,
					Ip:       tool.GetClientIp(r),
					Content:  secureLogContent,
					Created:  time.Now(),
					Updated:  time.Now(),
				}
				_, err = m.SecurityLogModel.Insert(m.ctx, nil, secureLog)
				if err != nil {
					println(fmt.Sprintf("middleware save secure_log falied. path:%s, user_id:%d", currentPath, user), err)
				}
			}
		}

		next(w, r)
	}
}
