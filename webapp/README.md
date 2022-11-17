```
docker run -p 4444:4444 -p 7900:7900 --shm-size="2g" --rm seleniarm/standalone-firefox:latest
APP_URL=http://host.docker.internal:3000 WEBDRIVER_HOST=localhost go test
```

Update `selenium.Capabilities{"browserName": "chrome/firefox"}`.
