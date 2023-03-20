package core

type DatabaseConnection interface {
    LoadDeliver(id int) (Deliver, error)
    LoadPayment(id int) (Payment, error)
    LoadItem(id int) (Item, error)
    LoadOrder(id int) (Order, error)

    InsertDeliver(d Deliver) (int, error)
    InsertPayment(p Payment) (int, error)
    InsertItem(i Item) (int, error)
    InsertOrder(o Order) (int, error)

    Ping() error
    Close()
}
