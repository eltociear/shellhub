package main

import (
	"context"
	"github.com/shellhub-io/shellhub/test/api/community"

	"github.com/shellhub-io/shellhub/test/api"
	"github.com/shellhub-io/shellhub/test/database"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func main() {
	ctx := context.Background()
	log.SetLevel(log.TraceLevel)

	r := &cobra.Command{Use: "test"}
	r.AddCommand(&cobra.Command{
		Use:     "api",
		Example: "test api <address>",
		Args:    cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tests := community.Tests(args[0])

			err := api.Init(args[0], tests)

			if err != nil {
				log.Errorln(err)
				return
			}
		},
	})

	d := &cobra.Command{Use: "database"}
	d.AddCommand(&cobra.Command{
		Use:     "populate",
		Example: "test database populate",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.Populate(ctx)
			if err != nil {
				log.Errorln(err)
				return
			}
		},
	})
	d.AddCommand(&cobra.Command{
		Use:     "clean",
		Example: "test database clean",
		Run: func(cmd *cobra.Command, args []string) {
			err := database.Clean(ctx)
			if err != nil {
				log.Errorln(err)
				return
			}
		},
	})

	r.AddCommand(d)
	if err := r.Execute(); err != nil {
		log.Error(err)
	}
}
