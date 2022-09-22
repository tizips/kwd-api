package basic

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/golang-module/carbon/v2"
	"kwd/app/constant"
	"kwd/app/service/helper"
	"kwd/kernel/app"
)

func JwtParseMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		if authorization := ctx.GetHeader(constant.JwtAuthorization); authorization != "" {

			var claims jwt.RegisteredClaims

			_, err := jwt.ParseWithClaims(authorization, &claims, func(token *jwt.Token) (any, error) {
				return []byte(app.Cfg.Jwt.Secret), nil
			})

			valid, _ := err.(*jwt.ValidationError)

			if (err == nil || valid.Errors > 0 && valid.Is(jwt.ErrTokenExpired)) && claims.Subject != "" {
				now := carbon.Now()
				if ok := claims.VerifyNotBefore(now.Carbon2Time(), false); ok {
					if ok = claims.VerifyExpiresAt(now.Carbon2Time(), true); ok { //	生效的授权令牌操作
						set(ctx, claims)
					} else if ok = claims.VerifyExpiresAt(now.SubHours(app.Cfg.Jwt.Lifetime).Carbon2Time(), true); ok { //	失效的授权令牌，重新发放
						refresh(ctx, claims)
					}
				}
			}
		}

		ctx.Next()
	}
}

func set(ctx *gin.Context, claims jwt.RegisteredClaims) {
	ctx.Set(constant.ContextID, claims.Subject)
	ctx.Set(constant.ContextJWT, claims)
}

func refresh(ctx *gin.Context, claims jwt.RegisteredClaims) {

	cache, _ := app.Redis.HGetAll(ctx, helper.Blacklist(constant.ContextAdmin, "refresh", claims.ID)).Result()

	now := carbon.Now()

	if len(cache) <= 0 {

		expires := claims.ExpiresAt
		id := claims.ID

		claims.NotBefore = jwt.NewNumericDate(now.Carbon2Time())
		claims.IssuedAt = jwt.NewNumericDate(now.Carbon2Time())
		claims.ExpiresAt = jwt.NewNumericDate(now.AddHours(app.Cfg.Jwt.Lifetime).Carbon2Time())
		claims.ID = helper.JwtToken(claims.Subject)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		if signed, err := token.SignedString([]byte(app.Cfg.Jwt.Secret)); err == nil {

			set(ctx, claims)

			ctx.Header(constant.JwtAuthorization, signed)

			if affected, err := app.Redis.HSet(ctx, helper.Blacklist(constant.ContextAdmin, "refresh", id), "token", signed, "created_at", now.ToDateTimeString()).Result(); err == nil && affected > 0 {
				app.Redis.ExpireAt(ctx, helper.Blacklist(constant.ContextAdmin, "refresh", id), carbon.Time2Carbon(expires.Time).AddHours(app.Cfg.Jwt.Lifetime).Carbon2Time())
			}
		}

	} else {

		diff := now.DiffAbsInSeconds(carbon.Parse(cache["created_at"]))

		if diff <= app.Cfg.Jwt.Leeway {
			ctx.Header(constant.JwtAuthorization, cache["token"])
		}
	}
}
