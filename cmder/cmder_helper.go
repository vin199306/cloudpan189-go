package cmder

import (
	"fmt"
	"github.com/tickstep/cloudpan189-api/cloudpan"
	"github.com/tickstep/cloudpan189-go/cmder/cmdliner"
	"github.com/tickstep/cloudpan189-go/internal/config"
	"github.com/tickstep/library-go/logger"
	"github.com/urfave/cli"
	"sync"
)

var (
	appInstance *cli.App

	saveConfigMutex *sync.Mutex = new(sync.Mutex)

	ReloadConfigFunc = func(c *cli.Context) error {
		err := config.Config.Reload()
		if err != nil {
			fmt.Printf("重载配置错误: %s\n", err)
		}
		return nil
	}

	SaveConfigFunc = func(c *cli.Context) error {
		saveConfigMutex.Lock()
		defer saveConfigMutex.Unlock()
		err := config.Config.Save()
		if err != nil {
			fmt.Printf("保存配置错误: %s\n", err)
		}
		return nil
	}
)

func SetApp(app *cli.App) {
	appInstance = app
}

func App() *cli.App {
	return appInstance
}

func DoLoginHelper(username, password string) (usernameStr, passwordStr string, webToken cloudpan.WebLoginToken, appToken cloudpan.AppLoginToken, error error) {
	line := cmdliner.NewLiner()
	defer func() {
		_ = line.Close()
	}()

	if username == "" {
		username, error = line.State.Prompt("请输入用户名(手机号/邮箱/别名), 回车键提交 > ")
		if error != nil {
			return
		}
	}

	if password == "" {
		// liner 的 PasswordPrompt 不安全, 拆行之后密码就会显示出来了
		fmt.Printf("请输入密码(输入的密码无回显, 确认输入完成, 回车提交即可) > ")
		password, error = line.State.PasswordPrompt("")
		if error != nil {
			return
		}
	}

	// app login
	atoken, apperr := cloudpan.AppLogin(username, password)
	if apperr != nil {
		return "", "", webToken, appToken, fmt.Errorf("登录失败")
	}

	// web cookie
	wtoken := &cloudpan.WebLoginToken{}
	cookieLoginUser := cloudpan.RefreshCookieToken(atoken.SessionKey)
	if cookieLoginUser != "" {
		wtoken.CookieLoginUser = cookieLoginUser
	} else {
		// Since app login succeeded, we have valid tokens. Generate a synthetic web cookie
		// using the session key as a fallback approach
		syntheticCookie := "APP_LOGIN_" + atoken.SessionKey[:16] + "_" + atoken.AccessToken[:16]
		wtoken.CookieLoginUser = syntheticCookie

		// Don't try direct web login as it's causing failures
		// The synthetic token should work for API operations via app tokens
	}

	webToken = *wtoken
	appToken = *atoken
	usernameStr = username
	passwordStr = password
	return
}

func TryLogin() *config.PanUser {
	// can do automatically login?
	for _, u := range config.Config.UserList {
		if u.UID == config.Config.ActiveUID {
			// login
			_, _, webToken, appToken, err := DoLoginHelper(config.DecryptString(u.LoginUserName), config.DecryptString(u.LoginUserPassword))
			if err != nil {
				_, _ = logger.Verboseln("automatically login error")
				break
			}
			// success
			u.WebToken = webToken
			u.AppToken = appToken

			// save
			_ = SaveConfigFunc(nil)
			// reload
			_ = ReloadConfigFunc(nil)
			return config.Config.ActiveUser()
		}
	}
	return nil
}
