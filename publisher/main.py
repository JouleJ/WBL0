import asyncio
import json
import nats
import os
import random
import time

def randhex():
    return hex(random.randint(0, 1e18))[2:]

def random_deliver():
    return dict(
        name='name_{}'.format(randhex()),
        phone='phone_{}'.format(randhex()),
        zip='zip_{}'.format(randhex()),
        city='city_{}'.format(randhex()),
        address='address_{}'.format(randhex()),
        region='region_{}'.format(randhex()),
        email='email_{}'.format(randhex()),
    )

def random_payment():
    return dict(
        transaction='transaction_{}'.format(randhex()),
        request_id='request_id_{}'.format(randhex()),
        currency='currency_{}'.format(randhex()),
        provider='provider={}'.format(randhex()),
        amount=random.randint(0, 10000),
        payment_dt=random.randint(0, 10000),
        bank='bank_{}'.format(randhex()),
        delivery_cost=random.randint(0, 10000),
        goods_total=random.randint(0, 10000),
        custom_fee=random.randint(0, 10000),
    )

def random_item():
    return dict(
        chrt_id=random.randint(0, 10000),
        track_number='track_number_{}'.format(randhex()),
        price=random.randint(0, 10000),
        rid='rid_{}'.format(randhex()),
        name='name_{}'.format(randhex()),
        sale=random.randint(0, 10000),
        size='size_{}'.format(randhex()),
        total_price=random.randint(0, 10000),
        nm_id=random.randint(0, 10000),
        brand='brand_{}'.format(randhex()),
        status=random.randint(0, 10000),
    )

def random_items():
    n = random.randint(1, 5)
    return [random_item() for _ in range(n)]

def random_order():
    return dict(
        order_uid='order_uid_{}'.format(randhex()),
        track_number='track_number_{}'.format(randhex()),
        entry='entry_{}'.format(randhex()),
        delivery=random_deliver(),
        payment=random_payment(),
        items=random_items(),
        locale='locale_{}'.format(randhex()),
        internal_signature='internal_signature_{}'.format(randhex()),
        customer_id='customer_id_{}'.format(randhex()),
        delivery_service='delivery_service_{}'.format(randhex()),
        shardkey='shardkey_{}'.format(randhex()),
        sm_id=random.randint(0, 10000),
        date_created='date_created_{}'.format(randhex()),
        oof_shard='oof_shard_{}'.format(randhex()),
    )

async def main():
    connection = await nats.connect(os.getenv('NATS_URI'))

    await connection.publish('delivers', json.dumps(random_deliver()).encode('utf-8'))
    await connection.publish('payments', json.dumps(random_payment()).encode('utf-8'))
    await connection.publish('items', json.dumps(random_item()).encode('utf-8'))
    await connection.publish('orders', json.dumps(random_order()).encode('utf-8'))

    await connection.drain()

if __name__ == '__main__':
    while True:
        time.sleep(5)
        asyncio.run(main())
