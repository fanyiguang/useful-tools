package style

import (
	"image/color"
	"useful-tools/app/usefultools/view/constant"

	"fyne.io/fyne/v2"
	fyneTheme "fyne.io/fyne/v2/theme"
)

func Apply(a fyne.App, style string) {
	if a == nil {
		return
	}
	switch style {
	case constant.StyleLowSaturationGreen:
		a.Settings().SetTheme(lowSaturationGreenTheme{base: fyneTheme.DefaultTheme()})
	case constant.StyleWarmLuxury:
		a.Settings().SetTheme(warmLuxuryTheme{base: fyneTheme.DefaultTheme()})
	case constant.StyleNeutralMinimal:
		a.Settings().SetTheme(neutralMinimalTheme{base: fyneTheme.DefaultTheme()})
	default:
		a.Settings().SetTheme(fyneTheme.DefaultTheme())
	}
}

type lowSaturationGreenTheme struct {
	base fyne.Theme
}

func (t lowSaturationGreenTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case fyneTheme.ColorNamePrimary:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 121, G: 151, B: 133, A: 255}
		}
		return color.NRGBA{R: 92, G: 122, B: 106, A: 255}
	case fyneTheme.ColorNameButton:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 103, G: 130, B: 114, A: 255}
		}
		return color.NRGBA{R: 180, G: 196, B: 186, A: 255}
	case fyneTheme.ColorNameFocus:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 135, G: 165, B: 146, A: 255}
		}
		return color.NRGBA{R: 138, G: 168, B: 150, A: 255}
	case fyneTheme.ColorNameHover:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 120, G: 146, B: 130, A: 255}
		}
		return color.NRGBA{R: 170, G: 190, B: 178, A: 255}
	case fyneTheme.ColorNameSelection:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 94, G: 118, B: 106, A: 255}
		}
		return color.NRGBA{R: 196, G: 210, B: 200, A: 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t lowSaturationGreenTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t lowSaturationGreenTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t lowSaturationGreenTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

type warmLuxuryTheme struct {
	base fyne.Theme
}

func (t warmLuxuryTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case fyneTheme.ColorNamePrimary:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 176, G: 148, B: 104, A: 255}
		}
		return color.NRGBA{R: 166, G: 140, B: 96, A: 255}
	case fyneTheme.ColorNameButton:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 150, G: 132, B: 96, A: 255}
		}
		return color.NRGBA{R: 220, G: 210, B: 188, A: 255}
	case fyneTheme.ColorNameFocus:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 193, G: 170, B: 126, A: 255}
		}
		return color.NRGBA{R: 196, G: 176, B: 132, A: 255}
	case fyneTheme.ColorNameHover:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 156, G: 136, B: 100, A: 255}
		}
		return color.NRGBA{R: 232, G: 222, B: 202, A: 255}
	case fyneTheme.ColorNameSelection:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 132, G: 114, B: 84, A: 255}
		}
		return color.NRGBA{R: 238, G: 230, B: 212, A: 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t warmLuxuryTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t warmLuxuryTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t warmLuxuryTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}

type neutralMinimalTheme struct {
	base fyne.Theme
}

func (t neutralMinimalTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case fyneTheme.ColorNamePrimary:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 150, G: 150, B: 150, A: 255}
		}
		return color.NRGBA{R: 96, G: 96, B: 96, A: 255}
	case fyneTheme.ColorNameButton:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 72, G: 72, B: 72, A: 255}
		}
		return color.NRGBA{R: 216, G: 216, B: 216, A: 255}
	case fyneTheme.ColorNameFocus:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 110, G: 110, B: 110, A: 255}
		}
		return color.NRGBA{R: 140, G: 140, B: 140, A: 255}
	case fyneTheme.ColorNameHover:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 88, G: 88, B: 88, A: 255}
		}
		return color.NRGBA{R: 232, G: 232, B: 232, A: 255}
	case fyneTheme.ColorNameSelection:
		if variant == fyneTheme.VariantDark {
			return color.NRGBA{R: 64, G: 64, B: 64, A: 255}
		}
		return color.NRGBA{R: 200, G: 200, B: 200, A: 255}
	default:
		return t.base.Color(name, variant)
	}
}

func (t neutralMinimalTheme) Font(style fyne.TextStyle) fyne.Resource {
	return t.base.Font(style)
}

func (t neutralMinimalTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return t.base.Icon(name)
}

func (t neutralMinimalTheme) Size(name fyne.ThemeSizeName) float32 {
	return t.base.Size(name)
}
