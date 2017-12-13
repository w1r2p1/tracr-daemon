package processors

import (
	log "github.com/inconshreveable/log15"
	"tracr-daemon/exchanges"
	"fmt"
	"tracr-store"
)

type MyTradeHistoryProcessor struct {
	exchange string
	Store    tracr_store.Store
}

func NewMyTradeHistoryProcessor(exchange string) *MyTradeHistoryProcessor {
	s, _ := tracr_store.NewStore()

	return &MyTradeHistoryProcessor{exchange, s}
}

func (self *MyTradeHistoryProcessor) Process(input interface{}) {
	log.Debug("processing", "key", self.Key(), "module", "processors")
	data := input.(map[string][]exchanges.Trade)

	for pair, trades := range data {
		log.Debug("replacing trades for pair", "module", "processors", "pair", pair, "numOfTrades", len(trades))
		self.Store.ReplaceTrades(self.exchange, pair, trades)
	}
}

func (self *MyTradeHistoryProcessor) Key() string {
	return fmt.Sprintf("MyTradeHistory-%s", self.exchange)
}
