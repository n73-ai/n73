# notes

error uploading directory: 
error uploading assets/Geist-Black-DyTs4Xsi.woff2: status 401, 
body: {"HttpCode":401,"Message":"Unauthorized"}

failed to create Pull Zone (status 400): 
{"ErrorKey":"pullzone.validation","Field":"",
"Message":"Object reference not set to an instance of an object."}

failed to create Pull Zone (status 400): 
{"ErrorKey":"pullzone.validation","Field":"",
"Message":"Nullable object must have a value."}

important data from bunny net
```sql
DROP TABLE IF EXISTS bunny;
CREATE TABLE bunny (
    -- the project id is the name of storage zone and pull zone
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    project_id VARCHAR(255) NOT NULL,

    storage_zone_id VARCHAR(255) DEFAULT '',
    storage_zone_region VARCHAR(255) DEFAULT '',
    storage_zone_password VARCHAR(255) DEFAULT '',
    pullzone_id VARCHAR(255) DEFAULT '',

    bunny_eu BOOLEAN DEFAULT false,
    bunny_us BOOLEAN DEFAULT false,
    bunny_asia BOOLEAN DEFAULT false,
    bunny_sa BOOLEAN DEFAULT false,
    bunny_af BOOLEAN DEFAULT false,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE   
);

```
     

como agrego las .env?
que la ai liste cuales son las .env nesesarias, y agregar un campo en settings
que diga enviroment variables
ahi puedo crear nuevas variables y ponerle los valores
una vez que se agregen llegan a un handler que hace
```bash
fly secrets set SECRET=***** --app $project_id
```
