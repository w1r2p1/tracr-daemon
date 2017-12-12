package command

import (
	log "github.com/inconshreveable/log15"
	"goku-bot/command/actions"
)

type Strategy struct {
	decisionTrees []*DecisionTree
	position      string
}

func NewStategy(position string) *Strategy {
	var trees []*DecisionTree
	return &Strategy{trees, position}
}

func (self *Strategy) AddTree(tree *DecisionTree) {
	self.decisionTrees = append(self.decisionTrees, tree)
}

func (self *Strategy) run(botActionChan chan<- *actions.ActionQueue) {
	log.Debug("running strategy", "module", "command")
	botActionQueue := actions.NewActionQueue() // the queue that will be sent back to the bot

	for _, tree := range self.decisionTrees {
		treeActionChan := make(chan *actions.ActionQueue)
		go tree.run(treeActionChan)

		treeActionQueue := <-treeActionChan // gets actions from tree

		for _, action := range treeActionQueue.Queue { // add actions from tree action queue to bot action queue
			botActionQueue.Push(action)
		}
	}

	botActionChan <- botActionQueue
}