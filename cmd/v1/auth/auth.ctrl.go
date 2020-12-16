package auth

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/kyfk/gin-jwt"
	"net/http"
)

type Role int

const (
	OPERATOR     Role = 0x1
	ADMIN        Role = 0x1 << 1
	SYSTEM_ADMIN Role = 0x1 << 2
)

func (r Role) IsOperator() bool {
	return r&OPERATOR != 0
}

func (r Role) IsAdmin() bool {
	return r&ADMIN != 0
}

func (r Role) IsSystemAdmin() bool {
	return r&SYSTEM_ADMIN != 0
}

var naiveDatastore = map[string]User{
	"operator":     {"operator", "o@o.com", "operator", OPERATOR},
	"admin":        {"admin", "a@a.com", "admin", OPERATOR | ADMIN},
	"system_admin": {"system_admin", "sa@sa.com", "system_admin", SYSTEM_ADMIN},
}

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"-"` // here is just for example
	Role     Role   `json:"-"`
}

func NewAuth() (jwt.Auth, error) {
	return jwt.New(jwt.Auth{
		SecretKey: []byte("must change here"),
		Authenticator: func(c *gin.Context) (jwt.MapClaims, error) {
			var req struct {
				Username string `json:"username"`
				Password string `json:"password"`
			}
			if err := c.ShouldBind(&req); err != nil {
				return nil, jwt.ErrorAuthenticationFailed
			}

			u := naiveDatastore[req.Username] // change here fetching from read datastore
			if u.Password != req.Password {
				return nil, jwt.ErrorAuthenticationFailed
			}

			return jwt.MapClaims{
				"username": u.Username,
				"role":     u.Role,
			}, nil
		},
		UserFetcher: func(c *gin.Context, claims jwt.MapClaims) (interface{}, error) {
			username, ok := claims["username"].(string)
			if !ok {
				return nil, nil
			}
			u, ok := naiveDatastore[username]
			if !ok {
				return nil, nil
			}
			return u, nil
		},
	})
}

func Operator(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims).IsOperator()
	})
}

func Admin(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims).IsAdmin()
	})
}

func SystemAdmin(m jwt.Auth) gin.HandlerFunc {
	return m.VerifyPerm(func(claims jwt.MapClaims) bool {
		return role(claims).IsSystemAdmin()
	})
}

func role(claims jwt.MapClaims) Role {
	return Role(claims["role"].(float64))
}

func SayHello(c *gin.Context) {
	u := jwt.User(c).(User)
	c.JSON(http.StatusOK, struct {
		User User `json:"user"`
	}{u})
}

// Register routes for creating account
func Register(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Account_created",
	})
}

// Login routes for creating account
func Login(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Logged",
	})
}
