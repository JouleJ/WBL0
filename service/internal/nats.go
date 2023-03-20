package internal

import (
    "encoding/json"
    "log"
    "os"

    "github.com/nats-io/nats.go"
    "github.com/JouleJ/WLB0/core"
)

type natsConnection struct {
    impl *nats.Conn
}

func ConnectNats() (core.StreamConnection, error) {
    streamImpl, err := nats.Connect(os.Getenv("NATS_URI"))
    if err != nil {
        return nil, err
    }

    log.Printf("Connection to NATS: %v\n", *streamImpl)

    return &natsConnection{impl: streamImpl}, nil
}

func (nc *natsConnection) SubscribeForDelivers(handler func (d core.Deliver)) {
    _, err := nc.impl.Subscribe("delivers", func (m *nats.Msg) {
        var d core.Deliver

        err := json.Unmarshal(m.Data, &d)
        if err == nil {
            handler(d)
        } else {
            log.Printf("Failed to unmarshal data=%v due to err=%v\n", m.Data, err)
        }
    })

    if err != nil {
        log.Printf("Failed to subscribe to delivers: %v\n", err)
    }
}

func (nc *natsConnection) SubscribeForPayments(hanlder func (p core.Payment)) {
    _, err := nc.impl.Subscribe("payments", func(m *nats.Msg) {
        var p core.Payment

        err := json.Unmarshal(m.Data, &p)
        if err == nil {
            hanlder(p)
        } else {
            log.Printf("Failed to unmarshal data=%v due to err=%v\n", m.Data, err)
        }
    })

    if err != nil {
        log.Printf("Failed to subscribe to payments: %v\n", err)
    }
}

func (nc *natsConnection) SubscribeForItems(handler func (i core.Item)) {
    _, err := nc.impl.Subscribe("items", func(m *nats.Msg) {
        var i core.Item

        err := json.Unmarshal(m.Data, &i)
        if err == nil {
            handler(i)
        } else {
            log.Printf("Failed to unmarshal data=%v due to err=%v\n", m.Data, err)
        }
    })

    if err != nil {
        log.Printf("Failed to subscribe to items: %v\n", err)
    }
}

func (nc *natsConnection) SubscribeForOrders(handler func (o core.Order)) {
    _, err := nc.impl.Subscribe("orders", func(m *nats.Msg) {
        var o core.Order

        err := json.Unmarshal(m.Data, &o)
        if err == nil {
            handler(o)
        } else {
            log.Printf("Failed to unmarshal data=%v due to err=%v\n", m.Data, err)
        }
    })

    if err != nil {
        log.Printf("Failed to subscribe to orders: %v\n", err)
    }
}
