package mfa

import (
	"testing"
	"time"

	coreAuth "github.com/skygeario/skygear-server/pkg/core/auth"
	"github.com/skygeario/skygear-server/pkg/core/config"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMaskAuthenticators(t *testing.T) {
	date := time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC)
	Convey("MaskAuthenticators", t, func() {
		input := []interface{}{
			TOTPAuthenticator{
				ID:          "totp",
				Type:        coreAuth.AuthenticatorTypeTOTP,
				CreatedAt:   date,
				ActivatedAt: &date,
				DisplayName: "totp",
			},
			OOBAuthenticator{
				ID:          "oobsms",
				Type:        coreAuth.AuthenticatorTypeOOB,
				CreatedAt:   date,
				ActivatedAt: &date,
				Channel:     coreAuth.AuthenticatorOOBChannelSMS,
				Phone:       "+85298765432",
			},
			OOBAuthenticator{
				ID:          "oobemail",
				Type:        coreAuth.AuthenticatorTypeOOB,
				CreatedAt:   date,
				ActivatedAt: &date,
				Channel:     coreAuth.AuthenticatorOOBChannelEmail,
				Email:       "johndoe@example.com",
			},
		}
		actual := MaskAuthenticators(input)
		expected := []interface{}{
			MaskedTOTPAuthenticator{
				ID:          "totp",
				Type:        coreAuth.AuthenticatorTypeTOTP,
				CreatedAt:   date,
				ActivatedAt: &date,
				DisplayName: "totp",
			},
			MaskedOOBAuthenticator{
				ID:          "oobsms",
				Type:        coreAuth.AuthenticatorTypeOOB,
				CreatedAt:   date,
				ActivatedAt: &date,
				Channel:     coreAuth.AuthenticatorOOBChannelSMS,
				MaskedPhone: "+8529876****",
			},
			MaskedOOBAuthenticator{
				ID:          "oobemail",
				Type:        coreAuth.AuthenticatorTypeOOB,
				CreatedAt:   date,
				ActivatedAt: &date,
				Channel:     coreAuth.AuthenticatorOOBChannelEmail,
				MaskedEmail: "joh****@example.com",
			},
		}
		So(actual, ShouldResemble, expected)
	})
}

func TestCanAddAuthenticator(t *testing.T) {
	type Existing struct {
		TOTP     int
		OOBSMS   int
		OOBEmail int
	}
	type Limit struct {
		Total    int
		TOTP     int
		OOBSMS   int
		OOBEmail int
	}
	type Case struct {
		Existing Existing
		New      interface{}
		Limit    Limit
		Expected bool
	}

	f := func(c Case) {
		var authenticators []interface{}
		for i := 0; i < c.Existing.TOTP; i++ {
			authenticators = append(authenticators, TOTPAuthenticator{})
		}
		for i := 0; i < c.Existing.OOBSMS; i++ {
			authenticators = append(authenticators, OOBAuthenticator{
				Channel: coreAuth.AuthenticatorOOBChannelSMS,
			})
		}
		for i := 0; i < c.Existing.OOBEmail; i++ {
			authenticators = append(authenticators, OOBAuthenticator{
				Channel: coreAuth.AuthenticatorOOBChannelEmail,
			})
		}

		newA := c.New

		maximum := &c.Limit.Total
		mfaConfiguration := config.MFAConfiguration{
			Maximum: maximum,
			TOTP: config.MFATOTPConfiguration{
				Maximum: c.Limit.TOTP,
			},
			OOB: config.MFAOOBConfiguration{
				SMS: config.MFAOOBSMSConfiguration{
					Maximum: c.Limit.OOBSMS,
				},
				Email: config.MFAOOBEmailConfiguration{
					Maximum: c.Limit.OOBEmail,
				},
			},
		}

		actual := CanAddAuthenticator(authenticators, newA, mfaConfiguration)
		So(actual, ShouldEqual, c.Expected)
	}

	cases := []Case{
		Case{
			Existing: Existing{
				TOTP:     0,
				OOBSMS:   0,
				OOBEmail: 0,
			},
			Limit: Limit{
				Total:    1,
				TOTP:     1,
				OOBSMS:   0,
				OOBEmail: 0,
			},
			New:      TOTPAuthenticator{},
			Expected: true,
		},

		Case{
			Existing: Existing{
				TOTP:     1,
				OOBSMS:   0,
				OOBEmail: 0,
			},
			Limit: Limit{
				Total:    1,
				TOTP:     1,
				OOBSMS:   0,
				OOBEmail: 0,
			},
			New:      TOTPAuthenticator{},
			Expected: false,
		},

		Case{
			Existing: Existing{
				TOTP:     1,
				OOBSMS:   0,
				OOBEmail: 0,
			},
			Limit: Limit{
				Total:    2,
				TOTP:     1,
				OOBSMS:   1,
				OOBEmail: 0,
			},
			New: OOBAuthenticator{
				Channel: coreAuth.AuthenticatorOOBChannelSMS,
			},
			Expected: true,
		},
	}
	Convey("CanAddAuthenticator", t, func() {
		for _, c := range cases {
			f(c)
		}
	})
}
