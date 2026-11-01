package admin

import (
	"github.com/RobertsMJ/simc-cloud/backend/platform/httpx"
	"go.uber.org/fx"
)

const adminModule = "admin"

var Module = fx.Module(
	adminModule,
	httpx.NewWithConfig(adminModule, NewAdminHandler, NewAdminServerConfig),
)
