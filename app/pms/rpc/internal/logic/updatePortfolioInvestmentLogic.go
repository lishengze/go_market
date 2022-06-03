package logic

import (
	"bcts/app/pms/model"
	"bcts/common/xerror"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"time"

	"bcts/app/pms/rpc/internal/svc"
	"bcts/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePortfolioInvestmentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePortfolioInvestmentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePortfolioInvestmentLogic {
	return &UpdatePortfolioInvestmentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePortfolioInvestmentLogic) UpdatePortfolioInvestment(in *pb.UpdatePortfolioInvestmentReq) (*pb.EmptyRsp, error) {
	res, err := l.svcCtx.PortfolioInvestmentsModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerror.ErrorDB, "UpdatePortfolioInvestment find portfolio_investments db err , id: %d , err: %v", in.Id, err)
	}
	if res == nil {
		return nil, errors.Wrapf(xerror.ErrorParamError, "portfolio_investments is not exist, id: %d", in.Id)
	}
	//判断此用户下，此name 是否存在
	whereBuilder := l.svcCtx.PortfolioInvestmentsModel.RowBuilder()
	whereBuilder = whereBuilder.Where(squirrel.Eq{"user_id": in.UserID})
	resAll, err := l.svcCtx.PortfolioInvestmentsModel.FindAll(l.ctx, whereBuilder, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Update PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	var gid []int64
	if len(resAll) > 0 {
		for _, v := range resAll {
			if v.Id == in.Id {
				continue
			}
			gid = append(gid, v.Id) //此用户的其他资产组
			if v.Name == in.PortfolioName {
				return nil, errors.Wrapf(xerror.ErrorParamError, "portfolio_investments name is already exist, id: %d", v.Id)
			}
		}
	}
	//校验资产组币种是否重复
	itemBuilderOther := l.svcCtx.PortfolioInvestmentItemsModel.RowBuilder()
	itemBuilderOther = itemBuilderOther.Where(squirrel.Eq{"gid": gid, "investment": in.Investments})
	itemRspOther, err := l.svcCtx.PortfolioInvestmentItemsModel.FindAll(l.ctx, itemBuilderOther, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Update PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	if len(itemRspOther) > 0 {
		return nil, errors.Wrapf(xerror.ErrorParamError, "investment already repeat, investment: %s", itemRspOther[0].Investment)
	}
	//更新当前资产组
	itemBuilder := l.svcCtx.PortfolioInvestmentItemsModel.RowBuilder()
	itemBuilder = itemBuilder.Where(squirrel.Eq{"gid": in.Id})
	itemRsp, err := l.svcCtx.PortfolioInvestmentItemsModel.FindAll(l.ctx, itemBuilder, "")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Update PortfolioInvestment err: %+v , in :%+v", err, in)
	}
	var currents []string //库中此组存在的币种
	if len(itemRsp) > 0 {
		for _, v := range itemRsp {
			currents = append(currents, v.Investment)
		}
	}
	addedInvestments := substr(in.Investments, currents)
	deletedInvestments := substr(currents, in.Investments)

	var delId []int64
	if len(deletedInvestments) > 0 {
		for _, v := range deletedInvestments {
			for _, item := range itemRsp {
				if v == item.Investment {
					delId = append(delId, item.Id)
				}
			}
		}
	}

	nowTime := time.Now().UTC()
	if err = l.svcCtx.PortfolioInvestmentsModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		//删除
		if len(delId) > 0 {
			for _, v := range delId {
				if err = l.svcCtx.PortfolioInvestmentItemsModel.Delete(ctx, session, v); err != nil {
					return errors.Wrapf(xerror.ErrorDB, "delete db portfolio_investment_items  err: %v,id: %d", err, v)
				}
			}
		}
		//更新 portfolio_investment
		res.Name = in.PortfolioName
		res.Updated = nowTime
		res.Remark = in.Remark
		res.LastOperator = in.UserID
		_, err = l.svcCtx.PortfolioInvestmentsModel.Update(ctx, session, res)
		if err != nil {
			return errors.Wrapf(xerror.ErrorDB, "db portfolio_investments Update err:%v,portfolioInvestment:%+v", err, res)
		}
		//插入新增加币种
		if len(addedInvestments) > 0 {
			for _, v := range addedInvestments {
				investmentItem := new(model.PortfolioInvestmentItems)
				investmentItem.Gid = in.Id
				investmentItem.Investment = v
				investmentItem.Created = nowTime
				investmentItem.Updated = nowTime
				if _, err := l.svcCtx.PortfolioInvestmentItemsModel.Insert(ctx, session, investmentItem); err != nil {
					return errors.Wrapf(xerror.ErrorDB, "db portfolio_investment_items Insert err:%v,PortfolioInvestmentItems:%v", err, investmentItem)
				}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &pb.EmptyRsp{}, nil
}

/**
求a中出现，b中未出现的
*/
func substr(a []string, b []string) []string {
	var c []string
	temp := map[string]struct{}{}

	for _, val := range b {
		if _, ok := temp[val]; !ok {
			temp[val] = struct{}{} // 空struct 不占内存空间
		}
	}

	for _, val := range a {
		if _, ok := temp[val]; !ok {
			c = append(c, val)
		}
	}

	return c
}
