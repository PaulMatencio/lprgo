/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"fmt"
	// lprlib "github.com/documatrix/go-lprlib"
	"github.com/paulmatencio/lprgo/lib"
	lprlib "github.com/paulmatencio/lprlib"
	logz "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var (
	serverCmd = &cobra.Command{
		Use:   "lpd",
		Short: "Start a lpd server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			StartLpd()
		},
	}

	lpr       lprlib.LprDaemon
	port      uint16
	ipAddress string
	waitTime  = 10 * time.Second
	saveDir   string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	initServer(serverCmd)

}

func initServer(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&ipAddress, "ip", "i", "0.0.0.0", "lpd server IP address")
	cmd.Flags().Uint16VarP(&port, "port", "p", 1515, "lpd server port")
	cmd.Flags().StringVarP(&saveDir, "save-directory", "S", "", "lpd server spool")
}

func StartLpd() {

	var (
		lprd = lprlib.LprDaemon{}
		err  error
	)

	lprlib.SetDebugLogger(log.Print)

	if saveDir != "" {

		if err = lib.CreateDir(saveDir); err == nil {
			lprd.InputFileSaveDir = saveDir
		} else {
			logz.Error().Err(err).Msg("Make SaveFile directory")
		}
	}

	err = lprd.Init(port, "")
	if err != nil {
		logz.Error().Err(err).Msg("Start lpdServer")
	}

	for {
		select {
		case conn := <-lprd.FinishedConnections():
			if conn.SaveName != "" {
				if fi, err := os.Stat(conn.SaveName); err == nil {
					fmt.Printf("Print Queue: %s - filename/filesize: %s/%d - Savename: %s\n", conn.PrqName, conn.Filename, fi.Size(), conn.SaveName)
				} else {

				}
			}
		case <-time.After(waitTime):
			fmt.Println("Waiting for a finished connection")
		}
	}
	return
}
