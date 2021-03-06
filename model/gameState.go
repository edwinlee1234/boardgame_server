package model

// GameState db table
type GameState struct {
	ID         int32
	Type       string
	State      int32
	Result     string
	Seat       int32
	InsertTime string
	UpdateTime string
}

// GameState

// NotOpen owner only
const NotOpen = 0

// Opening 開放玩家
const Opening = 1

// Playing 遊戲中
const Playing = 2

// Close 關
const Close = 4

// Abort 放棄
const Abort = 5

// GameStateList GameStateList
type GameStateList []GameState

// TODO 把這些的func都塞到struct裡面
// 在model/main.go那邊把model的struct都new出來

// CreateGame 新增一局遊戲
func CreateGame(game string, seat int32, insertTime int64) (int32, error) {
	stmt, err := DB.Prepare("INSERT INTO game_state (type, seat, create_timestamp) VALUES (?, ?, ?)")
	if err != nil {
		return 0, err
	}

	val, err := stmt.Exec(game, seat, insertTime)
	if err != nil {
		return 0, err
	}

	id, err := val.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}

// FindGameByGameID 用GameID去game_state搜尋遊戲資料
func FindGameByGameID(id int32) (gameType string, state int32, seat int32, time int32) {
	row := DB.QueryRow("SELECT `type`, `state`, `seat`, `create_timestamp` game_state FROM `game_state` WHERE `id` = ? LIMIT 1", id)
	row.Scan(&gameType, &state, &seat, &time)

	return gameType, state, seat, time
}

// ChangeGameStateDB 改變state
func ChangeGameStateDB(id int32, state int32) error {
	stmt, _ := DB.Prepare("UPDATE `game_state` set `state` = ? where `id` = ?")
	res, _ := stmt.Exec(state, id)

	affect, err := res.RowsAffected()
	if err != nil || affect == 0 {
		return err
	}

	return nil
}

// FindAllGameByState 找出全部遊戲By State
func FindAllGameByState(state int32) (GameStateList, error) {
	var res GameStateList

	rows, err := DB.Query(
		"SELECT `id`, `type`, `state`, `seat`, `create_timestamp`, `update_time` FROM `game_state` WHERE `state` = ?", state)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var row GameState

		err = rows.Scan(&row.ID, &row.Type, &row.State, &row.Seat, &row.InsertTime, &row.UpdateTime)
		if err != nil {
			return res, err
		}

		res = append(res, row)
	}

	return res, nil
}
