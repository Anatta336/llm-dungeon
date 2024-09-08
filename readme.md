# Dungeon
An LLM-powered text adventure game that attempts to have a coherent world model.

![image](https://github.com/user-attachments/assets/3c89a7f3-89f0-4496-9f31-96fa93780302)

## Status
Very early. Core game loop isn't implemented yet.

## Quickstart
- Have Ollama installed and set up.
- `ollama serve`
- `. ./build.sh && go build -o bin/main main.go && bin/main`

## Proposed Structure (not current implementation)
```mermaid
graph TB

world[("World model")]
planChange[["LLM: Plan changes and description"]]
check[["Code: Check consistency"]]
applyChange[["Code: Apply changes"]]
user(["User"])

user --"Freeform input"--> planChange
world --> planChange
world --> check
planChange --"Planned changes"--> check
check --"Error message"-->planChange
check --"Confirmed changes"--> applyChange
applyChange --"Update state"--> world
applyChange --"Display description"--> user
```

## Plan
- Actually update `<world></world>` state.
- Action encoder is doing odd things with notes.
    - Maybe remove those and include past descriptions in context?
- Error message to adjudicate LLM step.
    - Both syntax and "not allowed" errors.
    - Keep looping until adjudicator plus action-encode give something valid.
    - Limit loop count presumably.
- Some kind of player health/condition tracking?
    - Something more interesting that can take advantage of LLM fuzzy reasoning - Joie de vivre?
- Sanitise player input, remove `<world>` tags if given.
- Normally block narrator from giving clues, but from code side add them in if they appear stuck?
- Anything to optimise LLM steps.
