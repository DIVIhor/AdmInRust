package server

import (
	"adminrust/internal/database"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

// Command-related routes
func (s *Server) registerPluginCmdRoutes(r chi.Router) {
	r.Route("/commands", func(r chi.Router) {
		r.Get("/", s.getPluginCommands)

		r.Get("/add", s.addPluginCommandsForm)
		r.Post("/add", s.addPluginCommands)
	})
}

// Retrieve command list
func (s *Server) getPluginCommands(w http.ResponseWriter, r *http.Request) {
	pluginSlug := r.PathValue("pluginSlug")

	commands, err := s.db.Queries().GetPluginCommands(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	renderPage(w, "plugin_commands", "", commands, nil)
}

// Render a page with plugin commands addition form
func (s *Server) addPluginCommandsForm(w http.ResponseWriter, r *http.Request) {
	renderPage(w, "add_plugin_cmds", "Add Plugin Commands", nil, nil)
}

// Add plugin commands
func (s *Server) addPluginCommands(w http.ResponseWriter, r *http.Request) {
	// get plugin ID
	pluginSlug := r.PathValue("pluginSlug")
	pluginID, err := s.db.Queries().GetPluginID(r.Context(), pluginSlug)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
		return
	}

	// verify inputs
	rawCommands := r.FormValue("commands")
	if rawCommands == "" {
		log.Println("Empty command list")
		return
	}
	descrSep := r.FormValue("descr-sep")
	if descrSep == "" {
		descrSep = "-"
	}
	cmdSep := r.FormValue("cmd-sep")
	if cmdSep == "" {
		cmdSep = "\n"
	}

	// parse and convert commands
	commandsMap, err := parseCommands(rawCommands, descrSep, cmdSep)
	if err != nil {
		log.Println(err)
		return
	}
	var commandArgs []database.AddPluginCommandsParams
	for cmd, descr := range commandsMap {
		commandArgs = append(commandArgs, database.AddPluginCommandsParams{
			PluginID:    pluginID,
			Command:     cmd,
			Description: descr,
		})
	}

	// save commands to DB
	_, err = s.db.Queries().AddPluginCommands(r.Context(), commandArgs)
	if err != nil {
		log.Println(err)
		internalServerErr(w)
	}
}

// Split input string on rows then split each row on command and its description
func parseCommands(rawCommands, descrSep, cmdSep string) (commands map[string]string, err error) {
	commandRows := strings.Split(rawCommands, cmdSep)
	if len(commandRows) == 1 {
		err = fmt.Errorf("wrong command separator: <%s> or nothing to split", cmdSep)
		return
	}

	// extract and save commands and their descriptions
	commands = map[string]string{}
	for _, row := range commandRows {
		rowEls := strings.SplitN(row, descrSep, 2)
		if len(rowEls) == 2 {
			// clean and save command and description
			cmd := strings.Trim(rowEls[0], " ")
			descr := strings.Trim(rowEls[1], " ")
			commands[cmd] = descr
		}
	}

	return commands, err
}
