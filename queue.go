package library

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
)

func PublishToNats(conn *nats.Conn,service, name string, payload interface{}) error {

	js, _ := json.Marshal(payload)
	return conn.Publish(fmt.Sprintf("da.mno.%s.%s",service,name),js)
}
