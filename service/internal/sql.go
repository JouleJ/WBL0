package internal

import (
    "database/sql"
    "fmt"
    "os"

    "github.com/JouleJ/WLB0/core"
    "github.com/lib/pq"
)

type pgConnection struct {
    impl *sql.DB
}

func ConnectPostgres() (core.DatabaseConnection, error) {
    dbImpl, err := sql.Open(
        "postgres",
        fmt.Sprintf(
            "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            os.Getenv("PG_HOST"),
            os.Getenv("PG_PORT"),
            os.Getenv("PG_USER"),
            os.Getenv("PG_PASSWORD"),
            os.Getenv("PG_NAME")))

    if err != nil {
        return nil, err
    }

    return &pgConnection{impl: dbImpl}, nil
}

func (db *pgConnection) LoadDeliver(id int) (core.Deliver, error) {
    var d core.Deliver

    rows, err := db.impl.Query(
        `SELECT name, phone, zip, city, address, region, email FROM delivers WHERE id = $1;`,
        id)

    if err != nil {
        return d, err
    }

    if rows.Next() {
        rows.Scan(
            &d.Name,
            &d.Phone,
            &d.Zip,
            &d.City,
            &d.Address,
            &d.Region,
            &d.Email)

        return d, nil
    } else {
        return d, fmt.Errorf("No such id=%v", id)
    }
}

func (db *pgConnection) LoadPayment(id int) (core.Payment, error) {
    var p core.Payment

    rows, err := db.impl.Query(
        `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE id = $1;`,
        id)

    if err != nil {
        return p, err
    }

    if rows.Next() {
        rows.Scan(
            &p.Transaction,
            &p.RequestId,
            &p.Currency,
            &p.Provider,
            &p.Amount,
            &p.PaymentDt,
            &p.Bank,
            &p.DeliveryCost,
            &p.GoodsTotal,
            &p.CustomFee)

        return p, nil
    } else {
        return p, fmt.Errorf("No such id=%v", id)
    }
}

func (db *pgConnection) LoadItem(id int) (core.Item, error) {
    var i core.Item

    rows, err := db.impl.Query(
        `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM items WHERE id = $1;`,
        id)

    if err != nil {
        return i, err
    }

    if rows.Next() {
        rows.Scan(
            &i.ChrtId,
            &i.TrackNumber,
            &i.Price,
            &i.Rid,
            &i.Name,
            &i.Sale,
            &i.Size,
            &i.TotalPrice,
            &i.NmId,
            &i.Brand,
            &i.Status)

        return i, nil
    } else {
        return i, fmt.Errorf("No such id=%v\n", id)
    }
}

func (db *pgConnection) LoadOrder(id int) (core.Order, error) {
    var o core.Order

    rows, err := db.impl.Query(
        `SELECT order_uid,
                track_number,
                entry,
                deliver_id,
                payment_id,
                item_ids,
                locale,
                internal_signature,
                customer_id,
                delivery_service,
                shardkey,
                sm_id,
                date_created,
                oof_shard
         FROM orders
         WHERE id = $1`,
        id)

    if err != nil {
        return o, err
    }

    if rows.Next() {
        var deliverId, paymentId int
        var itemIds pq.Int32Array

        rows.Scan(
            &o.OrderUid,
            &o.TrackNumber,
            &o.Entry,
            &deliverId,
            &paymentId,
            &itemIds,
            &o.Locale,
            &o.InternalSignature,
            &o.CustomerId,
            &o.DeliveryService,
            &o.ShardKey,
            &o.SmId,
            &o.DateCreated,
            &o.OofShard)

        o.Deliver, err = db.LoadDeliver(deliverId)
        if err != nil {
            return o, err
        }

        o.Payment, err = db.LoadPayment(paymentId)
        if err != nil {
            return o, err
        }

        for _, itemId := range itemIds {
            item, err := db.LoadItem(int(itemId))
            if err != nil {
                return o, err
            }

            o.Items = append(o.Items, item)
        }

        return o, nil
    } else {
        return o, fmt.Errorf("No such id=%v\n", id)
    }
}

func (db *pgConnection) InsertDeliver(d core.Deliver) (int, error) {
    id := -1

    rows, err := db.impl.Query(
        `INSERT INTO delivers (name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id;`,
        d.Name,
        d.Phone,
        d.Zip,
        d.City,
        d.Address,
        d.Region,
        d.Email)

    if err != nil {
        return id, err
    }

    if rows.Next() {
        rows.Scan(&id)
        return id, nil
    } else {
        return id, fmt.Errorf("Failed to insert")
    }
}

func (db *pgConnection) InsertPayment(p core.Payment) (int, error) {
    id := -1

    rows, err := db.impl.Query(
        `INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
         RETURNING id;`,
         p.Transaction,
         p.RequestId,
         p.Currency,
         p.Provider,
         p.Amount,
         p.PaymentDt,
         p.Bank,
         p.DeliveryCost,
         p.GoodsTotal,
         p.CustomFee)

    if err != nil {
        return id, err
    }

    if rows.Next() {
        rows.Scan(&id)
        return id, nil
    } else {
        return id, fmt.Errorf("Failed to insert")
    }
}

func (db *pgConnection) InsertItem(i core.Item) (int, error) {
    id := -1

    rows, err := db.impl.Query(
        `INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
         RETURNING id;`,
         i.ChrtId,
         i.TrackNumber,
         i.Price,
         i.Rid,
         i.Name,
         i.Sale,
         i.Size,
         i.TotalPrice,
         i.NmId,
         i.Brand,
         i.Status)

    if err != nil {
        return id, err
    }

    if rows.Next() {
        rows.Scan(&id)
        return id, nil
    } else {
        return id, fmt.Errorf("Failed to insert")
    }
}

func (db *pgConnection) InsertOrder(o core.Order) (int, error) {
    id := -1

    deliverId, err := db.InsertDeliver(o.Deliver)
    if err != nil {
        return id, err
    }

    paymentId, err := db.InsertPayment(o.Payment)
    if err != nil {
        return id, err
    }

    itemIds := pq.Int32Array{}
    for _, i := range o.Items {
        itemId, err := db.InsertItem(i)
        if err != nil {
            return id, err
        }

        itemIds = append(itemIds, int32(itemId))
    }

    rows, err := db.impl.Query(
        `INSERT INTO orders (order_uid,
                             track_number,
                             entry,
                             deliver_id,
                             payment_id,
                             item_ids,
                             locale,
                             internal_signature,
                             customer_id,
                             delivery_service,
                             shardkey,
                             sm_id,
                             date_created,
                             oof_shard)
         VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
         RETURNING id;`,
        o.OrderUid,
        o.TrackNumber,
        o.Entry,
        deliverId,
        paymentId,
        itemIds,
        o.Locale,
        o.InternalSignature,
        o.CustomerId,
        o.DeliveryService,
        o.ShardKey,
        o.SmId,
        o.DateCreated,
        o.OofShard)

    if err != nil {
        return id, err
    }

    if rows.Next() {
        rows.Scan(&id)
        return id, nil
    } else {
        return id, fmt.Errorf("Failed to insert")
    }
}

func (db *pgConnection) Ping() error {
    return db.impl.Ping()
}

func (db *pgConnection) Close() {
    db.impl.Close()
}
