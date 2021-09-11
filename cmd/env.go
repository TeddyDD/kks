package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
)

func NewEnvCmd() *EnvCmd {
	c := &EnvCmd{
		Cmd: Cmd{
			fs:         flag.NewFlagSet("env", flag.ExitOnError),
			alias:      []string{""},
			usageStr:   "[options]",
			sessionReq: true,
		},
	}
	c.fs.BoolVar(&c.json, "json", false, "json output")
	return c
}

type EnvCmd struct {
	Cmd
	json bool
}

func (c *EnvCmd) Run() error {
	switch c.json {
	case true:
		j, err := json.MarshalIndent(
			map[string]string{
				"session": c.session,
				"client":  c.client,
			}, "", "  ")
		if err != nil {
			return err
		}
		fmt.Println(string(j))
	case false:
		fmt.Printf("session: %s\n", c.session)
		fmt.Printf("client: %s\n", c.client)
		// fmt.Printf("workdir: %s\n", cc.WorkDir)
		// fmt.Printf("buffer: %s\n", cc.Buffer)
	}
	return nil
}
