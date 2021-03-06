import psycopg2
import os
from faker import Faker
import datetime
import random
import requests
import json

random.seed("qoinefoin")

# db_host = os.getenv("POSTGRES_HOST")
# db_pwd = os.getenv("POSTGRES_PASSWORD")
# db_port = os.getenv("POSTGRES_PORT")

# conn_dns = 'host={} user={} password={} dbname={} port={}'.format(
#     db_host, "postgres", db_pwd, "moneybot", db_port)

# conn = psycopg2.connect(conn_dns)
# cur = conn.cursor()

fake = Faker()

start_date = datetime.date(2021, 1, 1)
number_of_days = 380
date_list = [(start_date + datetime.timedelta(days=day)).isoformat()
             for day in range(number_of_days)]

# fake_name = ['default']
fake_name = []
fake_cate = ['eat', 'game', 'life']
for i in range(5):
    fake_name.append(fake.name().replace(" ", ""))
# for i in range(2):
#     fake_cate.append(fake.name().replace(" ", ""))

insert_data = []
for d in date_list:
    if random.randint(0, 10) in [3, 5, 8]:
        continue
    for i in range(random.randint(0, 3)):
        insert_data.append({
            'amount':
            random.randint(-300, 300),
            'user_id':
            'Ub44a894f0f75f644a8d77566251e67c0',
            'cate':
            random.sample(fake_name, 1)[0],
            'tags':
            random.sample(fake_name, random.randint(1, 3)),
            'date':
            d,
        })

# for i in fake_cate:
#     cur.execute(
#         f"INSERT INTO cates (name, total, user_id) VALUES ('{i}', 0, 1)")
#     conn.commit()

for i in insert_data:
    print(i)
    r = requests.post("http://localhost:9090/v1/acc/insert",
                      data=json.dumps(i))
    print(r.status_code)
    # cur.execute(
    #     f"INSERT INTO accounts (created_at, amount, user_id, cate_id) VALUES ('{i['date']}',{i['amount']}, {i['user_id']}, '{i['cate_id']}') RETURNING id"
    # )
    # id = cur.fetchone()[0]
    # # print(id)
    # for tag in i['tags']:
    #     cur.execute(
    #         f"INSERT INTO tags (created_at, name, account_id, user_id) VALUES ('{i['date']}','{tag}', {id},{i['user_id']}) RETURNING id"
    #     )
    # conn.commit()