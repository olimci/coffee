package coffee

type BannerHandle struct {
	coffee *Coffee
	banner *Banner
}

func (c *Coffee) Banner(text string, opts ...SubmodelOption) (*BannerHandle, error) {
	banner := NewBanner(text)

	submodelOpts := append([]SubmodelOption{WithSection(SectionHeader), WithFocusBehind()}, opts...)
	if err := c.AddSubmodel(banner, submodelOpts...); err != nil {
		return nil, err
	}

	return &BannerHandle{
		coffee: c,
		banner: banner,
	}, nil
}

func (h *BannerHandle) Set(text string) error {
	return h.coffee.send(msgBannerSet{banner: h.banner, text: text})
}

func (h *BannerHandle) Clear() error {
	return h.coffee.send(msgBannerClear{banner: h.banner})
}
