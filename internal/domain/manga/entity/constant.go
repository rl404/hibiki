package entity

// Type is manga type.
type Type string

// Available manga type.
const (
	TypeManga      Type = "MANGA"
	TypeNovel      Type = "NOVEL"
	TypeOneShot    Type = "ONE_SHOT"
	TypeDoujinshi  Type = "DOUJINSHI"
	TypeManhwa     Type = "MANHWA"
	TypeManhua     Type = "MANHUA"
	TypeOEL        Type = "OEL"
	TypeLightNovel Type = "LIGHT_NOVEL"
)

// Status is manga publishing status.
type Status string

// Available manga publishing status.
const (
	StatusFinished     Status = "FINISHED"
	StatusReleasing    Status = "RELEASING"
	StatusNotYet       Status = "NOT_YET"
	StatusHiatus       Status = "HIATUS"
	StatusDiscontinued Status = "DISCONTINUED"
)

// Relation is manga relation type.
type Relation string

// Available manga relation.
const (
	RelationSequel             = "SEQUEL"
	RelationPrequel            = "PREQUEL"
	RelationAlternativeSetting = "ALTERNATIVE_SETTING"
	RelationAlternativeVersion = "ALTERNATIVE_VERSION"
	RelationSideStory          = "SIDE_STORY"
	RelationParentStory        = "PARENT_STORY"
	RelationSummary            = "SUMMARY"
	RelationFullStory          = "FULL_STORY"
	RelationSpinOff            = "SPIN_OFF"
	RelationAdaptation         = "ADAPTATION"
	RelationCharacter          = "CHARACTER"
	RelationOther              = "OTHER"
)

// SearchMode is search mode.
type SearchMode string

// Available search mode.
const (
	SearchModeAll    SearchMode = "all"
	SearchModeSimple SearchMode = "simple"
)
