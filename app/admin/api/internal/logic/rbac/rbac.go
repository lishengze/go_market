package rbac

import (
	"context"
	"market_server/app/admin/model"
	"market_server/common/utils"

	"github.com/Masterminds/squirrel"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type RbacLogic struct {
	logx.Logger
	ctx                context.Context
	conn               sqlx.SqlConn
	MenuModel          model.MenuModel
	RoleOperationModel model.RoleOperationModel
	SecurityLogModel   model.SecurityLogModel
	UserRoleModel      model.UserRoleModel
}

func NewRbacLogic(ctx context.Context, conn sqlx.SqlConn) *RbacLogic {
	return &RbacLogic{
		Logger:             logx.WithContext(ctx),
		ctx:                ctx,
		conn:               conn,
		MenuModel:          model.NewMenuModel(conn),
		RoleOperationModel: model.NewRoleOperationModel(conn),
		SecurityLogModel:   model.NewSecurityLogModel(conn),
		UserRoleModel:      model.NewUserRoleModel(conn),
	}
}

func (m *RbacLogic) FindUserRoleByUserId(userId int64) (menus []*model.UserRole, err error) {
	rowBuilder := m.UserRoleModel.RowBuilder().
		Where(squirrel.Eq{"user_id": userId})
	return m.UserRoleModel.FindAll(m.ctx, rowBuilder, "")
}

func (m *RbacLogic) FindUserOperationByRoleId(roleIds []int64) (menus []*model.RoleOperation, err error) {
	rowBuilder := m.RoleOperationModel.RowBuilder().
		PlaceholderFormat(squirrel.Question).
		Where(squirrel.Eq{"role_id": roleIds})

	return m.RoleOperationModel.FindAll(m.ctx, rowBuilder, "")
}

func (m *RbacLogic) FindUserOperations(userId int64) (menus []*model.Menu, err error) {
	userRoles, err := m.FindUserRoleByUserId(userId)
	if err != nil {
		return
	}
	var roleIdList = make([]int64, 0, len(userRoles))
	for _, userRole := range userRoles {
		roleIdList = append(roleIdList, userRole.RoleId)
	}

	var roleOperations []*model.RoleOperation
	if len(roleIdList) > 0 {
		roleOperations, err = m.FindUserOperationByRoleId(roleIdList)
		if err != nil {
			return
		}
	}

	if len(roleOperations) == 0 {
		roleOperations = make([]*model.RoleOperation, 0)
	}

	var operationIdList = make([]int64, 0, len(roleOperations))
	for _, roleOperation := range roleOperations {
		operationIdList = append(operationIdList, roleOperation.OperationId)
	}
	operationIdList = utils.UniqueInt64(operationIdList)

	if len(roleIdList) > 0 {
		menus, err = m.MenuModel.FindAll(m.ctx, m.MenuModel.RowBuilder().Where(squirrel.Eq{"id": operationIdList}), "")
		if err != nil {
			return
		}
	} else {
		menus = make([]*model.Menu, 0)
	}
	return
}
