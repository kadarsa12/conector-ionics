## BUILD:

```bash
GOOS=windows GOARCH=amd64 go build .
```

```
sc create InovaKPIService binPath="C:\path\to\inovakpi-service.exe -db_name=pgsql -db_host=host -db_sid=sid -db_username=dbuser -db_password=dbpass base_url=baseURL -execution_interval=0 batch_size=100 initial_date=YYYY/MM/DD -clientId=clientId -clientSecret=clientSecret"

sc start InovaKPIService
sc query InovaKPIService
sc stop InovaKPIService
```
