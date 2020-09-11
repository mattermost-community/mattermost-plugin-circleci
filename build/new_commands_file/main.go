// main handles creation of a template of a server/commands_template.go file
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

const commandsFileTemplate = `package main

import (
	"fmt"

	"github.com/jszwedko/go-circleci"

	"github.com/mattermost/mattermost-server/v5/model"
)

const (
	%sTrigger            = "TODO"
	%sHint               = "TODO"
	%sHelpText           = "TODO"
	// TODO:  add subcommands Triggers, Hints and HelpTexts here

	// TODO: add theses subCommands in getAutocompleteData()
)

func (p *Plugin) execute%s(args *model.CommandArgs, circleciToken string, split []string) (*model.CommandResponse, *model.AppError) {
	subcommand := commandHelpTrigger
	if len(split) > 0 {
		subcommand = split[0]
	}

	switch subcommand {
	// TODO: add cases with subcommand triggers here

	case commandHelpTrigger:
		return p.sendHelpResponse(args, %sTrigger)

	default:
		return p.sendIncorrectSubcommandResponse(args, %sTrigger)
	}
}

// TODO: implements the subcommands

`

func main() {
	if len(os.Args) <= 1 {
		panic("Please precise the name of the new command group (in one word)")
	}

	commandGroupName := os.Args[1]
	fileName := fmt.Sprintf("commands_%s.go", commandGroupName)

	ok := checkFileDoesNotExistInServerDir(fileName)
	if !ok {
		panic("This file already exists!")
	}

	err := createFileInServerDir(fileName, commandGroupName)
	if err != nil {
		panic(err)
	}
}

func checkFileDoesNotExistInServerDir(fileName string) bool {
	path := filepath.Join(".", "server", fileName)

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return false
	}

	return true
}

func createFileInServerDir(fileName string, commandGroupName string) error {
	path := filepath.Join(".", "server", fileName)

	// write generated code to file by using Go file template.
	if err := ioutil.WriteFile(
		path,
		[]byte(fmt.Sprintf(
			commandsFileTemplate,
			commandGroupName,
			commandGroupName,
			commandGroupName,
			strings.Title(commandGroupName),
			commandGroupName,
			commandGroupName,
		)),
		0600,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to write server/%s", fileName))
	}

	return nil
}
