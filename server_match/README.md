# SERVER MATCH

The matching server consumes orders from Kafka for matching. With 1 CPU and 128MB of memory, it can match 1000+ orders per second.

## Build Docker Image

```bash
docker build . -t market/match
```
