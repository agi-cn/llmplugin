# LLM Plugin

LLM Plugin system.

## 1. Plugins

### 2.1 Google Search

Get Google Search token: [https://docs.chatkit.app/tools/google-search.html](https://docs.chatkit.app/tools/google-search.html)

## 2. TESTING

1. OpenAI:
   ```bash
   cp .env.example .env
   ```

2. Google:
   ```bash
   cd plugins/google

   cp .env.example .env
   ```


Run test:

```bash
go test -v ./...
```


## 3. RELEASE

### v0.1.0

1. init project.
2. support plugin: Google for search, calculator for mathematical calculations.
