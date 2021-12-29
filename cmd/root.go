// Dinkur the task time tracking utility.
// <https://github.com/dinkur/dinkur>
//
// Copyright (C) 2021 Kalle Fagerberg
// SPDX-FileCopyrightText: 2021 Kalle Fagerberg
// SPDX-License-Identifier: GPL-3.0-or-later
//
// This program is free software: you can redistribute it and/or modify it under
// the terms of the GNU General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your option)
// any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS
// FOR A PARTICULAR PURPOSE.  See the GNU General Public License for more
// details.
//
// You should have received a copy of the GNU General Public License along with
// this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/dinkur/dinkur/internal/cfgpath"
	"github.com/dinkur/dinkur/pkg/dinkurdb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile = cfgpath.Path()
var db dinkurdb.Client

// RootCMD represents the base command when called without any subcommands
var RootCMD = &cobra.Command{
	Use:   "dinkur",
	Short: "The Dinkur CLI",
	Long:  `Through these subcommands you can access your time-tracked tasks.`,
	Run: func(cmd *cobra.Command, args []string) {
		connectAndMigrateDB()
		activeTask, err := db.ActiveTask()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error getting active task:", err)
			os.Exit(1)
		}
		if activeTask != nil {
			fmt.Println("Current task:", *activeTask)
		} else {
			fmt.Println("You have no active task.")
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	db = dinkurdb.NewClient()
	defer db.Close()
	err := RootCMD.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	RootCMD.PersistentFlags().StringVar(&cfgFile, "config", cfgFile, "config file")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetConfigFile(cfgFile)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	} else if !errors.As(err, &viper.ConfigFileNotFoundError{}) && !errors.Is(err, os.ErrNotExist) {
		fmt.Fprintln(os.Stderr, "Error reading config:", err)
		os.Exit(1)
	}
}

func connectAndMigrateDB() {
	if err := db.Connect("dinkur.db"); err != nil {
		fmt.Fprintln(os.Stderr, "Error connecting to database:", err)
		os.Exit(1)
	}
	if err := db.Migrate(); err != nil {
		fmt.Fprintln(os.Stderr, "Error migrating database:", err)
		os.Exit(1)
	}
}
