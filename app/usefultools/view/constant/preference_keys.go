package constant

const (
	NavStatePreferenceCurrentPage  = "currentPage"
	NavStatePreferenceProMode      = "proMode"
	NavStatePreferenceHideBody     = "hideBody"
	NavStatePreferenceSaveAesKey   = "saveAesKey"
	NavStatePreferenceCloseUpgrade = "closeUpgrade"

	CacheKeyAesKeyList = "aes-key-list"
)

func CacheKeys() []string {
	return []string{
		CacheKeyAesKeyList,
	}
}

type ViewMode int

const (
	ViewModeNormal ViewMode = iota
	ViewModePro
)
