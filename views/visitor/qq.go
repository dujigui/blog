package visitor

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/services/users"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	at = "https://graph.qq.com/oauth2.0/token?grant_type=authorization_code" +
		"&client_id=%s" +
		"&client_secret=%s" +
		"&code=%s" +
		"&redirect_uri=%s"
	me   = "https://graph.qq.com/oauth2.0/me?access_token=%s"
	info = "https://graph.qq.com/user/get_user_info" +
		"?access_token=%s" +
		"&oauth_consumer_key=%s" +
		"&openid=%s"
)

func qq(ctx iris.Context) {
	ac := ctx.URLParam("code")
	if ac == "" {
		Logger().Error("qq", "无法获取 Authorize Code", Params{"authorizeCode": ac})
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "Authorization Code 不能为空", nil))
		return
	}
	Logger().Trace("qq", "获取到 Authorize Code", Params{"authorizeCode": ac})

	ss := ctx.URLParam("state")
	state, err := DecodeState(ss)
	if err != nil {
		Logger().Error("qq", "State 检验失败", Params{"state": ss}.Err(err))
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(Result(false, "State 检验失败", nil))
		return
	}

	at, err := accessToken(ac)
	if err != nil {
		Logger().Error("qq", "无法获取 Access Token", Params{"authorizeCode": ac}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "无法获取 Access Token", nil))
		return
	}
	Logger().Trace("qq", "获取到 Access Token", Params{"authorizeCode": ac, "accessToken": at})

	oi, err := openID(at)
	if err != nil {
		Logger().Error("qq", "无法获取 open ID", Params{"authorizeCode": ac, "accessToken": at}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "无法获取 open ID", nil))
		return
	}
	Logger().Trace("qq", "获取到 open ID", Params{"authorizeCode": ac, "accessToken": at, "openID": oi})

	ui, err := userInfo(at, oi)
	if err != nil {
		Logger().Error("qq", "无法获取用户信息", Params{"authorizeCode": ac, "accessToken": at, "openID": oi}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "无法获取用户信息", nil))
		return
	}

	qi, err := QQTable().ByOpenID(ui.OpenID)
	if err != nil {
		Logger().Error("qq", "无法获取数据库用户信息", Params{"authorizeCode": ac, "accessToken": at, "openID": oi}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(Result(false, "无法获取数据库用户信息", nil))
		return
	}

	p := Params{
		"nickname":       ui.Nickname,
		"gender":         ui.Gender,
		"gender_type":    ui.GenderType,
		"province":       ui.Province,
		"city":           ui.City,
		"year":           ui.Year,
		"constellation":  ui.Constellation,
		"figureurl":      ui.Figureurl,
		"figureurl_1":    ui.Figureurl1,
		"figureurl_2":    ui.Figureurl2,
		"figureurl_qq_1": ui.FigureurlQq1,
		"figureurl_qq_2": ui.FigureurlQq2,
		"figureurl_qq":   ui.FigureurlQq,
		"figureurl_type": ui.FigureurlType,
		"open_id":        ui.OpenID,
		"access_token":   ui.AccessToken,
	}
	if qi.ID == 0 {
		id, err := QQTable().Create(p)
		if err != nil || id == 0 {
			Logger().Error("qq", "无法创建 QQ 用户", p.Err(err))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(Result(false, "无法创建 QQ 用户", nil))
			return
		}
		id1, err1 := UserTable().Create(Params{"type": ViaQQ, "qq_id": id, "avatar": ui.FigureurlQq, "nickname": ui.Nickname})
		if err1 != nil || id1 == 0 {
			Logger().Error("qq", "无法关联 QQ 用户", p.Err(err))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(Result(false, "无法关联 QQ 用户", nil))
			return
		}
		ctx.SetCookieKV("token", CreateToken(id1, false, CookieExpire), iris.CookieExpires(CookieExpire))
		ctx.Redirect(state.Redirect, iris.StatusFound)
	} else {
		if err := QQTable().Update(qi.ID, p); err != nil {
			Logger().Error("qq", "更新数据库失败", p.Err(err))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(Result(false, "更新数据库失败", nil))
			return
		}
		ui, err := UserTable().ByQQID(qi.ID)
		if err != nil {
			Logger().Error("qq", "此 QQ 号码未关联账号", p.Err(err))
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(Result(false, "此 QQ 号码未关联账号", nil))
			return
		}
		ctx.SetCookieKV("token", CreateToken(ui.ID, ui.Admin, CookieExpire), iris.CookieExpires(CookieExpire))
		ctx.Redirect(state.Redirect, iris.StatusFound)
	}
}

func accessToken(ac string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(at, Pref().QQAppID, Pref().QQAppKey, ac, Pref().QQRedirect))
	if err != nil {
		return "", err
	}

	// access_token=A91DA9A03B5F5003637DB3B1FD84D81C&expires_in=7776000&refresh_token=3EF3D730E9D42D4B61AB7D960092A85F
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := string(d)
	start := strings.Index(body, "access_token=")
	if start != -1 {
		start += len("access_token=")
		if len(body) > start {
			end := strings.Index(body[start:], "&")
			if end != -1 && len(body) > start+end {
				return body[start : start+end], nil
			}
		}
	}
	return "", errors.New("无法解析 body")
}

func openID(accessToken string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(me, accessToken))
	if err != nil || resp.StatusCode != iris.StatusOK {
		return "", err
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(d)
	if strings.Contains(s, "callback") {
		li := strings.Index(s, "(")
		ri := strings.Index(s, ")")
		s = s[li+1 : ri]
	}

	type Entity struct {
		ClientID string `json:"client_id"`
		OpenID   string `json:"openid"`
	}

	var e Entity
	if err := json.Unmarshal([]byte(s), &e); err != nil {
		return "", errors.New("无法解析 body")
	}

	return e.OpenID, nil
}

func userInfo(accessToken, openID string) (QQInfo, error) {
	var qu QQInfo

	resp, err := http.Get(fmt.Sprintf(info, accessToken, Pref().QQAppID, openID))
	if err != nil || resp.StatusCode != iris.StatusOK {
		return qu, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&qu); err != nil {
		return qu, err
	}
	if qu.Ret != 0 {
		return qu, errors.New(qu.Msg)
	}

	qu.OpenID = openID
	qu.AccessToken = accessToken
	return qu, nil
}
