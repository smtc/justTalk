package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/bitly/oauth2_proxy/cookie"
	"github.com/gin-gonic/gin"
	"github.com/guotie/config"
	"github.com/guotie/deferinit"
	"github.com/satori/go.uuid"
	"github.com/smtc/justTalk/providers"
)

var (
	_              = fmt.Printf
	_sid           string
	_cookieSecure  bool
	_cookieSeed    string
	_expireSecond  int
	_cookieExpire  time.Duration
	githubProvider providers.Provider
	githubCB       = "/oauth/callback/github"
)

type githubProfile struct {
	AccessToken  string
	ExpiresOn    time.Time
	RefreshToken string
	ProvideId    string `json:"email"`
	Name         string `json:"name"`
	AvatarUrl    string `json:"avatar_url"`
	Blog         string `json:"blog"`
	Location     string `json:"location"`
}

type cachedUserData struct {
	Id      string
	Profile interface{}
}

func init() {
	deferinit.AddInit(func() {
		githubProvider = providers.NewGitHubProvider(&providers.ProviderData{
			LoginUrl:     &url.URL{},
			RedeemUrl:    &url.URL{},
			ProfileUrl:   &url.URL{},
			ValidateUrl:  &url.URL{},
			ClientID:     config.GetStringDefault("clientId", "66a252c0d27dc279b7cb"),
			ClientSecret: config.GetStringDefault("clientSecret", "a313e648feff0e6b30794142ff9304e42cd50da1"),
		})
		_sid = config.GetStringDefault("cookieName", "sid")
		_cookieSecure = config.GetBooleanDefault("cookieSecure", false)
		_cookieSeed = config.GetStringDefault("cookieSeed", "cookieseed")
		_expireSecond = config.GetIntDefault("cookieExpire", 86400*30)
		_cookieExpire = time.Duration(_expireSecond) * time.Second
	}, nil, 0)
}

func oauthLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{
		"title":        "登陆",
		"github_login": githubProvider.GetLoginURL(hostname+githubCB, ""),
	})
}

// github授权后的页面，根据参数code获取用户
func githubLogin(c *gin.Context) {
	code := c.Query("code")
	sess, err := githubProvider.Redeem("", code)
	if err != nil {
		c.String(500, err.Error())
		return
	}
	profile, err := getGithubProfile(sess)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// 创建新用户，或者授权登陆
	exist, user, err := providerUserHasExist("github", profile.ProvideId)
	if err != nil {
		c.String(500, err.Error())
		return
	}

	if exist == false {
		// todo: 需要更加精细的控制：
		//     检查用户名是否有重合
		//     跳转到一个设置用户名的页面
		// 新用户
		user = &User{
			Name:         profile.Name,
			Email:        profile.ProvideId,
			Blog:         profile.Blog,
			City:         profile.Location,
			Avatar:       profile.AvatarUrl,
			ProviderName: "github",
			ProviderId:   profile.ProvideId,
		}
		err = createProviderUser(user)
		if err != nil {
			c.String(500, err.Error())
			return
		}
	}

	next := c.Query("next")
	if next == "" {
		next = "/"
	}

	ctoken, err := setUserCache(&cachedUserData{user.Id, profile})
	if err != nil {
		c.String(500, err.Error())
		return
	}

	// 设置cookie
	setCookie(c.Writer, c.Request, ctoken)

	c.String(200, ctoken)
}

func getGithubProfile(sess *providers.SessionState) (p *githubProfile, err error) {
	var (
		body    []byte
		profile githubProfile
	)
	// 从github接口得到用户的profile, 反序列化
	body, err = githubProvider.GetProfile(sess)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &profile)
	if err != nil {
		return
	}

	profile.AccessToken = sess.AccessToken
	profile.RefreshToken = sess.RefreshToken
	profile.ExpiresOn = sess.ExpiresOn

	return &profile, nil
}

// 创建一个随机数，以该随机数为key，用户id和用户access token作为value，存在redis中，并返回key
func setUserCache(data *cachedUserData) (string, error) {
	ctoken := uuid.NewV4().String()
	buf, err := json.Marshal(data)
	if err != nil {
		return "", err
	}
	err = SETEX(ctoken, _expireSecond, buf)
	if err != nil {
		return "", err
	}
	return ctoken, nil
}

// 生成cookie
func makeCookie(req *http.Request, value string, expiration time.Duration, now time.Time) *http.Cookie {
	domain := req.Host
	if h, _, err := net.SplitHostPort(domain); err == nil {
		domain = h
	}

	if value != "" {
		value = cookie.SignedValue(_cookieSeed, _sid, value, now)
	}

	return &http.Cookie{
		Name:     _sid,
		Value:    value,
		Path:     "/",
		Domain:   domain,
		HttpOnly: true,
		Secure:   _cookieSecure,
		Expires:  now.Add(expiration),
	}
}

// 清除cookie
func clearCookie(rw http.ResponseWriter, req *http.Request) {
	http.SetCookie(rw, makeCookie(req, "", time.Hour*-1, time.Now()))
}

// 设置cookie
func setCookie(rw http.ResponseWriter, req *http.Request, val string) {
	http.SetCookie(rw, makeCookie(req, val, _cookieExpire, time.Now()))
}
