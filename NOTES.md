# notes

fix this errors and add this new features:
- cuando sale el error try to fix, donde va el website iframe, sigue cargando
esto esta mal, deberia decir algo como, el primer deployment fallo
click try to fix to ask n83 to fix the error
- i need to do an npm install on ts-claude/ui-only to do a npm run build
it need to be fast so do the install while the project is being generated

  1. spawn the Vite dev server (npm run dev -- --host 0.0.0.0 --port 5173) in the background
  2. Poll http://0.0.0.0:5173 every 500ms (up to 30s) until it responds
  3. Take the screenshot
  4. devServer.kill() in finally — always killed regardless of success/failure
  5. npm run build → copy screenshot into dist → zip → upload to bunny

#### ==

# take screenshots

then edit readme.md and replace the word food for cat.",
c55ef272-d033-4814-8860-9b1afe9ecfdd

wget https://dl.google.com/linux/direct/google-chrome-stable_current_amd64.deb

- [ ] create app (fly/create.go)
- [ ] deploy manual (fly/create.go)

crear apps en fly manualmente para ver si funciona con: 
0: old-current
```bash
FROM nikolaik/python-nodejs:python3.10-nodejs22

RUN npm install -g @anthropic-ai/claude-code

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

EXPOSE 5173
EXPOSE 5000

CMD ["./start.sh"]

---

#!/bin/bash

mkdir /root/.claude
cp .credentials.json /root/.claude/.credentials.json

python main.py &

npm --prefix /app/ui-only install
npm run --prefix /app/ui-only build
python3 -m http.server 5173 --directory /app/ui-only/dist
```

1:
```Dockerfile
FROM debian:bullseye-slim

# Instalar dependencias mínimas
RUN apt-get update && apt-get install -y \
    wget \
    gnupg \
    unzip \
    curl \
    ca-certificates \
    --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Instalar Chrome
RUN wget -q -O - https://dl-ssl.google.com/linux/linux_signing_key.pub | apt-key add - && \
    echo "deb http://dl.google.com/linux/chrome/deb/ stable main" >> /etc/apt/sources.list.d/google.list && \
    apt-get update && \
    apt-get install -y google-chrome-stable --no-install-recommends && \
    rm -rf /var/lib/apt/lists/*

# Instalar ChromeDriver
RUN CHROMEDRIVER_VERSION=$(curl -sS chromedriver.storage.googleapis.com/LATEST_RELEASE) && \
    wget -q "https://chromedriver.storage.googleapis.com/${CHROMEDRIVER_VERSION}/chromedriver_linux64.zip" && \
    unzip chromedriver_linux64.zip && \
    mv chromedriver /usr/local/bin/ && \
    chmod +x /usr/local/bin/chromedriver && \
    rm chromedriver_linux64.zip

WORKDIR /app

# Si usas Python/Selenium
RUN apt-get update && apt-get install -y python3 python3-pip && \
pip3 install selenium

CMD ["/bin/bash"]
```

2: 
```Dockerfile
FROM zenika/alpine-chrome:with-chromedriver

WORKDIR /app
COPY . .

# Ya trae Chrome + ChromeDriver
```


---

checkout typescript: 
    "https://platform.claude.com/docs/en/agent-sdk/quickstart#typescript"

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
