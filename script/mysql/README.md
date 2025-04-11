## 生成百万条数据如何执行脚本
1. 安装依赖：
   首先确保已安装mysql-connector-python：
    ```c
    pip install mysql-connector-python
    ```
2. 修改配置：
   编辑脚本中的db_config部分，填入你的MySQL连接信息：
    ```c
    db_config = {
    'host': 'localhost',
    'user': 'yourusername',
    'password': 'yourpassword',
    'database': 'yourdatabase'
    }
    ```
3. 运行脚本
    ```c
    python3 million.py
    ```