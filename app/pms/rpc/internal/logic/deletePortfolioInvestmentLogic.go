package logic

import (
	"bcts/common/xerror"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"bcts/app/pms/rpc/internal/svc"
	"bcts/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePortfolioInvestmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePortfolioInvestmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePortfolioInvestmentLogic {
	return &DeletePortfolioInvestmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeletePortfolioInvestmentLogic) DeletePortfolioInvestment(in *pb.DeletePortfolioInvestmentReq) (*pb.EmptyRsp, error) {
	//查出item表的id
	itemWhereBuilder := l.svcCtx.PortfolioInvestmentItemsModel.RowBuilder()
	itemWhereBuilder = itemWhereBuilder.Where(squirrel.Eq{"gid": in.Id})
	itemRsp, err := l.svcCtx.PortfolioInvestmentItemsModel.FindAll(l.ctx, itemWhereBuilder, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Create PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	if err = l.svcCtx.PortfolioInvestmentsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		err = l.svcCtx.PortfolioInvestmentsModel.Delete(ctx, session, in.Id)
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "delete db portfolio_investments  err: %v, id: %d", err, in.Id)
		}
		if len(itemRsp) > 0 {
			for _, v := range itemRsp {
				if err = l.svcCtx.PortfolioInvestmentItemsModel.Delete(ctx, session, v.Id); err != nil {
					return errors.Wrapf(xerror.ErrorDB, "delete db portfolio_investment_items  err:%v,gid:%d", err, in.Id)
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return &pb.EmptyRsp{}, nil
}
