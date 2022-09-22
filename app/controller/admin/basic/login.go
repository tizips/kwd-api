package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"golang.org/x/crypto/bcrypt"
	"kwd/app/constant"
	"kwd/app/model"
	"kwd/app/request/admin/basic"
	res "kwd/app/response/admin/basic"
	helperService "kwd/app/service/helper"
	"kwd/kernel/app"
	"kwd/kernel/authorize"
	"kwd/kernel/response"
	"strconv"
)

func DoLoginByAccount(ctx *gin.Context) {

	var request basic.DoLoginByAccess

	if err := ctx.ShouldBind(&request); err != nil {
		response.FailByRequest(ctx, err)
		return
	}

	var admin model.SysAdmin

	app.Database.Find(&admin, "`username`=? and `is_enable`=?", request.Username, constant.IsEnableYes)

	if admin.Id <= 0 {
		response.Fail(ctx, "用户名或密码错误")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(request.Password)); err != nil {
		response.Fail(ctx, "用户名或密码错误")
		return
	}

	now := carbon.Now()

	claims := jwt.RegisteredClaims{
		Issuer:    constant.ContextAdmin,
		Subject:   strconv.Itoa(admin.Id),
		NotBefore: jwt.NewNumericDate(now.Carbon2Time()),
		IssuedAt:  jwt.NewNumericDate(now.Carbon2Time()),
		ExpiresAt: jwt.NewNumericDate(now.AddHours(app.Cfg.Jwt.Lifetime).Carbon2Time()),
		ID:        helperService.JwtToken(admin.Id),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(app.Cfg.Jwt.Secret))

	if err != nil {
		response.Fail(ctx, "用户名或密码错误")
		return
	}

	responses := res.DoLoginByAccess{
		Token:    signed,
		ExpireAt: now.AddHours(app.Cfg.Jwt.Lifetime).Timestamp(),
	}

	response.Success(ctx, responses)

}

func DoLoginByQrcode(ctx *gin.Context) {

}

func DoLogout(ctx *gin.Context) {

	claims := authorize.Jwt(ctx)

	ok := helperService.BlackJwt(ctx, constant.ContextAdmin, *claims)

	if !ok {
		response.Fail(ctx, "退出失败，请稍后重试！")
		return
	}

	response.Success[any](ctx)
}
