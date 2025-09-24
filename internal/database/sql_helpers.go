package database

import (
	"context"
	"fmt"
	"strings"
)

const addPluginCommands = `INSERT INTO plugin_commands(plugin_id, command, description, created_at, updated_at)
VALUES
%s
RETURNING id, plugin_id, command, description, created_at, updated_at
`

type AddPluginCommandsParams struct {
	PluginID    int64
	Command     string
	Description string
}

// Insert and return command list
func (q *Queries) AddPluginCommands(ctx context.Context, args []AddPluginCommandsParams) (items []PluginCommand, err error) {
	commandArgs := "(?, ?, ?, datetime('now'), datetime('now'))"
	var valueStrings = []string{}
	// create a vaulues slice for easier querying
	var values = []any{} // since there are several data types
	// form arg strings and values lists
	for _, arg := range args {
		valueStrings = append(valueStrings, commandArgs)
		values = append(values,
			arg.PluginID,
			arg.Command,
			arg.Description,
		)
	}
	// create a VALUES string for query population
	commands := strings.Join(valueStrings, ",\n")

	// run query using unpacked values
	rows, err := q.db.QueryContext(ctx, fmt.Sprintf(addPluginCommands, commands), values...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// process returned items
	for rows.Next() {
		var item PluginCommand
		if err := rows.Scan(
			&item.ID,
			&item.PluginID,
			&item.Command,
			&item.Description,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, err
}
