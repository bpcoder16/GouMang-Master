package bootstrap

import (
	"context"
	"goumang-master/route"
	"goumang-master/services/cron"
	"path"

	"github.com/bpcoder16/Chestnut/v2/appconfig"
	"github.com/bpcoder16/Chestnut/v2/appconfig/env"
	"github.com/bpcoder16/Chestnut/v2/bootstrap"
	"github.com/bpcoder16/Chestnut/v2/contrib/httphandler/gin"
	"github.com/bpcoder16/Chestnut/v2/core/gtask"
	"github.com/bpcoder16/Chestnut/v2/modules/httpserver"
)

func Start(ctx context.Context, config *appconfig.AppConfig) error {
	var g *gtask.Group
	g, ctx = gtask.WithContext(ctx)

	bootstrap.Start(ctx, config, g.Go)

	g.Go(func() error {
		return httpserver.NewManager(
			path.Join(env.ConfigDirPath(), "http.yaml"),
			gin.HTTPHandler(
				route.Api(),
			),
		).Run(ctx)
	})

	g.Go(func() error {
		cron.Run(ctx)
		return nil
	})

	return g.Wait()
}
