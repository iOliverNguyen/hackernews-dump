package hn

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func cmdSyncFiles(ctx *cli.Context) error {
	if flagDebug {
		setupDebug()
	}

	c := loadAllFiles(false)
	for _, file := range c.Files {
		if file.IsPartial && len(file.Items) == ItemsPerFile {
			must(0, saveDataFile(file))
			must(0, os.Remove(file.Path))
		}
	}
	maxItem := must(client.MaxItem())
	startLast := maxItem / ItemsPerFile * ItemsPerFile

	if UpdateBack > 0 && len(c.Files) > 0 {
		lastFile := c.Files[len(c.Files)-1]
		if !lastFile.IncludeData {
			c.Files[len(c.Files)-1] = must(loadDataFile(lastFile.Path, true))
			c.FileByNumber[lastFile.EndNumber] = c.Files[len(c.Files)-1]
			lastFile = c.Files[len(c.Files)-1]
		}
		t0 := lastFile.LastTime()
		t1 := t0.Add(-time.Duration(UpdateBack) * 24 * time.Hour)

		for i := len(c.Files) - 1; i >= 0; i-- {
			file := c.Files[i]
			if !file.IncludeData {
				c.Files[i] = must(loadDataFile(file.Path, true))
				c.FileByNumber[file.EndNumber] = c.Files[i]
				file = c.Files[i]
			}
			if file.LastTime().After(t1) {
				c.Files[i] = nil
				c.FileByNumber[file.EndNumber] = nil
			} else {
				break
			}
		}
	}

	wg := &sync.WaitGroup{}
	ch := make(chan int, Concurrent)
	for i := startLast; i > 0; i -= ItemsPerFile {
		ch <- i
		wg.Add(1)
		done := func() {
			<-ch
			wg.Done()
		}
		file := c.FileByNumber[i]
		if file != nil && !file.IsPartial {
			done()
			continue
		}
		if file == nil {
			file = &DataFile{
				EndNumber:   i,
				IsPartial:   true,
				IncludeData: true,
			}
		}
		fmt.Printf("loading %v-%v\n", file.EndNumber-ItemsPerFile+1, file.EndNumber)
		go func(file *DataFile) {
			defer done()
			loadItemsAndSaveToFile(file)
		}(file)
	}
	wg.Wait()
	return nil
}

func loadItemsAndSaveToFile(file *DataFile) {
	first := file.EndNumber - ItemsPerFile + 1
	last := first - 1
	if len(file.Items) > 0 {
		last = file.Items[len(file.Items)-1].ID
	}
	for i := last + 1; i <= file.EndNumber; i++ {
		item, err := client.GetItem(i)
		if err != nil {
			fmt.Println("error loading item (will try again)", i, err)
			i--
			time.Sleep(5 * time.Second)
			continue
		}
		if item == nil { // done
			return
		}
		file.Items = append(file.Items, *item)

		if i%SaveCheckpoint == 0 {
			must(0, saveDataFile(file))
		}
	}
}

func cmdLoadMem(ctx *cli.Context) error {
	t0 := time.Now()
	fmt.Println("start load", t0)

	c := loadAllFiles(true)
	t1 := time.Now()
	fmt.Println("load data", t1.Sub(t0))

	processData(c)
	t2 := time.Now()
	fmt.Println("fill root", t2.Sub(t1))
	return nil
}
