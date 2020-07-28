package glob

import "../def"

var PlayerTypes = []pTypeData{
	{PType: def.PLAYER_TYPE_NEW, PName: "New"},
	{PType: def.PLAYER_TYPE_NORMAL, PName: "Normal"},
	{PType: def.PLAYER_TYPE_VETERAN, PName: "Veteran"},
	{PType: def.PLAYER_TYPE_TRUSTED, PName: "Trusted"},

	{PType: def.PLAYER_TYPE_BUILDER, PName: "Builder"},
	{PType: def.PLAYER_TYPE_MODERATOR, PName: "Moderator"},
	{PType: def.PLAYER_TYPE_ADMIN, PName: "Admin"},
	{PType: def.PLAYER_TYPE_OWNER, PName: "Owner"},
}
