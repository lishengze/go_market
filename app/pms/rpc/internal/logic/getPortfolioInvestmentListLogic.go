package logic

import (
	"bcts/common/xerror"
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"

	"bcts/app/pms/rpc/internal/svc"
	"bcts/app/pms/rpc/types/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPortfolioInvestmentListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPortfolioInvestmentListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPortfolioInvestmentListLogic {
	return &GetPortfolioInvestmentListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPortfolioInvestmentListLogic) GetPortfolioInvestmentList(in *pb.PortfolioInvestmentListReq) (*pb.PortfolioInvestmentListRsp, error) {
	groupType := int64(in.GroupType)
	whereBuilder := l.svcCtx.PortfolioInvestmentsModel.RowBuilder()
	if in.GroupType == pb.PortfolioInvestmentListReq_ALL {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"status": []int64{1, 2}})
	} else {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"status": groupType})
	}
	if len(in.SearchContent) > 0 {
		whereBuilder = whereBuilder.Where(squirrel.Like{"remark": "%" + in.SearchContent + "%"})
	}

	res, err := l.svcCtx.PortfolioInvestmentsModel.FindAll(l.ctx, whereBuilder, "id ASC")
	if err != nil {
		return nil, errors.Wrapf(xerror.ErrorDB, "Get PortfolioInvestmentList err: %+v , in :%+v", err, in)
	}
	var resp []*pb.PortfolioInvestmentDetail
	if len(res) > 0 {
		for _, v := range res {
			portfolioInvestmentDetail := new(pb.PortfolioInvestmentDetail)
			itemWhereBuilder := l.svcCtx.PortfolioInvestmentItemsModel.RowBuilder()
			itemWhereBuilder = itemWhereBuilder.Where(squirrel.Eq{"gid": v.Id})
			itemRsp, err := l.svcCtx.PortfolioInvestmentItemsModel.FindAll(l.ctx, itemWhereBuilder, "investment ASC")
			if err != nil {
				return nil, errors.Wrapf(xerror.ErrorDB, "Get PortfolioInvestmentList err: %+v , in :%+v", err, in)
			}
			_ = copier.Copy(portfolioInvestmentDetail, v)
			var investments []string
			if len(itemRsp) > 0 {
				for _, k := range itemRsp {
					investments = append(investments, k.Investment)
				}
				portfolioInvestmentDetail.Investments = investments
				portfolioInvestmentDetail.Created = v.Created.Format("2006-01-02 15:04:05")
				portfolioInvestmentDetail.Updated = v.Updated.Format("2006-01-02 15:04:05")
				//todo  Operator
			}
			resp = append(resp, portfolioInvestmentDetail)
		}
	}

	return &pb.PortfolioInvestmentListRsp{
		Details: resp,
	}, nil
}
