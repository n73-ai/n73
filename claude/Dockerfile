FROM python:3.10-slim

RUN apt-get update && apt-get install -y \
    wget \
    xz-utils \
    git \
    && rm -rf /var/lib/apt/lists/*

# Variables para Node.js
ARG NODE_VERSION=22.17.1

RUN wget https://nodejs.org/dist/v${NODE_VERSION}/node-v${NODE_VERSION}-linux-x64.tar.xz \
    && tar -xJf node-v${NODE_VERSION}-linux-x64.tar.xz \
    && mv node-v${NODE_VERSION}-linux-x64 /usr/local/nodejs \
    && rm node-v${NODE_VERSION}-linux-x64.tar.xz

ENV PATH="/usr/local/nodejs/bin:$PATH"

RUN npm install -g @anthropic-ai/claude-code

RUN node --version && npm --version && python --version

WORKDIR /app

COPY requirements.txt ./

COPY . .

RUN pip install -r requirements.txt 

EXPOSE 5000

CMD ["./start.sh"]
