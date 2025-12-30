FROM nikolaik/python-nodejs:python3.10-nodejs22

RUN npm install -g @anthropic-ai/claude-code

WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .

EXPOSE 5173
EXPOSE 5000

CMD ["./start.sh"]
