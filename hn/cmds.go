package hn

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/urfave/cli/v2"
)

func cmdSyncFiles(ctx *cli.Context) error {
	c := loadAllFiles(false)
	for _, file := range c.Files {
		if file.IsPartial && len(file.Items) == ItemsPerFile {
			must(0, saveDataFile(file))
			must(0, os.Remove(file.Path))
		}
	}
	maxItem := must(client.MaxItem())
	startLast := maxItem / ItemsPerFile * ItemsPerFile

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
		file.Items = append(file.Items, item)

		if i%SaveCheckpoint == 0 {
			must(0, saveDataFile(file))
		}
	}
}
