package cmds

import (
	"fmt"
	"log"
	"os"

	"github.com/opcoder0/typesense-cli/pkg/config"
	"github.com/spf13/cobra"
	"github.com/typesense/typesense-go/typesense"
)

var (
	instance string
)

// New returns a new typesense shell command tree
func New() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "tscli",
		Short: "tscli is a Typesense command line client to query, create, update or delete documents",
	}
	rootCmd.PersistentFlags().StringVarP(&instance, "instance", "i", "default", "config.toml instance")

	collectionsCmd := &cobra.Command{
		Use:   "collections",
		Short: "collections operate on typesense collections",
	}
	rootCmd.AddCommand(collectionsCmd)

	listCollectionsCmd := &cobra.Command{
		Use:   "list",
		Short: "List collections",
		Run:   listCmd,
	}
	collectionsCmd.AddCommand(listCollectionsCmd)
	return rootCmd
}

func newClient() *typesense.Client {
	conf, err := config.Load(instance)
	if err != nil {
		log.Fatal(err)
	}
	client := typesense.NewClient(
		typesense.WithServer(fmt.Sprintf("http://%s:%d", conf.Host, conf.Port)),
		typesense.WithAPIKey(conf.APIKey))
	return client
}

func listCmd(cmd *cobra.Command, args []string) {
	client := newClient()
	collections, err := client.Collections().Retrieve()
	if err != nil {
		fmt.Println("List collections: ", err)
		os.Exit(1)
	}
	for _, c := range collections {
		fmt.Println(c.Name)
	}
}
