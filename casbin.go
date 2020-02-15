package gqlgen_casbin_RBAC_example
import (
	"github.com/casbin/casbin/v2"
)

var (
	enforcer *casbin.Enforcer
)
func init(){
	initEnforcer()
}

func initEnforcer() {
	e, err := casbin.NewEnforcer("./rbac_with_domains_model.conf", "./rbac_with_domains_policy.csv")
	if err!=nil{
		panic(err)
	}
	enforcer=e
}

func Enforcer() *casbin.Enforcer {
	return enforcer
}

