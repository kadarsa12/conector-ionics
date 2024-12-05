## BUILD:

```bash
GOOS=windows GOARCH=amd64 go build .
```

```
sc create IonicsBIService binPath="C:\path\to\ionicsbi-service.exe -db_name=pgsql -db_host=host -db_sid=sid -db_username=dbuser -db_password=dbpass base_url=baseURL -execution_interval=0 batch_size=100 initial_date=YYYY/MM/DD -client_id=clientId -client_secret=clientSecret"

sc start IonicsBIService
sc query IonicsBIService
sc stop IonicsBIService
```
