/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/

package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "lprgo",
	Short: "test lpr with go library and command pipe",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
