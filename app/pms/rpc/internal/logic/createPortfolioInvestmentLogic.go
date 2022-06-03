package logic

import (
	"context"
	"market_server/app/pms/model"
	"market_server/app/pms/rpc/internal/svc"
	"market_server/app/pms/rpc/types/pb"
	"market_server/common/xerror"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type CreatePortfolioInvestmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePortfolioInvestmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePortfolioInvestmentLogic {
	return &CreatePortfolioInvestmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreatePortfolioInvestmentLogic) CreatePortfolioInvestment(in *pb.CreatePortfolioInvestmentReq) (*pb.EmptyRsp, error) {
	//判断是否出现过同名的
	groupType := int64(in.GroupType)
	whereBuilder := l.svcCtx.PortfolioInvestmentsModel.RowBuilder()
	whereBuilder = whereBuilder.Where(squirrel.Eq{"name": in.PortfolioName, "status": groupType, "user_id": in.UserID})
	res, err := l.svcCtx.PortfolioInvestmentsModel.FindAll(l.ctx, whereBuilder, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Create PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	if len(res) > 0 {
		return nil, errors.Wrapf(xerror.ErrorParamError, "the PortfolioName has exist, please change it for another one")
	}
	//查询 当前的id
	idWhereBuilder := l.svcCtx.PortfolioInvestmentsModel.RowBuilder()
	idWhereBuilder = idWhereBuilder.Where(squirrel.Eq{"user_id": in.UserID})
	idRes, err := l.svcCtx.PortfolioInvestmentsModel.FindAll(l.ctx, idWhereBuilder, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Create PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	//判断币种是否重复
	var ids []int64
	if len(idRes) > 0 {
		for _, v := range idRes {
			ids = append(ids, v.Id)
		}
	}
	if len(ids) > 0 {
		itemWhereBuilder := l.svcCtx.PortfolioInvestmentItemsModel.RowBuilder()
		itemWhereBuilder = itemWhereBuilder.Where(squirrel.Eq{"gid": ids, "investment": in.Investments})
		itemRsp, err := l.svcCtx.PortfolioInvestmentItemsModel.FindAll(l.ctx, itemWhereBuilder, "")
		if err != nil {
			return nil, errors.Wrapf(xerror.ErrorDB, "Create PortfolioInvestment err: %+v , in :%+v", err, in)
		}
		if len(itemRsp) > 0 {
			return nil, errors.Wrapf(xerror.ErrorParamError, "investment already repeat")
		}
	}
	//创建 存入表
	nowTime := time.Now().UTC()
	if err = l.svcCtx.PortfolioInvestmentsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		portfolioInvestment := new(model.PortfolioInvestments)
		portfolioInvestment.Name = in.PortfolioName
		portfolioInvestment.Remark = in.Remark
		portfolioInvestment.Status = groupType
		portfolioInvestment.LastOperator = in.UserID
		portfolioInvestment.UserId = in.UserID
		portfolioInvestment.Created = nowTime
		portfolioInvestment.Updated = nowTime
		insertResult, err := l.svcCtx.PortfolioInvestmentsModel.Insert(ctx, session, portfolioInvestment)
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "create db portfolio_investments Insert err:%v,portfolioInvestment:%+v", err, portfolioInvestment)
		}
		lastId, err := insertResult.LastInsertId()
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "create db portfolio_investments insertResult.LastInsertId err:%v,user:%+v", err, portfolioInvestment)
		}
		for _, investment := range in.Investments {
			investmentItem := new(model.PortfolioInvestmentItems)
			investmentItem.Gid = lastId
			investmentItem.Investment = investment
			investmentItem.Created = nowTime
			investmentItem.Updated = nowTime
			if _, err := l.svcCtx.PortfolioInvestmentItemsModel.Insert(ctx, session, investmentItem); err != nil {
				return errors.Wrapf(xerror.ErrorDB, "create db portfolio_investment_items Insert err:%v,PortfolioInvestmentItems:%v", err, investmentItem)
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.EmptyRsp{}, nil
}
