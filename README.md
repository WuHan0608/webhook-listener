Webhook HTTP Listener
====

### Installation
```
# mkdir -p $GOPATH/src/github.com/WuHan0608/
# cd $GOPATH/src/github.com/WuHan0608/
# git clone https://github.com/WuHan0608/webhook-listener.git
# cd webhook-listener
# CGO_ENABLED=0 GOOS=linux go install -ldflags="-s -w"
```

### Systemd Service
- create systemd unit file: /etc/systemd/system/webhook-listener.service
  - LISTEN_ADDRESS: http server listen address, defaults to :80
  - TRACK_REPO_TAG: track the specified repository tag only, defaults to all tags
  - PUSHOVER_USER_KEYS: pushover user keys seperated by comma or semicolon
  - PUSHOVER_API_TOKEN: pushover api token

```
[Unit]
Description=Webhook HTTP Server
After=network.target

[Service]
User=centos
Environment=LISTEN_ADDRESS=:8001
Environment=TRACK_REPO_TAG="latest"
Environment=PUSHOVER_USER_KEYS=userKey1,userKey2,groupKey3,groupKey4
Environment=PUSHOVER_API_TOKEN=apiToken
ExecStart=/home/centos/gocode/bin/webhook-listener

[Install]
WantedBy=multi-user.target
```

- start systemd service
```
# sudo systemctl start webhook-listener
# sudo systemctl enable webhook-listener
```

### Docker container
- pull image
```
# docker pull wuhan/webhook-listener:latest
```

- start container
```
# docker run -d --restart=always --name=webhook-listener -p 8000:8000 -e LISTEN_ADDRESS=":8000" -e TRACK_REPO_TAG="latest" -e PUSHOVER_USER_KEYS="userKey1,userKey2,groupKey3,groupKey4" -e PUSHOVER_API_TOKEN="apiToken" wuhan/webhook-listener:latest
```

### Integrated Webhook Services
- DockerHub Automated Build Webhook
  - webhook url: `http://<public_ip>:<published_port>/webhook/dockerhub`
  - health status: `curl -s http://<public_ip>:<published_port>/webhook/dockerhub`

