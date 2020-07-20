package def

const DEFAULT_PORT = ":7777"
const MAX_DESCRIPTORS = 10000
const STRING_UNKNOWN = "unknown"

const SERVER_RUNNING = 0
const SERVER_PAUSED = 1
const SERVER_STOP = 2

/*Connection State*/
const CON_STATE_DISCONNECTED = -2
const CON_STATE_DISCONNECTING = -1

const CON_STATE_WELCOME = 0
const CON_STATE_ENTER_LOGIN = 10
const CON_STATE_INVALID_LOGIN = 20

const CON_STATE_PASSWORD = 30
const CON_STATE_INVALID_PASSWORD = 40

const CON_STATE_NEWS = 50
const CON_STATE_PLAYING = 100

/*New Users*/
const CON_STATE_NEW_LOGIN = 1
const CON_STATE_NEW_LOGIN_CONFIRM = 2
const CON_STATE_NEW_PASSWORD = 3
const CON_STATE_NEW_PASSWORD_CONFIRM = 4

/*Player States*/
const PLAYER_ALIVE = 0
const PLAYER_SIT = 10
const PLAYER_REST = 20
const PLAYER_SLEEP = 30
const PLAYER_STUNNED = 40
const PLAYER_DEAD = 50

/*Errors*/
const ERROR_NONFATAL = false
const ERROR_FATAL = true

/*Player Type*/
const PLAYER_TYPE_NEW = 0
const PLAYER_TYPE_NORMAL = 10
const PLAYER_TYPE_VETERAN = 20
const PLAYER_TYPE_TRUSTED = 30

const PLAYER_TYPE_BUILDER = 70
const PLAYER_TYPE_MODERATOR = 80
const PLAYER_TYPE_ADMIN = 90
const PLAYER_TYPE_OWNER = 100
