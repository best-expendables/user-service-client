package userclient

import "time"

// User roles which should be identical with user-service
const (
	RoleAdmin                  = "Admin"
	RoleShippingProvider       = "ShippingProvider"
	RolePlatformTransportation = "PlatformTransportation"
	RolePlatformWarehouse      = "PlatformWarehouse"
	RolePlatformOverview       = "PlatformOverview"
	RoleLocationTree           = "LocationTree"
	RoleAccessAdmin            = "AccessAdmin"
	RoleReturn                 = "Return"
	RmsAdmin                   = "RmsAdmin"
	RmsReturnOperator          = "RmsReturnOperator"
	RmsCbOperator              = "RmsCbOperator"
	RmsApi                     = "RmsApi"
	RoleRmsApplicationAdmin    = "RmsApplicationAdmin"
	RmsCs                      = "RmsCs"
	RmsAdminID                 = "ID_RMS_ADMIN"
	RmsAdminVN                 = "VN_RMS_ADMIN"
	RmsAdminTH                 = "TH_RMS_ADMIN"
	RmsAdminSG                 = "SG_RMS_ADMIN"
	RmsAdminPH                 = "PH_RMS_ADMIN"
	RmsAdminPK                 = "PK_RMS_ADMIN"
	RmsAdminBD                 = "BD_RMS_ADMIN"
	RmsAdminLK                 = "LK_RMS_ADMIN"
	RmsAdminNP                 = "NP_RMS_ADMIN"
	RmsAdminMM                 = "MM_RMS_ADMIN"
	RmsAdminMY                 = "MY_RMS_ADMIN"
	RmsReturnOperatorID        = "ID_RMS_RETURN-OPERATOR"
	RmsReturnOperatorVN        = "VN_RMS_RETURN-OPERATOR"
	RmsReturnOperatorTH        = "TH_RMS_RETURN-OPERATOR"
	RmsReturnOperatorSG        = "SG_RMS_RETURN-OPERATOR"
	RmsReturnOperatorPH        = "PH_RMS_RETURN-OPERATOR"
	RmsReturnOperatorMY        = "MY_RMS_RETURN-OPERATOR"
	RmsReturnOperatorPK        = "PK_RMS_RETURN-OPERATOR"
	RmsReturnOperatorBD        = "BD_RMS_RETURN-OPERATOR"
	RmsReturnOperatorLK        = "LK_RMS_RETURN-OPERATOR"
	RmsReturnOperatorNP        = "NP_RMS_RETURN-OPERATOR"
	RmsReturnOperatorMM        = "MM_RMS_RETURN-OPERATOR"
	RmsCsID                    = "ID_RMS_CS"
	RmsCsVN                    = "VN_RMS_CS"
	RmsCsTH                    = "TH_RMS_CS"
	RmsCsSG                    = "SG_RMS_CS"
	RmsCsPH                    = "PH_RMS_CS"
	RmsCsMY                    = "MY_RMS_CS"
	RmsCsPK                    = "PK_RMS_CS"
	RmsCsBD                    = "BD_RMS_CS"
	RmsCsLK                    = "LK_RMS_CS"
	RmsCsNP                    = "NP_RMS_CS"
	RmsCsMM                    = "MM_RMS_CS"
	RmsCrossborderOperatorID   = "ID_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorVN   = "VN_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorTH   = "TH_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorSG   = "SG_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorPH   = "PH_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorMY   = "MY_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorPK   = "PK_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorBD   = "BD_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorLK   = "LK_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorNP   = "NP_RMS_CROSSBORDER-OPERATOR"
	RmsCrossborderOperatorMM   = "MM_RMS_CROSSBORDER-OPERATOR"
)

type Platform struct {
	Name string `json:"name"`
}

type User struct {
	Id            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	Active        bool      `json:"active"`
	Password      string    `json:"password"`
	Roles         []string  `json:"roles"`
	UpdatedAt     time.Time `json:"updatedAt"`
	CreatedAt     time.Time `json:"createdAt"`
	PlatformNames []string  `json:"platforms"`
}

type RevokedToken struct {
	Token     string    `json:"token"`
	ExpiredAt time.Time `json:"expiredAt"`
}

func (t *RevokedToken) TTL() time.Duration {
	return t.ExpiredAt.Sub(time.Now())
}

func (u User) HasRole(roles ...string) bool {
	for _, r := range roles {
		for _, ur := range u.Roles {
			if ur == r {
				return true
			}
		}
	}

	return false
}

func (u *User) HasAccessToPlatform(name string) bool {
	if !u.HasPlatformRole() {
		return true
	}
	for _, platformName := range u.PlatformNames {
		if platformName == name {
			return true
		}
	}
	return false
}

func (u *User) HasPlatformRole() bool {
	if u.HasRole(RolePlatformTransportation, RolePlatformWarehouse, RolePlatformOverview) {
		return true
	}
	return false
}
