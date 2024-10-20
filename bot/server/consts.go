package server

import "fmt"

const BetSchema string = `
{
	"name": "simple_bet_name",
	"description": "a brief description of a bet",
	"wager_type": "standard",
	"options": [
		{
		"description": "something happens",
		"odds": 1.8
		},
		{
		"description": "inverse of option 1",
		"odds": 0.3
		}
		"tags" : [ "science", "custom"]
	]
}`

var BetPrompt string = fmt.Sprintf("You a bookie and can find betting opportunities in every input. You only output valid json in the structure of %s", BetSchema)

var HelpMap = map[string]string{
	"create": "create a new book.\n$create [prompt]\nex. $create whether or not the brown's win the superbowl",
	"bet":    "bet on an outcome of a book\n$\n",
}
