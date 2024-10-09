package cmd

import (
	"context"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.lsp.dev/jsonrpc2"

	"github.com/ArnauLlamas/terragrunt-ls/internal/handler"
	_ "github.com/ArnauLlamas/terragrunt-ls/internal/log"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Terragrunt language server",
	RunE: func(_ *cobra.Command, _ []string) error {
		conn := jsonrpc2.NewConn(jsonrpc2.NewStream(rwc{}))
		handler := handler.NewHandler(conn)
		handlerSrv := jsonrpc2.HandlerServer(handler)

		log.Info("Starting server")

		return handlerSrv.ServeStream(context.Background(), conn)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type rwc struct{}

func (rwc) Read(r []byte) (int, error) {
	return os.Stdin.Read(r)
}

func (rwc) Write(w []byte) (int, error) {
	return os.Stdout.Write(w)
}

func (rwc) Close() error {
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	return os.Stdout.Close()
}
