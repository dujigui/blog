package visitor

import (
	"encoding/json"
	"errors"
	"fmt"
	. "github.com/dujigui/blog/services"
	. "github.com/dujigui/blog/services/logs"
	. "github.com/dujigui/blog/utils"
	"github.com/kataras/iris/v12"
	"io/ioutil"
	"net/http"
	"net/url"
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
		ctx.WriteString("Authorization Code 不能为空")
		return
	}
	Logger().Trace("qq", "获取到 Authorize Code", Params{"authorizeCode": ac})

	ss := ctx.URLParam("state")
	state, err := DecodeState(ss)
	if err != nil {
		Logger().Error("qq", "State 检验失败", Params{"state": ss}.Err(err))
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("State 检验失败")
		return
	}

	at, err := accessToken(ac)
	if err != nil {
		Logger().Error("qq", "无法获取 Access Token", Params{"authorizeCode": ac}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("无法获取 Access Token")
		return
	}
	Logger().Trace("qq", "获取到 Access Token", Params{"authorizeCode": ac, "accessToken": at})

	oi, err := openID(at)
	if err != nil {
		Logger().Error("qq", "无法获取 open ID", Params{"authorizeCode": ac, "accessToken": at}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("无法获取 open ID")
		return
	}
	Logger().Trace("qq", "获取到 open ID", Params{"authorizeCode": ac, "accessToken": at})

	ui, err := userInfo(at, oi)
	if err != nil {
		Logger().Error("qq", "无法获取用户信息", Params{"authorizeCode": ac, "accessToken": at, "openID": oi}.Err(err))
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("无法获取 open ID")
		return
	}
	Logger().Trace("qq", "获取到用户信息", Params{"authorizeCode": ac, "accessToken": at, "openID": oi})

	d, err := json.MarshalIndent(ui, "", "  ")
	fmt.Println(d, err)
	ctx.Redirect(state.Redirect, iris.StatusFound)
}

func accessToken(ac string) (string, error) {
	resp, err := http.Get(fmt.Sprintf(at, Pref().QQAppID, Pref().QQKey, ac, Pref().QQRedirect))
	if err != nil || resp.StatusCode != iris.StatusFound {
		return "", err
	}

	referer := resp.Header.Get("Referer")
	if referer == "" {
		return "", errors.New("无法获取 Referer")
	}

	s, err := url.Parse(referer)
	if err != nil {
		return "", err
	}

	at := s.Query().Get("access_token")
	if ac == "" {
		return "", errors.New("无法从 Referer 获取 Access Token")
	}

	return at, nil
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

type qqInfo struct {
	NickName string `json:"nickname"`
	Avatar   string `json:"pic"`
}

func userInfo(accessToken, openID string) (qqInfo, error) {
	var qi qqInfo

	resp, err := http.Get(fmt.Sprintf(info, accessToken, Pref().QQAppID, openID))
	if err != nil || resp.StatusCode != iris.StatusOK {
		return qi, err
	}
	defer resp.Body.Close()

	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return qi, err
	}

	if err := json.Unmarshal(d, &qi); err != nil {
		return qi, err
	}

	return qi, nil
}
