package offlinetrade

import (
	"context"
	"errors"
	"fmt"
	"market_server/app/admin/api/internal/svc"
	"market_server/app/admin/api/internal/types"
	"market_server/app/admin/model"
	"strconv"
	"strings"

	"github.com/Masterminds/squirrel"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllOfflineTradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAllOfflineTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllOfflineTradeLogic {
	return &GetAllOfflineTradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAllOfflineTradeLogic) GetAllOfflineTrade(req *types.GetAllOfflineTradeReq) (resp *types.GetAllOfflineTradeReply, err error) {
	if err = l.BeforeQuery(req); err != nil {
		return
	}

	data, count, err := l.Query(req)
	if err != nil {
		return
	}

	resp = &types.GetAllOfflineTradeReply{
		Data:  ConvertModels(data),
		Count: count,
	}
	return
}

func (l *GetAllOfflineTradeLogic) BeforeQuery(req *types.GetAllOfflineTradeReq) error {
	return validOfflineInputStatus(req.Status)
}

func (l *GetAllOfflineTradeLogic) Query(req *types.GetAllOfflineTradeReq) (entities []*model.OfflineTradeInput, count int64, err error) {
	entities = make([]*model.OfflineTradeInput, 0)
	query, args := offlineInputCondition(req)
	rowBuilder := l.svcCtx.OfflineTradeInputModel.RowBuilder().
		PlaceholderFormat(squirrel.Question).
		Where(query, args...)
	entities, err = l.svcCtx.OfflineTradeInputModel.FindPageListByPage(l.ctx, rowBuilder, req.Page, req.PageSize, "operate_time DESC")
	if err != nil {
		return
	}

	countBuilder := l.svcCtx.OfflineTradeInputModel.CountBuilder("id").
		PlaceholderFormat(squirrel.Question).
		Where(query, args...)
	count, err = l.svcCtx.OfflineTradeInputModel.FindCount(l.ctx, countBuilder)
	if err != nil {
		return
	}

	return
}

func validOfflineInputStatus(status string) (err error) {
	if status != "" {
		statusAfter := strings.Split(status, ",")
		for _, v := range statusAfter {
			s, _ := strconv.Atoi(v)
			if s == 1 || s == 2 {
				return
			} else {
				err = errors.New(fmt.Sprintf("Status:%s is invalid", status))
				return
			}
		}
	}
	return
}

func offlineInputCondition(req *types.GetAllOfflineTradeReq) (query string, args []interface{}) {
	if req.UID != "" {
		query = " (`uid` LIKE ? or `user_name` LIKE ? ) "
		args = append(args, "%"+req.UID+"%", "%"+req.UID+"%")
	}

	if req.OfflineTradeId != "" {
		if query != "" {
			query += " AND "
		}
		query += " `offline_trade_id` = ? "
		args = append(args, req.OfflineTradeId)
	}

	if req.Symbol != "" {
		if query != "" {
			query += " AND "
		}
		query += " `symbol` = ? "
		args = append(args, req.Symbol)
	}

	//全部查询
	if req.Status != "" {
		if query != "" {
			query += " AND "
		}
		query += " `status` = ? "
		args = append(args, req.Status)
	}

	if req.StartTime != "" {
		if query != "" {
			query += " AND "
		}
		query += " `operate_time` >= ? "
		args = append(args, req.StartTime)
	}

	if req.EndTime != "" {
		if query != "" {
			query += " AND "
		}
		query += " `operate_time` <= ? "
		args = append(args, req.EndTime)
	}

	return
}
