package main

import (
    "bytes"
    "encoding/json"
    "io"
    "log"
    "net/http"
    "strconv"

    "github.com/JouleJ/WLB0/core"
    "github.com/JouleJ/WLB0/internal"
    "github.com/go-chi/chi/v5"
)

func getDatabaseConnection() (core.DatabaseConnection, error) {
    db, err := internal.ConnectPostgres()
    if err != nil {
        return nil, err
    }

    return internal.NewCachedConnection(db), nil
}

func getDeliverById(db core.DatabaseConnection, id int) interface{} {
    d, err := db.LoadDeliver(id)
    if err != nil {
        log.Printf("Failed to load deliver: %v\n", err)
        return nil
    }

    return &d
}

func getPaymentById(db core.DatabaseConnection, id int) interface{} {
    p, err := db.LoadPayment(id)
    if err != nil {
        log.Printf("Failed to load payment: %v\n", err)
        return nil
    }

    return &p
}

func getItemById(db core.DatabaseConnection, id int) interface{} {
    i, err := db.LoadItem(id)
    if err != nil {
        log.Printf("Failed to load item: %v\n", err)
        return nil
    }

    return &i
}

func getOrderById(db core.DatabaseConnection, id int) interface{} {
    o, err := db.LoadOrder(id)
    if err != nil {
        log.Printf("Failed to load order: %v\n", err)
        return nil
    }

    return &o
}

func makeChiHandler(f func(db core.DatabaseConnection, id int) interface{}) func(w http.ResponseWriter, r *http.Request) {
    return func (w http.ResponseWriter, r *http.Request) {
        id, err := strconv.Atoi(r.URL.Query().Get("id"))
        if err != nil {
            log.Printf("Invalid id: %v\n", err)
            io.WriteString(w, "Invalid id")
            return
        }

        db, err := getDatabaseConnection()
        if err != nil {
            log.Printf("Failed to connect to database: %v\n", err)
            io.WriteString(w, "Failed to connect to database")
            return
        }

        defer db.Close()

        value := f(db, id)
        jsonBytes, err := json.MarshalIndent(value, "", "    ")
        if err != nil {
            log.Printf("Failed to marshal: %v\n", err)
            io.WriteString(w, "Failed to marshal")
            return
        }

        var b bytes.Buffer
        json.HTMLEscape(&b, jsonBytes)
        io.WriteString(w, `<!DOCTYPE HTML>`)
        io.WriteString(w, `<html>`)

        io.WriteString(w, `<head>`)
        io.WriteString(w, `<style>`)

        io.WriteString(w, `body {
                             font-family: Arial, sans-serif;
                             background-color: #f2f2f2;
                             padding: 20px;
                           }`)

        io.WriteString(w, `div {
                             background-color: #fff;
                             border-radius: 5px;
                             padding: 20px;
                             box-shadow: 0px 0px 10px rgba(0,0,0,0.3);
                           }`)

        io.WriteString(w, `</style>`)
        io.WriteString(w, `</head>`)

        io.WriteString(w, `<body>`)
        io.WriteString(w, `<div>`)
        io.WriteString(w, `<pre>`)
        b.WriteTo(w)
        io.WriteString(w, `</pre>`)
        io.WriteString(w, `</div>`)
        io.WriteString(w, `</body>`)

        io.WriteString(w, `</html>`)
    }
}

func main() {
    var stream core.StreamConnection
    stream, err := internal.ConnectNats()
    if err != nil {
        log.Fatalf("Failed to connect to NATS: %v\n", err)
    }

    stream.SubscribeForDelivers(func (d core.Deliver) {
        db, err := getDatabaseConnection()
        if err != nil {
            log.Printf("Failed to open database connection: %v\n", err)
            return
        }

        defer db.Close()

        _, err = db.InsertDeliver(d)
        if err != nil {
            log.Printf("Failed insert into database: %v\n", err)
        }
    })

    stream.SubscribeForPayments(func (p core.Payment) {
        db, err := getDatabaseConnection()
        if err != nil {
            log.Printf("Failed to open database connection: %v\n", err)
            return
        }

        defer db.Close()

        _, err = db.InsertPayment(p)
        if err != nil {
            log.Printf("Failed to insert into database: %v\n", err)
        }
    })

    stream.SubscribeForItems(func (i core.Item) {
        db, err := getDatabaseConnection()
        if err != nil {
            log.Printf("Failed to open database connection: %v\n", err)
            return
        }

        defer db.Close()

        _, err = db.InsertItem(i)
        if err != nil {
            log.Printf("Failed to insert into database: %v\n", err)
        }
    })

    stream.SubscribeForOrders(func (o core.Order) {
        db, err := getDatabaseConnection()
        if err != nil {
            log.Printf("Failed to open database connection: %v\n", err)
            return
        }

        defer db.Close()

        _, err = db.InsertOrder(o)
        if err != nil {
            log.Printf("Failed to insert into database: %v\n", err)
        }
    })

    r := chi.NewRouter()
    r.Get("/deliver", makeChiHandler(getDeliverById))
    r.Get("/payment", makeChiHandler(getPaymentById))
    r.Get("/item", makeChiHandler(getItemById))
    r.Get("/order", makeChiHandler(getOrderById))

    http.ListenAndServe(":80", r)
}
