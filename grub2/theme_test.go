package main

import (
	"encoding/json"
	. "launchpad.net/gocheck"
)

var theme *Theme

func init() {
	theme = &Theme{}
	Suite(theme)
}

func (theme *Theme) TestGetplJsonData(c *C) {
	testJsonContent := `{
    "DefaultTplValue":{"Background":"background_default.jpg","ItemColor":"#a6a6a6","SelectedItemColor":"#fefefe"},
    "LastTplValue":{"Background":"background.jpg","ItemColor":"#a6a6a0","SelectedItemColor":"#fefef0"}
}`
	wantJsonData := &TplJsonData{}
	wantJsonData.DefaultTplValue = TplValues{"background_default.jpg", "#a6a6a6", "#fefefe"}
	wantJsonData.LastTplValue = TplValues{"background.jpg", "#a6a6a0", "#fefef0"}

	jsonData, err := theme.getTplJsonData([]byte(testJsonContent))
	if err != nil {
		c.Error(err)
	}
	c.Check(*jsonData, Equals, *wantJsonData)
}

func (theme *Theme) TestGetNewBgFileName(c *C) {
	tests := []struct {
		s, want string
	}{
		{"/a/b/c/d/image.png", "background.png"},
		{"/image2.jpg", "background.jpg"},
	}
	for _, t := range tests {
		c.Check(theme.getNewBgFileName(t.s), Equals, t.want)
	}
}

func (theme *Theme) TestGetCustomizedThemeContent(c *C) {
	testThemeTplContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "{{.Background}}"
desktop-color: "#000000"
terminal-box: "terminal_*.png"
terminal-font: "Fixed Regular 13"

# Show the boot menu
+ boot_menu {
  left = 15%
  top = 20%
  width = 70%
  height = 60%
  item_font = "Courier 10 Pitch Bold 16"
  selected_item_font = "Courier 10 Pitch Bold 24"
  item_color = "{{.ItemColor}}"
  selected_item_color = "{{.SelectedItemColor}}"
  item_spacing = 0
  menu_pixmap_style = "empty_*.png"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "sb_th_*.png"
}
`
	testThemeTplJSON := `{"Background": "background.jpg","ItemColor":"#a6a6a6","SelectedItemColor":"#fefefe"}`
	wantThemeTxtContent := `# GRUB2 gfxmenu Linux Deepin theme
# Designed for 1024x768 resolution
# Global Property
title-text: ""
desktop-image: "background.jpg"
desktop-color: "#000000"
terminal-box: "terminal_*.png"
terminal-font: "Fixed Regular 13"

# Show the boot menu
+ boot_menu {
  left = 15%
  top = 20%
  width = 70%
  height = 60%
  item_font = "Courier 10 Pitch Bold 16"
  selected_item_font = "Courier 10 Pitch Bold 24"
  item_color = "#a6a6a6"
  selected_item_color = "#fefefe"
  item_spacing = 0
  menu_pixmap_style = "empty_*.png"
  scrollbar = true
  scrollbar_width = 7
  scrollbar_thumb = "sb_th_*.png"
}
`
	tplValues := TplValues{}
	err := json.Unmarshal([]byte(testThemeTplJSON), &tplValues)
	if err != nil {
		c.Error(err)
	}

	s, _ := theme.getCustomizedThemeContent([]byte(testThemeTplContent), tplValues)
	c.Check(string(s), Equals, wantThemeTxtContent)
}
