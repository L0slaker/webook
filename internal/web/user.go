package web

import (
	"net/http"
	"time"

	"github.com/l0slakers/webook/internal/domain"
	pkgTime "github.com/l0slakers/webook/internal/pkg/time"
	"github.com/l0slakers/webook/internal/service"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// TODO key从配置文件中读取
const JwtKey = "U2FsdGVkX19IDF17ov2HRI/9TlXkROBA"

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,30}$`
)

const (
	systemErr            = "系统错误！"
	emailMatchErr        = "邮箱格式不正确！"
	pwdMatchErr          = "密码格式不正确：长度在8~30个字符，至少包含1个数字，一个字母和一个特殊字符"
	nicknameMatchErr     = "昵称长度不正确，需在3-15之间"
	birthdayMatchErr     = "时间格式不正确，参考（1900-01-01）"
	introductionMatchErr = "个人简介长度不正确，需在100以下"
)

type UserHandler struct {
	svc            *service.UserService
	passwordRexExp *regexp.Regexp
	emailRexExp    *regexp.Regexp
}

func NewUserHandler(userSvc *service.UserService) *UserHandler {
	return &UserHandler{
		svc:            userSvc,
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
	}
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type UserInfo struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req UserInfo
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}

	// 正则校验
	isEmailMatch, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isEmailMatch {
		ctx.String(http.StatusOK, emailMatchErr)
		return
	}
	isPwdMatch, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isPwdMatch {
		ctx.String(http.StatusOK, pwdMatchErr)
		return
	}

	err = h.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功！")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, service.ErrDuplicateEmail.Error())
	default:
		ctx.String(http.StatusOK, systemErr)
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type LoginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginInfo
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}

	// 正则校验
	isEmailMatch, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isEmailMatch {
		ctx.String(http.StatusOK, emailMatchErr)
		return
	}
	isPwdMatch, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isPwdMatch {
		ctx.String(http.StatusOK, pwdMatchErr)
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 设置登录态
		claim := UserClaim{
			RegisteredClaims: jwt.RegisteredClaims{
				// 十五分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
			},
			UserID: u.ID,
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
		tokenStr, err := token.SignedString([]byte(JwtKey))
		if err != nil {
			ctx.String(http.StatusOK, systemErr)
			return
		}
		// 写入头部返回给前端
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登陆成功:"+tokenStr)
	case service.ErrUnknownEmail:
		ctx.String(http.StatusOK, service.ErrUnknownEmail.Error())
	case service.ErrWrongInfo:
		ctx.String(http.StatusOK, service.ErrWrongInfo.Error())
	default:
		ctx.String(http.StatusOK, systemErr)
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type LoginInfo struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginInfo
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}

	// 正则校验
	isEmailMatch, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isEmailMatch {
		ctx.String(http.StatusOK, emailMatchErr)
		return
	}
	isPwdMatch, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	if !isPwdMatch {
		ctx.String(http.StatusOK, pwdMatchErr)
		return
	}

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 设置登录态
		sess := sessions.Default(ctx)
		sess.Set("userId", u.ID)
		// 设置登陆时间（以实际业务为准）
		// TODO 通过配置文件来设置
		sess.Options(sessions.Options{
			MaxAge: 900, // 单位:秒
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, systemErr)
			return
		}
		ctx.String(http.StatusOK, "登陆成功！")
	case service.ErrUnknownEmail:
		ctx.String(http.StatusOK, service.ErrUnknownEmail.Error())
	case service.ErrWrongInfo:
		ctx.String(http.StatusOK, service.ErrWrongInfo.Error())
	default:
		ctx.String(http.StatusOK, systemErr)
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type UserInfo struct {
		Nickname     string `json:"nickname"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	var req UserInfo
	if err := ctx.ShouldBind(&req); err != nil {
		return
	}

	// 校验
	if req.Nickname != "" {
		if len(req.Nickname) < 3 || len(req.Nickname) > 15 {
			ctx.String(http.StatusOK, nicknameMatchErr)
			return
		}
	}
	if req.Birthday != "" {
		if !pkgTime.IsValidDate(pkgTime.YYYYMMDD, req.Birthday) {
			ctx.String(http.StatusOK, birthdayMatchErr)
			return
		}
	}
	if req.Introduction != "" {
		if len(req.Introduction) > 100 {
			ctx.String(http.StatusOK, introductionMatchErr)
			return
		}
	}

	// 使用session存储
	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)

	val, ok := ctx.Get("user")
	if !ok {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	claim := val.(UserClaim)

	err := h.svc.Edit(ctx, domain.User{
		ID:           claim.UserID,
		Nickname:     req.Nickname,
		Birthday:     req.Birthday,
		Introduction: req.Introduction,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "修改成功！")
	default:
		ctx.String(http.StatusOK, systemErr)
	}
}

func (h *UserHandler) Info(ctx *gin.Context) {
	// 使用session存储
	//sess := sessions.Default(ctx)
	//uid := sess.Get("userId").(int64)

	val, ok := ctx.Get("user")
	if !ok {
		ctx.String(http.StatusOK, systemErr)
		return
	}
	claim := val.(UserClaim)

	u, err := h.svc.Info(ctx, claim.UserID)
	if err != nil {
		ctx.String(http.StatusOK, systemErr)
		return
	}

	type User struct {
		Nickname     string `json:"nickname"`
		Email        string `json:"email"`
		Birthday     string `json:"birthday"`
		Introduction string `json:"introduction"`
	}
	ctx.JSON(http.StatusOK, User{
		Nickname:     u.Nickname,
		Email:        u.Email,
		Birthday:     u.Birthday,
		Introduction: u.Introduction,
	})
}

type UserClaim struct {
	jwt.RegisteredClaims
	UserID int64
}
