package def

const VERSION = "v0.0.12 8-2-2020-958"
const DEFAULT_PORT = ":7777"
const DEFAULT_PORT_SSL = ":7778"

//SSL
const SSL_PEM = "server.pem"
const SSL_KEY = "server.key"

//Maximums
const MAX_SECTORS = 10000
const MAX_USERS = 1000
const MAX_DESC = 950
const MAX_MLE = 100

const MAX_INPUT_LENGTH = 4096
const MAX_OUTPUT_LENGTH = 16384
const MAX_INPUT_LINES = 100
const MAX_CMATCH_SEARCH = 100

//Milliseconds
const ROUND_LENGTH_uS = 250000
const ROUND_REST_MS = 3
const CONNECT_THROTTLE_MS = 500
const WELCOME_TIMEOUT_S = 30

//Player/sector defaults
const PLAYER_START_SECTOR = 1
const PLAYER_START_ROOM = 1

const PFILE_VERSION = "0.0.1"
const SECTOR_VERSION = "0.0.1"
const HELPS_VERSION = "0.0.1"

const PASSWORD_HASH_COST = 10
const MAX_PLAYER_NAME_LENGTH = 25
const MIN_PLAYER_NAME_LENGTH = 2
const STRING_UNKNOWN = "unknown"

/*Dir & File*/
const DATA_DIR = "data/"
const PLAYER_DIR = "players/"
const SECTOR_DIR = "sectors/"
const PSECTOR_DIR = "psectors/"
const TEXTS_DIR = "texts/"

const SECTOR_PREFIX = "sec-"
const FILE_SUFFIX = ".dat"

const GREET_FILE = "greet.txt"
const AUREVOIR_FILE = "aurevoir.txt"
const NEWS_FILE = "news.txt"
const HELPS_FILE = "help.txt"

/*Server mode*/
const SERVER_RUNNING = 0
const SERVER_BOOTING = 1
const SERVER_CLOSING = 2
const SERVER_CLOSED = 3
const SERVER_PAUSED = 4
const SERVER_PRIVATE = 5

/*Connection State*/
const CON_STATE_DISCONNECTED = -3
const CON_STATE_DISCONNECTING = -2
const CON_STATE_RELOG = -1

const CON_STATE_WELCOME = 0
const CON_STATE_PASSWORD = 100

const CON_STATE_NEWS = 200
const CON_STATE_RECONNECT_CONFIRM = 300
const CON_STATE_PLAYING = 1000

//New Users
const CON_STATE_NEW_LOGIN = 400
const CON_STATE_NEW_LOGIN_CONFIRM = 500
const CON_STATE_NEW_PASSWORD = 600
const CON_STATE_NEW_PASSWORD_CONFIRM = 700

/*Player States*/
const PLAYER_UNLINKED = -1
const PLAYER_ALIVE = 0
const PLAYER_SIT = 100
const PLAYER_REST = 200
const PLAYER_SLEEP = 300
const PLAYER_STUNNED = 400
const PLAYER_DEAD = 1000

/*Errors*/
const ERROR_NONFATAL = false
const ERROR_FATAL = true

/*Player Type*/
const PLAYER_TYPE_NEW = 0
const PLAYER_TYPE_NORMAL = 100
const PLAYER_TYPE_VETERAN = 200
const PLAYER_TYPE_TRUSTED = 300

const PLAYER_TYPE_BUILDER = 700
const PLAYER_TYPE_MODERATOR = 800
const PLAYER_TYPE_ADMIN = 900
const PLAYER_TYPE_OWNER = 1000

/*OLC */
const OLC_NONE = 0
const OLC_ROOM = 100
const OLC_OBJECT = 200
const OLC_TRIGGER = 300
const OLC_MOBILE = 400
const OLC_QUEST = 500
const OLC_SECTOR = 600
const OLC_EXITS = 700

const SETTING_TYPE_BOOL = 0
const SETTING_TYPE_INT = 100
const SETTING_TYPE_STRING = 200
const SETTING_TYPE_INDEX = 300

const LINESEPA = "-------------------------------------------------------------------------------\r\n"
const LINESEPB = "_______________________________________________________________________________\r\n"

/*Mle Editor*/
const MLE_ADD = 100
const MLE_REMOVE = 200
const MLE_INSERT = 300
const MLE_REPLACE = 400

/*Objects*/
const OBJ_TYPE_NORMAL = 0
