package dm

import (
	"encoding/json"
	"fmt"
	"samdriver/dungeon/llm"
	"samdriver/dungeon/world"
	"strings"
)

type DmResponseMessage struct {
	UserInput          string `json:"user-input"`
	InputState         world.State
	RawAdjudicate      string `json:"raw-adjudicate"`
	AdjudicateThoughts string `json:"adjudicate-thoughts"`
	Description        string `json:"description"`
	RawEncode          string `json:"raw-action-encode"`
	OutputState        world.State
}

func Process(state world.State, input string) (DmResponseMessage, error) {
	resp := DmResponseMessage{
		UserInput:  input,
		InputState: state,
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

	stateJson, err := json.Marshal(dmResponse.InputState)
	if err != nil {
		return fmt.Errorf("error marshaling InputState to JSON: %v", err)
	}

	systemPrompt := adjudicateSystem + "\nCurrent world JSON:\n" + string(stateJson)

	// TODO: include recent history of messages in context?

	request := llm.Request{
		System:      systemPrompt,
		User:        dmResponse.UserInput,
		Model:       "llama3.1:latest",
		Temperature: 0.7,
	}

	response, err := request.Process()

	if err != nil {
		return fmt.Errorf("error with Adjudicate LLM: %v", err)
	}

	thinkingStart := strings.Index(response.Result, "THINKING:")
	outputStart := strings.Index(response.Result, "OUTPUT:")
	thoughts := func() string {
		if thinkingStart == -1 || outputStart == -1 {
			return ""
		}
		return response.Result[thinkingStart+len("THINKING:") : outputStart]
	}()

	description := func() string {
		if outputStart == -1 {
			// Assume everything is the description we want.
			return response.Result
		}
		return response.Result[outputStart+len("OUTPUT:") : len(response.Result)]
	}()

	dmResponse.RawAdjudicate = response.Result
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.Description = description

	return nil
}

func (dmResponse *DmResponseMessage) fillEncodedActions() error {
	stateJson, err := json.Marshal(dmResponse.InputState)
	if err != nil {
		return fmt.Errorf("error marshaling InputState to JSON: %v", err)
	}

	systemPrompt := actionEncodeSystem + "\nCurrent world JSON:\n" + string(stateJson)

	request := llm.Request{
		System:      systemPrompt,
		User:        dmResponse.Description,
		Model:       "llama3.1:latest",
		Temperature: 0.2,
	}

	response, err := request.Process()

	if err != nil {
		return fmt.Errorf("error with action encode LLM: %v", err)
	}

	thinkingStart := strings.Index(response.Result, "THINKING:")
	outputStart := strings.Index(response.Result, "OUTPUT:")
	thoughts := func() string {
		if thinkingStart == -1 || outputStart == -1 {
			return ""
		}
		return response.Result[thinkingStart+len("THINKING:") : outputStart]
	}()

	jsonString := func() string {
		if outputStart == -1 {
			// Invalid, so make no changes.
			return ""
		}
		return response.Result[outputStart+len("OUTPUT:") : len(response.Result)]
	}()

	var newState world.State
	err = json.Unmarshal([]byte(jsonString), &newState)
	if err != nil {
		// Invalid JSON, so make no changes.
		newState = dmResponse.InputState
	}

	dmResponse.RawEncode = response.Result
	dmResponse.InputState = newState
	dmResponse.AdjudicateThoughts = thoughts
	dmResponse.OutputState = newState

	return nil
}

const adjudicateSystem = `
You are playing the role of a dungeon master in a text-based role-playing game. Your job is to adjudicate the actions of the player and describe the world around them based on the provided world state.

**Instructions**:

- **THINKING**:
  - Begin each response with "THINKING:" followed by your internal reasoning.
  - Consider what the player is trying to achieve and whether it is possible given the world state information and the abilities of the player's character.
  - **Do not** add actions beyond what the player is specifically asking for.
  - Plan the consequences of the player's actions, including any changes to the world.
  - This section is **not** shown to the player.

- **OUTPUT**:
  - After the "THINKING:" section, write "OUTPUT:" followed by the description presented to the player.
  - Stay in character as the dungeon master.
  - Refer to the player as "you".
  - Describe the results of the player's actions and the state of the world.
  - Do not suggest actions; allow the player to decide their next move.
  - Provide realistic consequences based on the world state.
  - Do not reject actions due to danger; allow the player to take risks.
  - If an action is impossible, gently explain why.

- Only reference objects, characters, and locations included in the world state.
- Only provide details about objects the player is specifically asking about.
- Only describe actions the player specifically requests. Do not provide extra information or actions.
- Be sure to mention anything in the scene that presents an immediate danger to the player.
- As well as the player's direct actions, also consider what any active objects in the scene might be doing.
`

const actionEncodeSystem = `
Your task is to update the JSON representation of the world state based on new descriptions provided. Read the input description and modify the world state accordingly. Ensure that any changes, additions, or removals are accurately reflected in the JSON.

**Instructions**:
- **THINKING**:
  - Begin each response with "THINKING:" followed by your internal reasoning.
  - Consider the description provided and how it should be reflected in changes to the world state.

- **OUTPUT**:
  - After the "THINKING:" section, write "OUTPUT:" followed by the updated JSON representation of the world state.
  - Ensure that the JSON structure is maintained.
  - If the description mentions new items, add them to the appropriate section of the JSON.
  - If items are removed, remove them from the JSON.
  - Remove any descriptions that are no longer accurate.
  - The JSON should be the only text after "OUTPUT:".

- Only change the world state if new information is provided or changes to the scene, objects, or characters are described.
`
