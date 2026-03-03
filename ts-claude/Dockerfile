FROM node:20-alpine3.19

# Instalar Python 3.11 y Chromium
RUN apk add --no-cache \
    python3 \
    chromium \
    chromium-chromedriver \
    nss \
    freetype \
    harfbuzz \
    ca-certificates \
    ttf-freefont

ENV CHROME_BIN=/usr/bin/chromium-browser \
    CHROME_PATH=/usr/lib/chromium/

WORKDIR /app
COPY . .

# Build en imagen, no en runtime
# RUN npm --prefix /app install
# RUN npm --prefix /app/ui-only install
# RUN npm run --prefix /app/ui-only build

EXPOSE 5000
CMD ["./init.sh"]
