# Dungeon
An LLM-powered text adventure game that attempts to have a coherent world model.

![image](https://github.com/user-attachments/assets/3c89a7f3-89f0-4496-9f31-96fa93780302)

## Status
Very early. Core game loop isn't implemented yet.

## Quickstart
- Have Ollama installed and set up.
- `ollama serve`
- `. ./build.sh && go build -o bin/main main.go && bin/main`

## Process structure
```mermaid
graph TB

user(["User"])
llmAdjudicate[["LLM: Adjudicate"]]
world[("World model")]
llmEncode[["LLM: Encode Actions"]]
codeActions[["Code: Apply Actions (NYI)"]]

user --"Freeform input"--> llmAdjudicate
world --> llmAdjudicate
llmAdjudicate --"Freeform description"--> llmEncode
world --> llmEncode
llmEncode --"Encoded actions"--> codeActions
world --> codeActions
codeActions --"Update"--> world
llmAdjudicate --"Freeform description"--> user
```

## Plan
- Do not rely on LLM to update JSON manually. Instead have it provide a diff.
- Convince the adjudicator to stop inventing objects.
- Some kind of player health/condition tracking?
    - Something more interesting that can take advantage of LLM reasoning.
- Anything to optimise LLM steps.
