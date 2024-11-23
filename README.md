## BUILD:

```bash
GOOS=windows GOARCH=amd64 go build .
```

```
sc create InovaKPIService binPath="C:\path\to\inovakpi-service.exe -db_name=pgsql -db_host=host -db_sid=sid -db_username=dbuser -db_password=dbpass -execution_interval=0 -clientId=clientId -clientSecret=clientSecret"

sc start InovaKPIService
sc query InovaKPIService
sc stop InovaKPIService
```
