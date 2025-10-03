![badge](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/marioweid/b1bfef0cff3b03f048d6c065fba5cbee/raw/action_badge.json)
![goRAGnarok Banner](/img/banner.png)

# goRAGnarok

The best RAG (Retrieval-Augmented Generation) application written in Go.

This project aims to stay as close to **plain Go** as possible, with minimal external dependencies.  
Currently, the **only external library** in use is the PostgreSQL driver for database access.


## Ollama
Download models:

- Gemma3:4b: `ollama pull gemma3:4b`
- all-minilm: `ollama pull all-minilm`

---

## Prepare Database

- Follow this [guide](https://ai.pydantic.dev/examples/rag/) to load the data into the PostgreSQL database provided in our Docker Compose setup.  
- Good luck with that 😉
