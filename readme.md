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
- Actually update `<world></world>` state.
- Action encoder is doing odd things with notes.
    - Maybe remove those and include past descriptions in context?
- Errors from applying actions?
    - Syntax errors just need the action encoder to try again.
    - Perhaps "not allowed" detection in code too? That'd have to start whole cycle again.
    - Keep looping until adjudicator plus action-encode give something valid.
    - Limit loop count presumably.
- Some kind of player health/condition tracking?
    - Something more interesting that can take advantage of LLM fuzzy reasoning - Joie de vivre?
- Sanitise player input, remove `<world>` tags if given.
- Normally block narrator from giving clues, but from code side add them in if they appear stuck?
- Anything to optimise LLM steps.
