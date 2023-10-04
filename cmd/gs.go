/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/

package cmd

import (
	"bufio"
	"bytes"
	// lprlib  "github.com/documatrix/go-lprlib"
	lprlib "github.com/paulmatencio/lprlib"
	logz "github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"os/exec"
)

// gsCmd represents the gs command
var (
	gsCmd = &cobra.Command{
		Use:   "gs",
		Short: "execute ghostscript to convert Pdf or ps to PCL4/5",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			gs(args)
		},
	}
	gsLprCmd = &cobra.Command{
		Use:   "gsLpr",
		Short: "execute ghostscript to convert Pdf or ps file  to PCL format  then send the result to a remote printer server",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			gsLpr(args)
		},
	}
	inputFile, inputFile1, outputFile, outputDevice string
)

func init() {
	rootCmd.AddCommand(gsCmd)
	rootCmd.AddCommand(gsLprCmd)
	initGS(gsCmd)
	initGSLpr(gsLprCmd)
}

func initGS(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&inputFile, "if", "i", "", "input file")
	cmd.Flags().StringVarP(&outputFile, "of", "o", "", "output file")
	cmd.Flags().StringVarP(&outputDevice, "od", "d", "ljet4", "output device")
}

func initGSLpr(cmd *cobra.Command) {
	cmd.Flags().Uint16VarP(&port, "port", "p", 1515, "lpd server port")
	cmd.Flags().StringVarP(&inputFile, "file-path", "f", "", "input file path")
	cmd.Flags().StringVarP(&prqName, "prt-queue", "q", "p0001", "print queue")
	cmd.Flags().StringVarP(&hostName, "host-name", "H", "localhost", "host name")
	cmd.Flags().StringVarP(&userName, "user-name", "U", "paul", "user name")
	cmd.Flags().StringVarP(&outputDevice, "od", "d", "ljet4", "output device")
}

/*
			gs Input : input file  or  stdin if input file is missing  -> gs
	        if gs Input = stdin
				Cat args[0] (input file)  | gs
*/

/*
	call ghostscript
		output to file |  output to stdout
		if output to stdout ->  cat the result to a file
		Purpose :  Just to show how to work with Pipe.io
*/

func gs(args []string) {
	var (
		err           error
		ipf           string
		gsStdin       io.WriteCloser
		catStdout     io.ReadCloser
		catCMD, gsCMD *exec.Cmd
	)
	if len(args) > 0 {
		inputFile1 = args[0]
	}
	if inputFile != "" {
		if _, err = os.Stat(inputFile); err != nil {
			logz.Warn().Msgf("Stat input file:%s - error:%v", inputFile, err)
			return
		}
		ipf = "-f" + inputFile
	} else if inputFile1 != "" {
		logz.Info().Msg("input file missing - stdin will be used")
		ipf = "-"
	}
	if outputFile == "" {
		logz.Warn().Msg("Output file missing")
		return
	}

	/*    Switches
	    -dSAFER      	Restricts file operations the job can perform
	    -dNOPAUSE 		Disables the prompt and pause at the end of each page.
		-dBATCH  		exit after last file
	    -sDEVICE  		select device
	    -sOutputFile	output file
		-dQUIETE   		quit  mode
	    -dFirstPage  	first page
	    -dLastPage     	last Page
		-sPAPERSIZE		select a specific paper size

	*/

	gsCMD = exec.Command("gs", "-q",
		"-dNOPAUSE", "-dBATCH", "-dSAFER",
		"-sDEVICE="+outputDevice,
		"-sOutputFile="+outputFile, ipf)

	if ipf == "-" {
		/*
			get stdin of the gs command
		*/
		if gsStdin, err = gsCMD.StdinPipe(); err != nil {
			logz.Error().Err(err).Msg("gsCMD stdinPipe")
		}
		/*  issue the cat command */
		catCMD = exec.Command("cat", inputFile1)
		/*  get the stdout of the cat command */
		catStdout, err = catCMD.StdoutPipe()
		if err != nil {
			logz.Error().Err(err).Msg("cat  stdoutPipe")
		}
		//  gs stdin = cat stdout
		gsCMD.Stdin = catStdout

		defer gsStdin.Close()
		defer catStdout.Close()
	}

	if err = gsCMD.Start(); err == nil {
		if ipf == "-" {
			err1 := catCMD.Start()
			if err1 != nil {
				logz.Error().Err(err1).Msg("catDMD start")
				return
			}
			if err1 = catCMD.Wait(); err1 != nil {
				logz.Error().Err(err1).Msg("catCMD wait")
				return
			}
		}

		err = gsCMD.Wait()
		if err != nil {
			logz.Error().Err(err).Msg("gsCMD Wait")
		}
	} else {
		logz.Error().Err(err).Msg("gsCMD Start")
	}

}

/*
	exec ghostscript  to convert
		PDF file or PS file to PCL
		Send the result ( PCL file) to a remote printer server (lpd)

*/

func gsLpr(args []string) {
	var (
		err      error
		ipf      string
		gsCMD    *exec.Cmd
		gsOutput io.ReadCloser
	)

	if inputFile != "" {
		if _, err = os.Stat(inputFile); err != nil {
			logz.Warn().Msgf("Stat input file:%s - error:%v", inputFile, err)
			return
		}
		ipf = "-f" + inputFile
	} else {
		logz.Info().Msg("input file is missing")
		return
	}

	gsCMD = exec.Command("gs", "-q",
		"-dNOPAUSE", "-dBATCH", "-dSAFER",
		"-sDEVICE="+outputDevice,
		"-sOutputFile=-", ipf)

	// connect to the  gs stdout
	gsOutput, err = gsCMD.StdoutPipe()
	defer gsOutput.Close() // not required  since Wait will close it
	gsCMD.Stderr = os.Stderr
	gsCMD.Stdin = os.Stdin
	// Start ghostscript
	err = gsCMD.Start()
	if err != nil {
		logz.Error().Err(err).Msg("Start gsCMD")
		return
	}

	lprlib.SetDebugLogger(log.Print)
	// lpr Init
	if err = lprlib.SendStdin(inputFile, gsOutput, hostName, port, prqName, userName, timeOut); err != nil {
		logz.Error().Err(err).Msg("Connection error")
	}

	// Wait ghostscript to finish
	err = gsCMD.Wait()
	if err != nil {
		logz.Error().Err(err).Msg("Wait gsCMD")
		return
	}

}

func readGsOutput(stdout io.ReadCloser) ([]byte, error) {
	//
	// You  can  use  io.ReadAll instead of readGs
	//
	var (
		n   int
		b1  bytes.Buffer
		err error
	)
	b := make([]byte, 4096)
	reader := bufio.NewReader(stdout)
	for {
		n, err = reader.Read(b)

		if err == nil || err == io.EOF {
			if n > 0 {
				b1.Write(b[:n])
			}
		} else {
			break
		}
		if err == io.EOF {
			// fmt.Println(string(b1.Bytes()))
			err = nil
			break
		}
	}
	return b1.Bytes(), err
}

// write to /temp/file  just for testing
func writeGsResult(bi []byte) (err error) {
	var f1 *os.File
	f1, err = os.Open("/tmp/t1")
	if err == nil {
		defer f1.Close()
		f1.Write(bi)
	}
	return
}
