package cmd

import (
	// "github.com/documatrix/go-lprlib"
	lprlib "github.com/paulmatencio/lprlib"
	logz "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

// clientCmd represents the client command
var (
	clientCmd = &cobra.Command{
		Use:   "lpr",
		Short: "Send a file to an lpd server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			SendFile()
		},
	}
	filePath, prqName, userName string
	hostName                    string
	timeOut                            = 10 * time.Second
	maxSize                     uint64 = 32 * 1024
)

func init() {
	rootCmd.AddCommand(clientCmd)
	initLp(clientCmd)
}

func initLp(cmd *cobra.Command) {
	cmd.Flags().Uint16VarP(&port, "port", "p", 1515, "lpd server port")
	cmd.Flags().StringVarP(&filePath, "file-path", "f", "", "input file path")
	cmd.Flags().StringVarP(&prqName, "prt-queue", "q", "", "print queue")
	cmd.Flags().StringVarP(&hostName, "host-name", "H", "localhost", "host name")
	cmd.Flags().StringVarP(&userName, "user-name", "U", "", "user name")
}

func SendFile() {
	var err error
	if len(userName) == 0 {
		logz.Warn().Msg("user name is missing")
		return
	}
	lprlib.SetDebugLogger(log.Print)
	// lpr Init
	if err = lprlib.Send(filePath, hostName, port, prqName, userName, timeOut); err != nil {
		logz.Error().Err(err).Msg("Connection error")
	}
	return

}

func SendFile1() {
	var (
		err      error
		fileInfo os.FileInfo
		lpr      = lprlib.LprSend{}
	)
	if len(userName) == 0 {
		logz.Warn().Msg("user name is missing")
		return
	}
	lprlib.SetDebugLogger(log.Print)
	//  Init the connection
	logz.Info().Msg("Init a connection")
	err = lpr.Init(hostName, filePath, port, prqName, userName, timeOut)
	if err != nil {
		logz.Error().Err(err).Msg("Init connection error")
		return
	}

	// sending the configuration
	logz.Info().Msg("Sending the configuration")
	err = lpr.SendConfiguration()
	if err != nil {
		logz.Error().Err(err).Msg("Sending configuration error")
		return
	}
	defer lpr.Close()
	//sending the data
	if fileInfo, err = os.Stat(filePath); err != nil {
		logz.Error().Err(err).Msgf("Can't stat file %s", filePath)
		return
	}
	logz.Info().Msgf("Sending file %s - Size %d", filePath, fileInfo.Size())
	err = lpr.SendFile()
	if err != nil {
		logz.Error().Err(err).Msg("Sending file error")
		return
	}

	return
}
