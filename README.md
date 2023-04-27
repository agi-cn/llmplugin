# LLM Plugin

LLM Plugin system.


## TESTING

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


# RELEASE

## v0.1.0

1. init project.
2. support plugin: Google for search, calculator for mathematical calculations.