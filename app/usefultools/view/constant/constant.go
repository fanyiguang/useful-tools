package constant

const (
	NavStatePreferenceCurrentPage  = "currentPage"
	NavStatePreferenceProMode      = "proMode"
	NavStatePreferenceHideBody     = "hideBody"
	NavStatePreferenceSaveAesKey   = "saveAesKey"
	NavStatePreferenceCloseUpgrade = "closeUpgrade"
)

type ViewMode int

const (
	ViewModeNormal ViewMode = iota
	ViewModePro
)
