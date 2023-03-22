package hn

import (
	"sync"
	"time"
)

type Time int

func (t Time) GoTime() time.Time {
	return time.Unix(int64(t), 0)
} // Unix time

type List []*Item

type Item struct {
	ID          int    `json:"id,omitempty"`          // The item's unique id.
	Deleted     bool   `json:"deleted,omitempty"`     // true if the item is deleted.
	Type        string `json:"type,omitempty"`        // The type of item. One of "job", "story", "comment", "poll", or "pollopt".
	By          string `json:"by,omitempty"`          // The username of the item's author.
	Time        Time   `json:"time,omitempty"`        // Creation date of the item, in Unix Time.
	Text        string `json:"text,omitempty"`        // The comment, story or poll text. HTML.
	Dead        bool   `json:"dead,omitempty"`        // true if the item is dead.
	Parent      int    `json:"parent,omitempty"`      // The comment's parent: either another comment or the relevant story.
	Poll        int    `json:"poll,omitempty"`        // The pollopt's associated poll.
	Kids        []int  `json:"kids,omitempty"`        // The ids of the item's comments, in ranked display order.
	URL         string `json:"url,omitempty"`         // The URL of the story.
	Score       int    `json:"score,omitempty"`       // The story's score, or the votes for a pollopt.
	Title       string `json:"title,omitempty"`       // The title of the story, poll or job. HTML.
	Parts       []int  `json:"parts,omitempty"`       // A list of related pollopts, in display order.
	Descendants int    `json:"descendants,omitempty"` // In the case of stories or polls, the total comment count

	// Extra fields
	Root int `json:"root,omitempty"`
}

type DataFile struct {
	EndNumber int
	Path      string
	Items     []Item

	// flag for loading file with .partial.json or exclude data
	IsPartial   bool
	IncludeData bool
}

func (c *DataFile) LastTime() time.Time {
	return c.Items[len(c.Items)-1].Time.GoTime()
}

type FileCollection struct {
	sync.Mutex

	Files        []*DataFile
	FileByNumber map[int]*DataFile
	MaxNumber    int
}

func (c *FileCollection) GetFile(id int) *DataFile {
	startNumber := (id-1)/ItemsPerFile*ItemsPerFile + 1
	endNumber := startNumber + ItemsPerFile - 1
	return c.FileByNumber[endNumber]
}

func (c *FileCollection) GetItem(id int) *Item {
	file := c.GetFile(id)
	if file == nil {
		return nil
	}
	idx := id - file.EndNumber + ItemsPerFile - 1
	if idx >= len(file.Items) {
		return nil
	}
	return &file.Items[idx]
}

func (c *FileCollection) Add(file *DataFile) {
	c.Lock()
	defer c.Unlock()

	if c.FileByNumber[file.EndNumber] == nil {
		c.Files = append(c.Files, file)
		c.FileByNumber[file.EndNumber] = file
	}
	if file.EndNumber > c.MaxNumber {
		c.MaxNumber = file.EndNumber
	}
}
