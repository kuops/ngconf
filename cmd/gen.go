package cmd

import (
	"github.com/kuops/ngconf/pkg/gen"

	"github.com/spf13/cobra"
)

var (
	genOptions gen.Options
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:                   "gen [--cloud=bool] [--domains=\"domain1,domain2,...\"] [--https=bool] [--nodes=\"node1\",\"node2\"] [--port=int] --projectname=\"value\" [--file=demo.conf] [--preview=bool]",
	DisableFlagsInUseLine: true,
	Short:                 "genterate a nginx.conf for proxy",
	Long:                  "genterate proxy nginx config files",
	Example: `
	ngconf gen --cloud=false --domains=kuops.com,dns.kuops.com --https=false --nodes=g1-sre-jenkins-v01,g1-sre-jenkins-v02 --port=8080 --projectname=demo --file=demo.conf --preview=true
	`,
	PreRun: func(cmd *cobra.Command, args []string) {
		return
	},
	Run: func(cmd *cobra.Command, args []string) {
		gen.GenConfig(&genOptions)
	},
}

func init() {
	rootCmd.AddCommand(genCmd)
	dopts := *gen.DefaultOptions()

	genCmd.PersistentFlags().StringVar(&genOptions.ProjectName, "projectname", dopts.ProjectName, "projectname is git repo name, use to domain prefix.")
	genCmd.PersistentFlags().StringSliceVar(&genOptions.Domains, "domains", dopts.Domains, "domain name")
	genCmd.PersistentFlags().BoolVar(&genOptions.IsCloud, "cloud", dopts.IsCloud, "is cloud project include medusa-online-service")
	genCmd.PersistentFlags().StringSliceVar(&genOptions.Nodes, "nodes", dopts.Nodes, "nginx upstream backend server")
	genCmd.PersistentFlags().UintVar(&genOptions.BackendPort, "port", dopts.BackendPort, "port is backend port")
	genCmd.PersistentFlags().BoolVar(&genOptions.ForceHTTPS, "https", dopts.ForceHTTPS, "use https redirect")
	genCmd.PersistentFlags().StringVar(&genOptions.WriteFileName, "file", dopts.WriteFileName, "write to files")
	genCmd.PersistentFlags().BoolVar(&genOptions.Preview, "preview", dopts.Preview, "is preview config files")
}
