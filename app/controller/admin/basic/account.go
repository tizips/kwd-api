package basic

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"kwd/app/constant"
	"kwd/app/model"
	"kwd/app/request/admin/basic"
	res "kwd/app/response/admin/basic"
	"kwd/kernel/app"
	"kwd/kernel/authorize"
	"kwd/kernel/response"
)

func ToAccountByInformation(ctx *gin.Context) {

	admin := authorize.Admin(ctx)

	if admin.Id <= 0 {
		response.FailByLogin(ctx)
		return
	}

	responses := res.ToAccountByInformation{
		Username: admin.Username,
		Nickname: admin.Nickname,
		Avatar:   admin.Avatar,
		Mobile:   admin.Mobile,
	}

	response.Success(ctx, responses)
}

func ToAccountByModule(ctx *gin.Context) {

	responses := make([]res.ToAccountByModule, 0)

	var modules []model.SysModule

	tx := app.Database.
		Where("`is_enable` = ?", constant.IsEnableYes)

	tc := app.Database.
		Select("1").
		Model(model.SysPermission{}).
		Where(fmt.Sprintf("`%s`.`id`=`%s`.`module_id`", model.TableSysModule, model.TableSysPermission))

	if !authorize.Root(authorize.Id(ctx)) {
		tc = tc.
			Joins(fmt.Sprintf("left join `%s` on `%s`.`id`=`%s`.`permission_id`", model.TableSysRoleBindPermission, model.TableSysPermission, model.TableSysRoleBindPermission)).
			Joins(fmt.Sprintf("left join `%s` on `%s`.`role_id`=`%s`.`role_id` and `%s`.`admin_id`=?", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole, model.TableSysAdminBindRole), authorize.Id(ctx)).
			Where(fmt.Sprintf("`%s`.`id` is not null and `%s`.`deleted_at` is null and `%s`.`deleted_at` is null", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole))
	}

	tx.
		Where("exists (?)", tc).
		Order("`order` asc").
		Find(&modules)

	for _, item := range modules {
		responses = append(responses, res.ToAccountByModule{
			Id:   item.Id,
			Slug: item.Slug,
			Name: item.Name,
		})
	}

	response.Success(ctx, responses)
}

func ToAccountByPermission(ctx *gin.Context) {

	var request basic.ToAccountByPermission

	if err := ctx.ShouldBindQuery(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var responses = make([]string, 0)

	var permissions []model.SysPermission

	tx := app.Database.
		Where("`module_id` = ? and `method` <> ? and `path` <> ?", request.Module, "", "")

	if !authorize.Root(authorize.Id(ctx)) {
		tx = tx.
			Joins(fmt.Sprintf("left join `%s` on `%s`.`id`=`%s`.`permission_id`", model.TableSysRoleBindPermission, model.TableSysPermission, model.TableSysRoleBindPermission)).
			Joins(fmt.Sprintf("left join `%s` on `%s`.`role_id`=`%s`.`role_id` and `%s`.`admin_id`=?", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole, model.TableSysAdminBindRole), authorize.Id(ctx)).
			Where(fmt.Sprintf("`%s`.`id` is not null and `%s`.`deleted_at` is null and `%s`.`deleted_at` is null", model.TableSysAdminBindRole, model.TableSysRoleBindPermission, model.TableSysAdminBindRole))
	}

	tx.
		Find(&permissions)

	for _, item := range permissions {
		responses = append(responses, item.Slug)
	}

	response.Success(ctx, responses)
}

func DoAccountByUpdate(ctx *gin.Context) {

	var request basic.DoAccountByUpdate

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	admin := authorize.Admin(ctx)

	admin.Avatar = request.Avatar

	if request.Password != "" {

		password, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		admin.Password = string(password)
	}

	if t := app.Database.Save(&admin); t.RowsAffected <= 0 {
		response.Fail(ctx, "修改失败")
		return
	}

	response.Success[any](ctx)
}
