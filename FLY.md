# fly

## database
```bash
fly postgres create
```

it will return the credentials like:
```bash
Postgres cluster app-name created
  Username:    **** 
  Password:    **** 
  Hostname:    **** 
  Flycast:     **** 
  Proxy port:  ****
  Postgres port:  **** 
  Connection string: ****
```

you can connect to it with the command:
```bash
fly postgres connect -a app-name
```

once inside you have to create the db tables, copy-paste the content of 
the file tables.sql.


## backend

create the x73-app:
```bash
fly apps create n73-app
```

create the volume for the app:
```bash
fly volumes create go_data --size 5 --region arn --app n73-app
```

create the fly token:
```bash
fly tokens create org --name "x73" --expiry 720h
```

add the require enviroment variables
```bash
fly secrets set DB_USER=***** --app n73-app
fly secrets set DB_PASSWORD=***** --app n73-app
fly secrets set DB_HOST=***** --app n73-app
fly secrets set DB_PORT=***** --app n73-app
fly secrets set DB_NAME=***** --app n73-app

fly secrets set EMAIL_SECRET_KEY=***** --app n73-app
fly secrets set PORT=**** --app n73-app
fly secrets set ROOT_PATH=**** --app n73-app
fly secrets set SECRET_KEY=***** --app n73-app
fly secrets set ADMIN_JWT=***** --app n73-app
fly secrets set DOMAIN=***** --app n73-app

fly secrets set FLY_API_TOKEN="****" --app x73-app
fly secrets set CLOUDFLARE_API_TOKEN="****" --app n73-app
fly secrets set GH_TOKEN="****" --app n73-app
```

deploy the app(ha==high avilability, just 1 machine)
```bash
fly deploy --ha=false
```

now the last step, create a new account and update the account to be admin:
```sql
UPDATE users SET role = 'admin' WHERE email = 'agustfricke@gmail.com';
```

do:
mkdir /data/fly_configs && /data/projects
