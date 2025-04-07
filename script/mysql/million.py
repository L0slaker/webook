import mysql.connector
import random
from datetime import datetime, timedelta
import time

# 数据库配置
db_config = {
    'host': 'localhost',
    'port': 13316,
    'user': 'root',
    'password': 'root',
    'database': 'webook'
}

# 中文名字生成
first_names = ["张", "李", "王", "赵", "刘", "陈", "杨", "黄", "吴", "周"]
last_names = ["伟", "芳", "娜", "秀英", "敏", "静", "丽", "强", "磊", "军"]
hobbies = ["游泳", "跑步", "篮球", "足球", "网球", "乒乓球", "羽毛球", "围棋", "象棋", "阅读"]

def generate_user(id):
    """生成一个用户数据"""
    first_name = random.choice(first_names)
    last_name = random.choice(last_names)
    nickname = first_name + last_name
    email = f"user{id}@example.com"

    # 生成生日(1950-1999)
    birth_year = random.randint(1950, 1999)
    birth_month = random.randint(1, 12)
    birth_day = random.randint(1, 28)
    birthday = f"{birth_year}-{birth_month:02d}-{birth_day:02d}"

    # 生成个人介绍
    hobby1 = random.choice(hobbies)
    hobby2 = random.choice(hobbies)
    while hobby2 == hobby1:
        hobby2 = random.choice(hobbies)
    introduction = f"我喜欢{hobby1}和{hobby2}"

    # 时间戳(秒)
    now = int(time.time())

    return (
        email,
        "$2a$10$" + ''.join(random.choices('abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789', k=50)),
        nickname,
        birthday,
        introduction,
        now,
        now
    )

def insert_users(total_records, batch_size=1000):
    """插入用户数据"""
    try:
        # 连接数据库
        db = mysql.connector.connect(**db_config)
        cursor = db.cursor()

        # 创建表(如果不存在)
        create_table_sql = """
        CREATE TABLE IF NOT EXISTS users (
            id BIGINT AUTO_INCREMENT PRIMARY KEY,
            email VARCHAR(255) NOT NULL UNIQUE,
            password VARCHAR(255) NOT NULL,
            nickname VARCHAR(100),
            birthday VARCHAR(20),
            introduction TEXT,
            created_at BIGINT,
            updated_at BIGINT,
            INDEX idx_email (email)
        ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
        """
        cursor.execute(create_table_sql)

        # 清空表(可选)
        try:
            cursor.execute("TRUNCATE TABLE users")
        except mysql.connector.Error as err:
            print(f"清空表时出错: {err}")

        # 准备插入SQL
        insert_sql = """
        INSERT INTO users
        (email, password, nickname, birthday, introduction, created_at, updated_at)
        VALUES (%s, %s, %s, %s, %s, %s, %s)
        """

        print(f"开始插入{total_records}条数据...")
        start_time = time.time()

        # 分批插入
        for i in range(0, total_records, batch_size):
            batch = []
            current_batch_size = min(batch_size, total_records - i)

            # 生成当前批次的数据
            for j in range(1, current_batch_size + 1):
                user_id = i + j
                batch.append(generate_user(user_id))

            # 执行批量插入
            cursor.executemany(insert_sql, batch)
            db.commit()

            # 打印进度
            if (i // batch_size) % 10 == 0:
                elapsed = time.time() - start_time
                print(f"已插入{i + current_batch_size}条记录, 耗时: {elapsed:.2f}秒")

        total_time = time.time() - start_time
        print(f"完成! 共插入{total_records}条记录, 总耗时: {total_time:.2f}秒")

    except mysql.connector.Error as err:
        print(f"数据库错误: {err}")
    finally:
        if 'db' in locals() and db.is_connected():
            cursor.close()
            db.close()

if __name__ == "__main__":
    # 插入100万条记录，每批1000条
    insert_users(1000000, batch_size=1000)