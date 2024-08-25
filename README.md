# Stress Test

## Description

This is a stress test for the http server. It will send a lot of requests to the server and check if the server can handle it. In the end, it will create a report with the results.

## Docker
### Build
```shell
docker pull mayckol/loadtester
```

### Run
```shell
docker run mayckol/loadtester --url=http://example.com --requests=10 --concurrency=5
```
**Parameters:**
- `url`: The url to test
- `requests`: The number of requests to send
- `concurrency`: The number of requests to send concurrently

### Reports sample

**Success:**
```
🏁 Starting the load test for http://example.com...

===== 📝 Load Test Report =====
⏳ Total time: 510.069834ms
📊 Total requests: 10
✅ Successful requests (HTTP 200): 10

⚡ Requests per second: 19.61
```

**Failure:**
```
🏁 Starting the load test for http://example.com/404...

===== 📝 Load Test Report =====
⏳ Total time: 511.376792ms
📊 Total requests: 10
✅ Successful requests (HTTP 200): 0

📉 Distribution of other HTTP status codes:
  ❌ Failed requests (HTTP 404): 10
```