package constant

const (
	NavStatePreferenceCurrentPage  = "currentPage"
	NavStatePreferenceProMode      = "proMode"
	NavStatePreferenceHideBody     = "hideBody"
	NavStatePreferenceSaveAesKey   = "saveAesKey"
	NavStatePreferenceCloseUpgrade = "closeUpgrade"
	NavStatePreferenceStyle        = "style"

	CacheKeyAesKeyList = "aes-key-list"
)

const (
	StyleDefault            = "default"
	StyleLowSaturationGreen = "low_saturation_green"
	StyleWarmLuxury         = "warm_luxury"
	StyleNeutralMinimal     = "neutral_minimal"
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
