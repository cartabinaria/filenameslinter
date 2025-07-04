// SPDX-FileCopyrightText: 2023 - 2024 Eyad Issa <eyadlorenzo@gmail.com>
// SPDX-FileCopyrightText: 2023 Luca Tagliavini <luca@teapot.ovh>
// SPDX-FileCopyrightText: 2023 Samuele Musiani <samu@teapot.ovh>
// SPDX-FileCopyrightText: 2024 Stefano Volpe <foxy@teapot.ovh>
// SPDX-FileCopyrightText: 2023 konova8 <canovasimo@gmail.com>
//
// SPDX-License-Identifier: AGPL-3.0-only

package main

import (
	"flag"
	"fmt"

	"os"
	"path"
	"path/filepath"
	"strings"

	log "log/slog"

	"github.com/cartabinaria/synta"

	"github.com/cartabinaria/filenameslinter"
)

var (
	recursive         = flag.Bool("recursive", true, "Recursively check all files")
	ensureKebabCasing = flag.Bool("ensure-kebab-casing", true, "Check if directory names are in kebab-case")
	ignoreDotfiles    = flag.Bool("ignore-dotfiles", true, "Ignore files and folders that start with a dot")
	syntaDefinition   = flag.String("definition", "", "Synta definition file to check filenames against")
	failFast          = flag.Bool("failfast", false, "Stop checking as soon as an error is found")
)

func main() {
	flag.Parse()

	githubActions := false
	if os.Getenv("GITHUB_ACTIONS") == "true" {
		githubActions = true
		log.Info("running in github actions")
	}

	pwd, err := os.Getwd()
	if err != nil {
		log.Error("could not get current working directory", "err", err)
		os.Exit(1)
	}

	dirPath := "."
	parent := pwd
	if len(flag.Args()) > 0 {
		absDir := path.Join(pwd, flag.Arg(0))
		parent = filepath.Dir(strings.TrimSuffix(absDir, string(os.PathSeparator)))
		dirPath, err = filepath.Rel(parent, absDir)

		if err != nil {
			log.Error("could not make the path relative", "err", err)
			os.Exit(2)
		}
	}

	var syntaFile *synta.Synta = nil
	if *syntaDefinition != "" {
		data, err := os.ReadFile(*syntaDefinition)
		if err != nil {
			log.Error("could not read synta definition file", "err", err)
			os.Exit(3)
		}
		s, err := synta.ParseSynta(string(data))
		if err != nil {
			log.Error("invalid synta definiton file", "err", err)
			os.Exit(4)
		}
		syntaFile = &s
	}

	fs := os.DirFS(parent)
	_, err = filenameslinter.ReadDir(fs, dirPath)
	// Sloppy: if the directory passed as argument can't be properly examined in
	// first place (say, because it does not exist), the check passes
	if err != nil {
		os.Exit(0)
	}

	opts := filenameslinter.Options{
		Recursive:         *recursive,
		EnsureKebabCasing: *ensureKebabCasing,
		IgnoreDotfiles:    *ignoreDotfiles,
		FailFast:          *failFast,
	}
	errs := filenameslinter.CheckDir(syntaFile, fs, dirPath, &opts)
	if len(errs) > 0 {
		log.Error("error while checking directory", "path", dirPath, "errors", len(errs))
		for _, e := range errs {
			log.Error("found error", "err", e)

			if regexErr, ok := e.(filenameslinter.RegexMatchError); ok && githubActions {
				fmt.Printf("::error file=%s::%s\n", path.Join(pwd, regexErr.Path), e.Error())
			}
		}
		os.Exit(5)
	}

	os.Exit(0)
}
