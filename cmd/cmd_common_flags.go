/*
Copyright 2024 4rcadia

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"errors"

	"github.com/urfave/cli/v3"
	"lucy/lucytypes"
)

var flagJsonOutput = &cli.BoolFlag{
	Name:  "json",
	Usage: "Print raw JSON response",
	Value: false,
}

var flagLongOutput = &cli.BoolFlag{
	Name:    "long",
	Usage:   "Show hidden or collapsed output",
	Value:   false,
	Aliases: []string{"l"},
}

func sourceFlag(absent lucytypes.Source) *cli.StringFlag {
	return &cli.StringFlag{
		Name:    "source",
		Aliases: []string{"s"},
		Usage:   "To fetch info from `SOURCE`",
		Value:   absent.String(),
		Validator: func(s string) error {
			if s != lucytypes.UnknownSource.String() {
				return errors.New("unknown source " + s)
			}
			return nil
		},
	}
}

var flagNoStyle = &cli.BoolFlag{
	Name:  "no-style",
	Usage: "Disable colored and styled output",
	Value: false,
}
