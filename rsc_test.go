package main

import (
	//"log"

	"log"
	"math"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveDataIntoFile(t *testing.T) {
	// t.SkipNow()

	fileInfo := &FileInfo{Filename: "testCreatefile", Filepath: "localdir_test", Filedata: "Test create data"}

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	//Create directory with proper permissions, if need be
	pathErr := os.MkdirAll(fileInfo.Filepath, 0755)
	//Log path error & return error status code
	if pathErr != nil {
		log.Println(pathErr)
	}

	fileInfo.Filepath = filepath.Join(currDir, fileInfo.Filepath+"/"+fileInfo.Filename)

	if err := fileInfo.saveDataIntoFile(); err != nil {
		log.Println(err)
		t.Errorf("File saving returned an error: %s", err.Error())
	}
}

func TestUpdateDataInFile(t *testing.T) {
	// t.SkipNow()

	fileInfo := &FileInfo{Filename: "testUpdatefile", Filepath: "localdir_test", Filedata: "Test update data"}

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	//Create directory with proper permissions, if need be
	pathErr := os.MkdirAll(fileInfo.Filepath, 0755)
	//Log path error & return error status code
	if pathErr != nil {
		log.Println(pathErr)
	}

	fileInfo.Filepath = filepath.Join(currDir, fileInfo.Filepath+"/"+fileInfo.Filename)

	if err := fileInfo.saveDataIntoFile(); err != nil {
		log.Println(err)
		t.Errorf("File saving returned an error: %s", err.Error())
	}

	fileInfo.Filedata = "Test update data 01"
	if err := fileInfo.replaceDataInFile(); err != nil {
		log.Println(err)
		t.Errorf("File updating returned an error: %s", err.Error())
	}
}

func TestReadDataFromFile(t *testing.T) {
	// t.SkipNow()

	fileData := "Test read data"
	fileInfo := &FileInfo{Filename: "testReadfile", Filepath: "localdir_test", Filedata: fileData}

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Fatal(err)
	}
	//Create directory with proper permissions, if need be
	pathErr := os.MkdirAll(fileInfo.Filepath, 0755)
	//Log path error & return error status code
	if pathErr != nil {
		log.Println(pathErr)
	}

	fileInfo.Filepath = filepath.Join(currDir, fileInfo.Filepath+"/"+fileInfo.Filename)

	if err := fileInfo.saveDataIntoFile(); err != nil {
		log.Println(err)
		t.Errorf("File saving returned an error: %s", err.Error())
	}

	if err := fileInfo.readAllDataFromFile(); err != nil {
		log.Println(err)
		t.Errorf("File reading returned an error: %s", err.Error())
	}

	if fileData != fileInfo.Filedata {
		t.Errorf("File reading returned data was incorrect, got: %s, want: %s.", fileInfo.Filedata, fileData)
	}
}

func TestCalculateAvgAndStdDev(t *testing.T) {
	// t.SkipNow()

	arrayInts := []int{2, 3, 4}
	actMean := float64(3)
	actStddev := float64(0.81649658092773)

	mean, stddev := calculateAvgAndStdDev(arrayInts)

	isMeanDiff := float64Equal(actMean, mean)
	isStddevDiff := float64Equal(actStddev, stddev)

	if !isMeanDiff {
		t.Errorf("Mean calcution returned data was incorrect, got: %f, want: %f.", mean, actMean)
	}

	if !isStddevDiff {
		t.Errorf("Standard Deviation calcution returned data was incorrect, got: %f, want: %f.", stddev, actStddev)
	}
}

func float64Equal(a, b float64) bool {
	ba := math.Float64bits(a)
	bb := math.Float64bits(b)
	diff := ba - bb
	if diff < 0 {
		diff = -diff
	}
	// accept 3 bits difference
	return diff < 100
}
