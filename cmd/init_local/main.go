package main

import (
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"github.com/Kazuki-Ya/wmd-server/log-server/agent"
	"github.com/Kazuki-Ya/wmd-server/web-server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	cli := &cli{}
	cmd := &cobra.Command{
		Use:     "log-inference-setup",
		PreRunE: cli.setupConfig,
		RunE:    cli.run,
	}

	if err := setupFlags(cmd); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

type cli struct {
	cfg cfg
}

type cfg struct {
	agent.Config
}

func setupFlags(cmd *cobra.Command) error {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Flags().String("config-file", "", "Path to config file.")

	dataDir := path.Join(os.TempDir(), "wmd-server")
	cmd.Flags().String(
		"data-dir",
		dataDir,
		"Directory to store log and Raft data.",
	)
	cmd.Flags().String("node-name", hostname, "Unique server ID.")
	cmd.Flags().String("bind-addr",
		"127.0.0.1:8401",
		"Address to bind serf on.")
	cmd.Flags().Int("rpc-port",
		8400,
		"Port for RPC clients (and Raft) connections.")
	cmd.Flags().StringSlice("start-join-addrs",
		nil,
		"Serf addresses to join")
	cmd.Flags().Bool("bootstrap", true, "Bootstrap the cluster")

	return viper.BindPFlags(cmd.Flags())
}

func (c *cli) setupConfig(cmd *cobra.Command, args []string) error {
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return err
	}

	viper.SetConfigFile(configFile)

	if err = viper.ReadInConfig(); err != nil {
		// 設定ファイルは、存在しなくても問題ない
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	c.cfg.DataDir = viper.GetString("data-dir")
	c.cfg.NodeName = viper.GetString("node-name")
	c.cfg.BindAddr = viper.GetString("bind-addr")
	c.cfg.RPCPort = viper.GetInt("rpc-port")
	c.cfg.StartJoinAddrs = viper.GetStringSlice("start-join-addrs")
	c.cfg.Bootstrap = viper.GetBool("bootstrap")

	return nil
}

func (c *cli) run(cmd *cobra.Command, args []string) error {
	agent, err := agent.New(c.cfg.Config)
	web.WebInit()
	if err != nil {
		return err
	}
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGINT, syscall.SIGTERM)
	<-sigc
	return agent.Shutdown()
}
