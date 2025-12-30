ARG GO_VERSION=1.22
FROM golang:${GO_VERSION}-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o x73-app cmd/main.go

FROM debian:bookworm
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates curl gnupg git && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL https://deb.nodesource.com/setup_20.x | bash - && \
    apt-get install -y --no-install-recommends nodejs && \
    rm -rf /var/lib/apt/lists/*

RUN mkdir -p -m 755 /etc/apt/keyrings && \
    curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg \
    > /etc/apt/keyrings/githubcli-archive-keyring.gpg && \
    chmod go+r /etc/apt/keyrings/githubcli-archive-keyring.gpg && \
    echo "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/githubcli-archive-keyring.gpg] \
    https://cli.github.com/packages stable main" > /etc/apt/sources.list.d/github-cli.list && \
    apt-get update && \
    apt-get install -y --no-install-recommends gh && \
    rm -rf /var/lib/apt/lists/*

ENV FLYCTL_INSTALL=/root/.fly
ENV PATH=/root/.fly/bin:$PATH

RUN curl -L https://fly.io/install.sh | sh

# RUN mkdir /data/projects
# RUN mkdir /data/fly_configs

RUN ln -s /data/projects /app/projects
RUN ln -s /data/fly_configs /app/fly_configs

RUN npm install -g wrangler@latest
RUN git config --global credential.helper '!gh auth git-credential'

RUN git config --global user.email "agustfricke@gmail.com"
RUN git config --global user.name "x73"

COPY --from=builder /app/x73-app /app/x73-app
COPY --from=builder /app/templates /app/templates
COPY --from=builder /app/claude /app/claude
COPY --from=builder /app/scripts /app/scripts

EXPOSE 8080

CMD ["./x73-app"]
