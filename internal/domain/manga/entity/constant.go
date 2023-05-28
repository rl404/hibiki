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
	RelationSequel             Relation = "SEQUEL"
	RelationPrequel            Relation = "PREQUEL"
	RelationAlternativeSetting Relation = "ALTERNATIVE_SETTING"
	RelationAlternativeVersion Relation = "ALTERNATIVE_VERSION"
	RelationSideStory          Relation = "SIDE_STORY"
	RelationParentStory        Relation = "PARENT_STORY"
	RelationSummary            Relation = "SUMMARY"
	RelationFullStory          Relation = "FULL_STORY"
	RelationSpinOff            Relation = "SPIN_OFF"
	RelationAdaptation         Relation = "ADAPTATION"
	RelationCharacter          Relation = "CHARACTER"
	RelationOther              Relation = "OTHER"
)

// SearchMode is search mode.
type SearchMode string

// Available search mode.
const (
	SearchModeAll    SearchMode = "ALL"
	SearchModeSimple SearchMode = "SIMPLE"
)
