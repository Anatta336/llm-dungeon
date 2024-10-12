package world

type State struct {
	MetaSetting []string            `json:"meta-setting"`
	Player      []string            `json:"player"`
	Scene       map[string][]string `json:"scene"`
	Objects     map[string][]string `json:"objects"`
}

func DungeonCell() State {
	return State{
		MetaSetting: []string{
			"Realistic historical fiction set in the 17th Century. There is no magic in this setting.",
		},
		Player: []string{
			"human woman",
		},
		Scene: map[string][]string{
			"room": {
				"dungeon cell",
				"dimly lit",
				"solid stone construction",
			},
			"north doorway": {
				"leads to a lit corridor",
				"blocked by sturdy wooden door",
			},
			"floor": {
				"covered in a thin layer of dirt and straw",
			},
		},
		Objects: map[string][]string{
			"sturdy wooden door": {
				"locked",
				"iron keyhole",
				"can be opened with the cell key",
				"leads out to a lit corridor",
				"a small gap between the door and floor allows in a little light",
			},
			"wooden bedframe": {
				"bedding made of straw and tattered cloth",
				"south side of room",
			},
			"wooden table": {
				"west side of room",
				"fragile",
			},
			"candle": {
				"not lit",
				"on the wooden table",
			},
			"loose stone": {
				"can be found on the east wall",
				"not immediately visible",
			},
			"cell key": {
				"hidden behind the loose stone",
				"only visible after moving stone",
				"can open the sturdy wooden door",
			},
			"small leather satchel": {
				"carried by player",
			},
			"flint and steel": {
				"in small leather satchel",
				"can set flammable items on fire",
			},
		},
	}
}
