# Messagebird Proxy API

## How to run? 
- Image is available on Docker Hub: `https://hub.docker.com/r/jazzerjazzer/messagebird_proxy/`

- Pull the latest:
```
docker pull jazzerjazzer/messagebird_proxy
```

- Run the docker image with API_KEY variable:
```
docker run -e API_KEY=<API_KEY> --expose 8080 -p 8080:8080 -t jazzerjazzer/messagebird_proxy
``` 

- Docker container should be running on port :8080. Requests now can be sent to: 
``` 
http://localhost:8080/sendMessage
```

## API

API docs can be found in `/docs/swagger.yaml`

## Libraries used
- Phone number validation: `https://github.com/nyaruka/phonenumbers`
- Unit testing: `https://github.com/stretchr/testify`
- Mocking: `https://github.com/vektra/mockery`
- Messagebird Go Client: `https://github.com/messagebird/go-rest-api`

## Tests
```
go test ./src/app/...
```
## Known Issues

- Sending multipart Unicode datacoded messages: Unicode characters are swallowed when the datacoding of messages are set to `auto` or `plain`. If the datacoding is set to `unicode`, then the message is displayed with the binary body.