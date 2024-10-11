package dm

import (
	"fmt"
	"samdriver/dungeon/llm"
	"strings"
)

type DmResponseMessage struct {
	UserInput          string `json:"user-input"`
	RawAdjudicate      string `json:"raw-adjudicate"`
	AdjudicateThoughts string `json:"adjudicate-thoughts"`
	Description        string `json:"description"`
	RawActionEncode    string `json:"raw-action-encode"`
	Actions            string `json:"actions"`
}

func Process(input string) (DmResponseMessage, error) {
	resp := DmResponseMessage{
		UserInput: input,
	}
	err := resp.fillDescription()
	if err != nil {
		return resp, err
	}

	err = resp.fillEncodedActions()
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func (dmResponse *DmResponseMessage) fillDescription() error {

	request := llm.Request{
		System:      adjudicateSystem,
		User:        dmResponse.UserInput,
		Model:       "llama3.2:3b",
		Temperature: 0.7,
	}

	response, err := request.Process()

	if err != nil {
		return fmt.Errorf("error with Adjudicate LLM: %v", err)
	}

	thinkingStart := strings.Index(response.Result, "<thinking>")
	thinkingEnd := strings.Index(response.Result, "</thinking>")
	thoughts := func() string {
		if thinkingStart == -1 || thinkingEnd == -1 {
			return ""
		}
		return response.Result[thinkingStart+len("<thinking>") : thinkingEnd]
	}()

	outputStart := strings.Index(response.Result, "<output>")
	outputEnd := strings.Index(response.Result, "</output>")
	description := func() string {
		if outputStart == -1 || outputEnd == -1 {
			// Assume everything after the </thinking> is the description we want.
			return response.Result[thinkingEnd+len("</thinking>"):]
		}
		return response.Result[outputStart+len("<output>") : outputEnd]
	}()

	dmResponse.RawAdjudicate = response.Result
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.Description = description

	return nil
}

func (dmResponse *DmResponseMessage) fillEncodedActions() error {
	request := llm.Request{
		System:      actionEncodeSystem,
		User:        dmResponse.Description,
		Model:       "llama3.2:3b",
		Temperature: 0.2,
	}

	response, err := request.Process()

	if err != nil {
		return fmt.Errorf("error with action encode LLM: %v", err)
	}

	thinkingStart := strings.Index(response.Result, "<thinking>")
	thinkingEnd := strings.Index(response.Result, "</thinking>")
	thoughts := func() string {
		if thinkingStart == -1 || thinkingEnd == -1 {
			return ""
		}
		return response.Result[thinkingStart+len("<thinking>") : thinkingEnd]
	}()

	outputStart := strings.Index(response.Result, "<output>")
	outputEnd := strings.Index(response.Result, "</output>")
	actions := func() string {
		if outputStart == -1 || outputEnd == -1 {
			return ""
		}
		return response.Result[outputStart+len("<output>") : outputEnd]
	}()

	dmResponse.RawActionEncode = response.Result
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.Actions = actions

	return nil
}

const adjudicateSystem = `
You are a world-class AI system playing the role of a dungeon master in a text adventure game, capable of complex reasoning and reflection.

Always start your response with <thinking></thinking> tags where you reason through the user's request. This reasoning should include if what they're asking for is possible.
Check over your reasoning and correct any mistakes.
After the </thinking> section completes, provide your final response inside <output> tags.
The user will only see the <output></output> section so make sure it includes everything they need to see.

You may make some lighthearted remarks to keep the game fun. Do not make any sexually suggestive content.
You are not a general purpose assistant.

You help the user by telling them about the game world that is described within <world></world> tags.
If the <world></world> state indicate that an item cannot be seen (for example if it is inside a closed container) do not mention it to the player in your output.
Stay in character and avoid mentioning "the description", instead just say that something isn't possible or can't be found.
Allow the user to make choices about what their character will do. Do not suggest actions.
Refer to the user as "you".
Give realistic consequences to attempted actions, obeying information given in the <world></world> tags.
If a user attempts to do something impossible given the situation that has been described, DO NOT agree with them. Instead explain why it is not possible.
You must only reference objects, characters, and locations that are included within the <world></world> tags.

<world>
{
    "meta-setting": "Realistic historical fiction set in the 17th Century. There is no magic in this setting."
    "environment": "A damp dungeon cell. Walls, floor, and roof all made of rough cut stone.",
    "player": {
        "wearing": [
            "skirt",
            "shirt",
        ],
        "carrying": [
            "small leather satchel": {
                "contains": [
                    "flint and steel",
                ],
            },
        ],
        "description": [
            "tall woman",
        ],
    },
    "scene": [
        "north wall": {
            "in it": [
                "sturdy wooden door": {
                    "location": "middle",
                },
            ],
        },
        "floor": {
            "on it": [
                "wooden bedframe": {
                    "on it": [
                        "straw mattress": {
                            "in it": [
                                "loose straw": {
                                    "in it": [
                                        "cell key",
                                    ],
                                },
                            ],
                        },
                    ],
                    "location": {
                        "south east corner",
                    }
                },
                "wooden table": {
                    "on it": [
                        "candle",
                    ],
                },
                "pool of foul water": {
                    "location": {
                        "south west corner",
                    },
                },
            ],
        },
    ],
    "item details": [
        "flint and steel": {
            "interactions": [
                "can set flammable items on fire"
            ],
        },
        "sturdy wooden door": {
            "state": [
                "locked",
            ],
            "interactions: [
                "can be opened with the cell key",
            ],
            "notes": [
                "a very small gap between the door and floor allows in a little light",
            ]
            "location": "middle",
        },
        "candle": {
            "state": [
                "unlit",
            ],
        },
    ],
}
</world>
`

const actionEncodeSystem = `
You are a world-class AI system capable of complex reasoning and reflection. Your purpose is to turn provided text into a set of specially encoded statements.

Always start your response with <thinking></thinking> tags where you reason through how you will encode the request.
Check over your reasoning to correct any mistakes and note if any statements aren't actually needed.
After the </thinking> section completes, provide your final response inside <output> tags.
Inside the <output></output> tags should only be specially encoded actions.
When actions correspond with items in the world, you should use the same name they have in the <world></world> tags.

The only valid encode statements are:

- 'move {item} {location}' when an item changes its location in the scene.
- 'set state {item} {state}' when some ongoing state about an object changes.
- 'unset state {item} {state}' when state or information no longer applies to an object.
- 'add note {item} {note}' when the provided text reveals new information about the object that's not already implicitly provided in the <world></world> description.
- 'add {item} {location}' when a completely new item is introduced into the scene. Only use this if there is definitely a new item.
- 'remove {item}' when an item is completely destroyed or removed from the scene.

Some examples of properly encoded actions:
- The player picks up a stone off the ground: 'move "stone" "player"."carrying"."stone"'
- Stone leaves the scene completely (only use this if the item is completely gone, not just moved somewhere else): 'remove "stone"'
- Stone is moved from the floor to a desk: 'move "stone" "scene"."floor"."on it"."desk"."on it"'
- Shirt the player is wearing is splashed with water: 'set state "shirt" "wet"'
- Stone becomes hot: 'set state "stone" "hot"'
- Candle that was lit is extinguished: 'unset state "candle" "lit" set state "candle" "unlit"'
- Chest is described as being inlaid with gold decorations: 'add note "chest" "inlaid with gold decorations"'
- Chest is unlocked and opened: 'clear state "chest" "locked" set state "chest" "open"'

<world>
{
    "meta-setting": "Realistic historical fiction set in the 17th Century. There is no magic in this setting."
    "environment": "A damp dungeon cell. Walls, floor, and roof all made of rough cut stone.",
    "player": {
        "wearing": [
            "skirt",
            "shirt",
        ],
        "carrying": [
            "small leather satchel": {
                "contains": [
                    "flint and steel",
                ],
            },
        ],
        "description": [
            "tall woman",
        ],
    },
    "scene": [
        "north wall": {
            "in it": [
                "sturdy wooden door": {
                    "location": "middle",
                },
            ],
        },
        "floor": {
            "on it": [
                "wooden bedframe": {
                    "on it": [
                        "straw mattress": {
                            "in it": [
                                "loose straw": {
                                    "in it": [
                                        "cell key",
                                    ],
                                },
                            ],
                        },
                    ],
                    "location": {
                        "south east corner",
                    }
                },
                "wooden table": {
                    "on it": [
                        "candle",
                    ],
                },
                "pool of foul water": {
                    "location": {
                        "south west corner",
                    },
                },
            ],
        },
    ],
    "item details": [
        "flint and steel": {
            "interactions": [
                "can set flammable items on fire"
            ],
        },
        "sturdy wooden door": {
            "state": [
                "locked",
            ],
            "interactions: [
                "can be opened with the cell key",
            ],
            "notes": [
                "a very small gap between the door and floor allows in a little light",
            ]
            "location": "middle",
        },
        "candle": {
            "state": [
                "unlit",
            ],
        },
    ],
}
</world>
`
