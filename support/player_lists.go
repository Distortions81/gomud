package support

import (
	"../def"
	"../glob"
)

var PlayerTypes = []glob.PTypeData{
	{PType: def.PLAYER_TYPE_NEW, PName: "New"},
	{PType: def.PLAYER_TYPE_NORMAL, PName: "Normal"},
	{PType: def.PLAYER_TYPE_VETERAN, PName: "Veteran"},
	{PType: def.PLAYER_TYPE_TRUSTED, PName: "Trusted"},

	{PType: def.PLAYER_TYPE_BUILDER, PName: "Builder"},
	{PType: def.PLAYER_TYPE_MODERATOR, PName: "Moderator"},
	{PType: def.PLAYER_TYPE_ADMIN, PName: "Admin"},
	{PType: def.PLAYER_TYPE_OWNER, PName: "Owner"},
}

var WearLocationsList = []glob.WearLocations{
	{Name: "head", ID: def.OBJ_WEAR_HEAD},
	{Name: "face", ID: def.OBJ_WEAR_FACE},

	{Name: "left eye", ID: def.OBJ_WEAR_LEYE, ConflictLocationA: def.OBJ_WEAR_EYES},
	{Name: "right eye", ID: def.OBJ_WEAR_REYE, ConflictLocationA: def.OBJ_WEAR_EYES},
	{Name: "eyes", ID: def.OBJ_WEAR_EYES, ConflictLocationA: def.OBJ_WEAR_LEYE, ConflictLocationB: def.OBJ_WEAR_REYE},

	{Name: "left ear", ID: def.OBJ_WEAR_LEAR, ConflictLocationA: def.OBJ_WEAR_EARS},
	{Name: "right ear", ID: def.OBJ_WEAR_REAR, ConflictLocationA: def.OBJ_WEAR_EARS},
	{Name: "ears", ID: def.OBJ_WEAR_EARS, ConflictLocationA: def.OBJ_WEAR_LEAR, ConflictLocationB: def.OBJ_WEAR_REAR},
}
