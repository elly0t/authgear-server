package userverify

import (
	"fmt"

	"github.com/skygeario/skygear-server/pkg/auth/dependency/provider/password"
	"github.com/skygeario/skygear-server/pkg/core/auth/authinfo"
	"github.com/skygeario/skygear-server/pkg/core/config"
)

type AutoUpdateUserVerifyFunc func(*authinfo.AuthInfo, []*password.Principal)

func CreateAutoUpdateUserVerifyfunc(tConfig config.TenantConfiguration) AutoUpdateUserVerifyFunc {
	if !tConfig.UserConfig.UserVerification.AutoUpdate {
		return nil
	}

	switch tConfig.UserConfig.UserVerification.Criteria {
	case "all":
		return func(authInfo *authinfo.AuthInfo, principals []*password.Principal) {
			allVerified := true
			for _, principal := range principals {
				for _, keyConfig := range tConfig.UserConfig.UserVerification.Keys {
					if principal.LoginIDKey == keyConfig.Key && !authInfo.VerifyInfo[principal.LoginID] {
						allVerified = false
						break
					}
				}
			}

			authInfo.Verified = allVerified
		}
	case "any":
		return func(authInfo *authinfo.AuthInfo, principals []*password.Principal) {
			for _, principal := range principals {
				for _, keyConfig := range tConfig.UserConfig.UserVerification.Keys {
					if principal.LoginIDKey == keyConfig.Key && authInfo.VerifyInfo[principal.LoginID] {
						authInfo.Verified = true
						return
					}
				}
			}

			authInfo.Verified = false
		}
	default:
		panic(fmt.Errorf("unexpected verify criteria `%s`", tConfig.UserConfig.UserVerification.Criteria))
	}
}
