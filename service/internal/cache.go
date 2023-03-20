package internal

import (
    "github.com/JouleJ/WLB0/core"
)

var (
    delivers = map[int]core.Deliver{}
    payments = map[int]core.Payment{}
    items = map[int]core.Item{}
    orders = map[int]core.Order{}
)

type cache struct {
    db core.DatabaseConnection
}

func NewCachedConnection(db core.DatabaseConnection) core.DatabaseConnection {
    return &cache{db: db}
}

func (c *cache) LoadDeliver(id int) (core.Deliver, error) {
    var d core.Deliver

    d, ok := delivers[id]
    if ok {
        return d, nil
    }

    d, err := c.db.LoadDeliver(id)
    if err != nil {
        return d, err
    }

    delivers[id] = d
    return d, err
}

func (c *cache) LoadPayment(id int) (core.Payment, error) {
    var p core.Payment

    p, ok := payments[id]
    if ok {
        return p, nil
    }

    p, err := c.db.LoadPayment(id)
    if err != nil {
        return p, err
    }

    payments[id] = p
    return p, err
}

func (c *cache) LoadItem(id int) (core.Item, error) {
    var i core.Item

    i, ok := items[id]
    if ok {
        return i, nil
    }

    i, err := c.db.LoadItem(id)
    if err != nil {
        return i, err
    }

    items[id] = i
    return i, err
}

func (c *cache) LoadOrder(id int) (core.Order, error) {
    var o core.Order

    o, ok := orders[id]
    if ok {
        return o, nil
    }

    o, err := c.db.LoadOrder(id)
    if err != nil {
        return o, err
    }

    orders[id] = o
    return o, err
}

func (c *cache) InsertDeliver(d core.Deliver) (int, error) {
    id, err := c.db.InsertDeliver(d)
    if err != nil {
        return id, err
    }

    delivers[id] = d
    return id, err
}

func (c *cache) InsertPayment(p core.Payment) (int, error) {
    id, err := c.db.InsertPayment(p)
    if err != nil {
        return id, err
    }

    payments[id] = p
    return id, err
}

func (c *cache) InsertItem(i core.Item) (int, error) {
    id, err := c.db.InsertItem(i)
    if err != nil {
        return id, err
    }

    items[id] = i
    return id, err
}

func (c *cache) InsertOrder(o core.Order) (int, error) {
    id, err := c.db.InsertOrder(o)
    if err != nil {
        return id, err
    }

    orders[id] = o
    return id, err
}

func (c *cache) Ping() error {
    return c.db.Ping()
}

func (c *cache) Close() {
    c.db.Close()
}
