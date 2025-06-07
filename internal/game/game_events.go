package game

type GameEvent int
type GameWebSocketMessage struct {
	Event   GameEvent   `json:"event"`
	Payload interface{} `json:"payload"`
}

const (
	Authentication  GameEvent = 0
	JoinGame        GameEvent = 1
	LeaveGame       GameEvent = 2
	GameStateUpdate GameEvent = 3
	PlayerJoin      GameEvent = 4
	PlayerLeave     GameEvent = 5
	PlayerMove      GameEvent = 6
	PlayerDash      GameEvent = 7
	Error           GameEvent = 252
	Forbidden       GameEvent = 253
	Unauthorized    GameEvent = 254
	ServerError     GameEvent = 255
)

// String returns the string representation of GameEvent
func (g GameEvent) String() string {
	switch g {
	case Authentication:
		return "Authentication"
	case JoinGame:
		return "JoinGame"
	case LeaveGame:
		return "LeaveGame"
	case GameStateUpdate:
		return "GameStateUpdate"
	case PlayerJoin:
		return "PlayerJoin"
	case PlayerLeave:
		return "PlayerLeave"
	case PlayerMove:
		return "PlayerMove"
	case PlayerDash:
		return "PlayerDash"
	case Error:
		return "Error"
	case Forbidden:
		return "Forbidden"
	case Unauthorized:
		return "Unauthorized"
	case ServerError:
		return "ServerError"
	default:
		return "Unknown"
	}
}
