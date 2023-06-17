package burningtext

import "image/color"

type BurningTextOptions struct {
	Text      string
	Speed     BurningSpeed
	FontChain string

	FlameColor  color.Color
	TextColor   color.Color
	RandomColor bool
}

func (settings BurningTextOptions) GetFontChain() []string {
	if settings.FontChain == "" {
		return getDefaultFontChain()
	}

	if fontChains == nil {
		loadFontChainsData()
	}

	return fontChains[settings.FontChain]
}
