package hn

func processData(c *FileCollection) {
	for _, file := range c.Files {
		for _, item := range file.Items {
			if item.Parent != 0 {
				parent := c.GetItem(item.Parent)
				if parent != nil {
					item.Root = parent.Root
				}
			}
		}
	}
}
