package receivers

import (
	"goku-bot/channels"
	"github.com/gorilla/websocket"
	"net/url"
	"time"
	"errors"
	"fmt"
	log "github.com/inconshreveable/log15"
)

type Receiver interface {
	Start()
	Key() string
}

type AllReceivers map[string]Receiver

var receivers AllReceivers

type SortedReceiverKeys []string

func (a SortedReceiverKeys) Len() int           { return len(a) }
func (a SortedReceiverKeys) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortedReceiverKeys) Less(i, j int) bool { return a[i] < a[j] }

func Init() {
	receivers = make(AllReceivers)

	obr := NewOrderBookReceiver("poloniex", "USDT_BTC")
	receivers[obr.Key()] = obr

	tr := NewTickerReceiver("poloniex", "USDT_BTC")
	receivers[tr.Key()] = tr

}

func Start() error {
	for _, receiver := range receivers {
		log.Info("Starting receiver", "key", receiver.Key(), "module", "receivers")
		go receiver.Start()
	}

	return nil
}

func sendToProcessor(key string, output interface{}) {
	log.Debug("sending to processor module", "key", key, "module", "receivers")

	processorInput := channels.ReceiverOutputProcessorInput{output, key}

	channels.ReceiverProcessorChan <- processorInput
}

func websocketConnect(address string, retries int) (*websocket.Conn, error) {
	u := url.URL{Scheme: "wss", Host: address}

	var connection *websocket.Conn
	var err error
	retriesLeft := retries
	for retriesLeft > 0 {
		connection, _, err = websocket.DefaultDialer.Dial(u.String(), nil)

		if err != nil {
			log.Warn("error connecting - retrying after 5 seconds", "module", "receivers", "error", err, "retriesLeft", retriesLeft)
			retriesLeft--

			timer := time.NewTimer(time.Second * 5)
			<-timer.C
		} else {
			break
		}
	}

	if retriesLeft == 0 {
		return connection, errors.New(fmt.Sprintf("Could not connect after %d attemps", retries))
	}

	return connection, nil
}