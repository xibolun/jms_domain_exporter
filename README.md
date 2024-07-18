# jms_domain_exporter

## Quick start
```
curl -sSL https://github.com/xibolun/jms_domain_exporter/install.sh | bash
```

## Build By Yourself
1. `make build` you will get a binary file jms_domain_exporter.
2. just replace conf.yml jms_addr and jms_token
3. start server by `jms_domain_exporter -c conf.yml`
4. you can access the metrics at `http://localhost:8080/metrics`