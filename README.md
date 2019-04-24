
LineBot with Kafka Flow:

![Screenshot](flow.png)

Update line secret and access token in docker-compose.yml

```
  linewebhook:
    .
    .
    .
    environment:
        CHANNEL_SECRET: <your secret> 
        CHANNEL_ACCESS_TOKEN: <your access token>
```

RUN: 
```
    docker-compose up -d
```
![Screenshot](docker-compose-up.png)
```
    docker container ls
```
![Screenshot](docker_screenshot.png)


webhook listen port 8080