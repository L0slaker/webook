package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/l0slakers/webook/internal/domain"
	"github.com/l0slakers/webook/internal/service"
	"net/http"
)

const (
	emailRegexPattern    = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,30}$`
)

const (
	systemErr     = "系统错误！"
	emailMatchErr = "邮箱格式不正确！"
	pwdMatchErr   = "密码格式不正确：长度在8~30个字符，至少包含1个数字，一个字母和一个特殊字符"
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
		ctx.String(http.StatusInternalServerError, systemErr)
		return
	}
	if !isEmailMatch {
		ctx.String(http.StatusBadRequest, emailMatchErr)
		return
	}
	isPwdMatch, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, systemErr)
		return
	}
	if !isPwdMatch {
		ctx.String(http.StatusBadRequest, pwdMatchErr)
		return
	}

	// TODO 其他校验（例如数字长度、格式等）

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
		ctx.String(http.StatusInternalServerError, systemErr)
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

	// TODO 其他校验（例如数字长度、格式等）

	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		// 设置登录态
		sess := sessions.Default(ctx)
		sess.Set("userId", u.ID)
		// 设置登陆时间（以实际业务为准）
		sess.Options(sessions.Options{
			MaxAge: 900, // 十五分钟
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusInternalServerError, systemErr)
		}
		ctx.String(http.StatusOK, "登陆成功！")
	case service.ErrUnknownEmail:
		ctx.String(http.StatusOK, service.ErrUnknownEmail.Error())
	case service.ErrWrongInfo:
		ctx.String(http.StatusOK, service.ErrWrongInfo.Error())
	default:
		ctx.String(http.StatusInternalServerError, systemErr)
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) Info(ctx *gin.Context) {
	ctx.String(http.StatusOK, "用户个人信息：")
}
