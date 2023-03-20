package core

type StreamConnection interface {
    SubscribeForDelivers(handler func (d Deliver))
    SubscribeForPayments(hanlder func (p Payment))
    SubscribeForItems(hanlder func (i Item))
    SubscribeForOrders(handler func (o Order))
}
