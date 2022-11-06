package main

import (
	"context"
	"fmt"
	"gmail-sender/internal/factory"
	"time"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type StartOption struct {
	Logger    *zap.Logger
	Port      string
	TokenHost string
	TokenPort string
}

var startOpt StartOption

func init() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = jst
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	SilenceErrors: true,
	SilenceUsage:  true,
	RunE: func(cmd *cobra.Command, args []string) error {
		return start(&startOpt)
	},
}

func start(opts *StartOption) error {
	l, err := factory.NewLogger()
	if err != nil {
		fmt.Println(err)
		return err
	}

	tr := factory.NewTokenRepo(opts.TokenHost, opts.TokenPort)
	gc := factory.NewGmailClient()
	df := factory.NewDefaultContent()
	uc := factory.NewUsecase(l, tr, gc, df)
	sv := factory.NewServer(l, opts.Port, uc)
	ctx := context.Background()

	uc.Logger.Info("load setting", zap.String("default_from", uc.DefaultContent().From), zap.String("default_to", uc.DefaultContent().To), zap.String("default_title", uc.DefaultContent().Title))
	return sv.Start(ctx)
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.Flags().StringVar(&startOpt.Port, "port", "80", "listen port")
	startCmd.Flags().StringVar(&startOpt.TokenHost, "token-host", "token-repository-api", "token-repository host")
	startCmd.Flags().StringVar(&startOpt.TokenPort, "token-port", "80", "token-repository port")
}
