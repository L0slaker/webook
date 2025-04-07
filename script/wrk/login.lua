wrk.method="POST"
wrk.headers["Content-Type"] = "application/json"
-- 这个要改为你的注册的数据
wrk.body='{"email":"user999004@example.com", "password": "hello@world123"}'