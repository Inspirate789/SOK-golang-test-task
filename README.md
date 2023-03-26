# SOK-golang-test-task
Simple transaction server based on Golang web server, PostgreSQL and Apache Kafka

### Run with docker compose:

```bash
CUR_IP=$(nmcli device show | grep IP4.ADDRESS | head -1 | awk '{print $2}' | rev | cut -c 4- | rev)
HOST_IP=$CUR_IP docker-compose -f docker-compose-kafka.yml up -d
HOST_IP=$CUR_IP docker-compose up -d --build
```

It will build and run the API and 2 Worker inside the docker.