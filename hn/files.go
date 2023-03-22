package hn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func getFileName(fileNumber int) string {
	return spr("%v", fileNumber/ItemsPerFile*ItemsPerFile)
}

func getDirName(fileNumber int) string {
	return spr("%v", fileNumber/ItemsPerFile/FilesPerDir)
}

func saveDataFile(file *DataFile) error {
	if len(file.Items) == 0 {
		return fmt.Errorf("no data to save")
	}

	baseFileName := getFileName(file.EndNumber)
	fileDir := getDirName(file.EndNumber)
	if len(file.Items) > ItemsPerFile {
		panic("invalid items per file")
	}
	fileName := baseFileName + ".jsonl"
	partialFileName := baseFileName + ".partial.jsonl"
	file.IsPartial = len(file.Items) < ItemsPerFile
	if file.IsPartial {
		fileName = partialFileName
	}

	must(0, os.MkdirAll(spr("%v/%v/%v", projectRoot, DataDir, fileDir), 0755))
	file.Path = spr("%v/%v/%v/%v", projectRoot, DataDir, fileDir, fileName)
	partialFilePath := spr("%v/%v/%v/%v", projectRoot, DataDir, fileDir, partialFileName)

	fmt.Printf("save file %v %v%%\n", file.Path, len(file.Items)*100/ItemsPerFile)
	data := must(encodeJSONLines(file.Items))
	err := os.WriteFile(file.Path, data, 0644)
	if err != nil {
		return err
	}
	if !file.IsPartial {
		_ = os.Remove(partialFilePath)
	}
	return nil
}

func loadDataFile(path string, includeData bool) (_ *DataFile, err error) {
	file := &DataFile{
		Path:      path,
		IsPartial: strings.Contains(path, "partial"),
	}
	file.EndNumber, err = strconv.Atoi(strings.Split(filepath.Base(path), ".")[0])
	if err != nil {
		fmt.Println("failed to parse file name", path)
		panic(err)
	}
	if file.IsPartial || includeData {
		file.IncludeData = true
		data := must(os.ReadFile(path))
		file.Items = must(parseJSONLines(data))
		if len(file.Items) > ItemsPerFile {
			return file, fmt.Errorf("invalid items per file %v %v", path, len(file.Items))
		}
	}
	return file, nil
}

func loadAllFiles(includeData bool) *FileCollection {
	c := &FileCollection{
		FileByNumber: map[int]*DataFile{},
	}
	wg := &sync.WaitGroup{}
	ch := make(chan int, Concurrent)
	must(0, filepath.Walk(filepath.Join(projectRoot, DataDir), func(path string, info fs.FileInfo, err error) error {
		must(0, err)
		if !info.IsDir() && strings.Contains(path, "json") {
			ch <- 1
			wg.Add(1)
			go func() {
				defer func() { <-ch }()
				defer wg.Done()
				file := must(loadDataFile(path, includeData))
				c.Add(file)
			}()
		}
		return nil
	}))
	wg.Wait()

	sort.Slice(c.Files, func(i, j int) bool {
		return c.Files[i].EndNumber < c.Files[j].EndNumber
	})
	return c
}

func parseJSONLines(data []byte) (out []Item, _ error) {
	lines := bytes.Split(data, []byte("\n"))
	out = make([]Item, 0, min(len(lines), ItemsPerFile))
	for _, line := range lines {
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}
		var item Item
		if err := json.Unmarshal(line, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}

func encodeJSONLines(items []Item) ([]byte, error) {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	for _, item := range items {
		err := enc.Encode(item)
		if err != nil {
			return nil, err
		}
	}
	return b.Bytes(), nil
}
