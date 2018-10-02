package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"unicode/utf8"

	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Token is a custom type of alphabetic word
type Token []rune

var alphaNumRegExp = regexp.MustCompile("[a-zA-Z0-9]+") //for alphanumeric characters

type FileInfo struct {
	Filename string `json:"filename"`
	Filepath string `json:"filepath"`
	Filedata string `json:"filedata"`
}

type DirStats struct {
	NumFiles           int     `json:"total_number_of_files"`
	NumBytes           int64   `json:"total_number_of_bytes"`
	AvgAlphanumChar    float64 `json:"avg_alphanum_char"`
	StdDevAlphanumChar float64 `json:"stddev_alphanum_char"`
	AvgWordLen         float64 `json:"avg_wordlength_char"`
	StdDevWordLen      float64 `json:"stddev_wordlength_char"`
}

// Create File: Entry point delegate function
func createFile(ctx *fasthttp.RequestCtx) {
	filename := ctx.UserValue("filename").(string)
	fmt.Fprintf(ctx, "filename, %s!\n", filename)
	recReqFileInfo := &FileInfo{Filename: filename}
	httpStatusCode := recReqFileInfo.createLocalFile(ctx)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("Create File returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

// Creates file and returns HTTP Response Status Code
func (fileInfo *FileInfo) createLocalFile(ctx *fasthttp.RequestCtx) (respStatusCode int) {

	if err := fileInfo.fetchFilePath(ctx); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}

	if err := fileInfo.saveDataIntoFile(); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}

	return fasthttp.StatusOK
}

// saved fileData into the file
func (fileInfo *FileInfo) saveDataIntoFile() error {
	fileHandle, err := os.OpenFile(fileInfo.Filepath, os.O_CREATE|os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fileHandle.Close()

	fileHandle.WriteString(fileInfo.Filedata)
	return nil
}

// Update File: Entry point delegate function
func updateFile(ctx *fasthttp.RequestCtx) {
	filename := ctx.UserValue("filename").(string)
	fmt.Fprintf(ctx, "filename, %s!\n", filename)
	recReqFileInfo := &FileInfo{Filename: filename}
	httpStatusCode := recReqFileInfo.updateLocalFile(ctx)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("Update File returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

// Updates file and returns HTTP Response Status Code
func (fileInfo *FileInfo) updateLocalFile(ctx *fasthttp.RequestCtx) (respStatusCode int) {

	if err := fileInfo.fetchFilePath(ctx); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}

	if err := fileInfo.replaceDataInFile(); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}

	return fasthttp.StatusOK
}

// Replaces fileData in the file
func (fileInfo *FileInfo) replaceDataInFile() error {
	fileHandle, err := os.OpenFile(fileInfo.Filepath, os.O_WRONLY|os.O_SYNC|os.O_TRUNC, 0666)
	if err != nil {
		log.Println(err)
		return err
	}
	defer fileHandle.Close()

	fileHandle.WriteString(fileInfo.Filedata)
	return nil
}

// Read File: Entry point delegate function
func readFile(ctx *fasthttp.RequestCtx) {
	filename := ctx.UserValue("filename").(string)

	recReqFileInfo := &FileInfo{Filename: filename}

	httpStatusCode := recReqFileInfo.readLocalFile(ctx)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("Read File returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		respJSONBytes, _ := json.Marshal(recReqFileInfo.Filedata)
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(respJSONBytes)
	}
}

// Reads file contents and returns the file contents & HTTP Response Status Code
func (fileInfo *FileInfo) readLocalFile(ctx *fasthttp.RequestCtx) (respStatusCode int) {
	if err := fileInfo.fetchFilePath(ctx); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}
	if err := fileInfo.readAllDataFromFile(); err != nil {
		log.Println(err)
		return fasthttp.StatusInternalServerError
	}

	return fasthttp.StatusOK
}

// Reads all File data and returns File data & error
func (fileInfo *FileInfo) readAllDataFromFile() error {
	data, err := ioutil.ReadFile(fileInfo.Filepath)
	if err != nil {
		log.Println(err)
		fileInfo.Filedata = ""
		return err
	}
	fileInfo.Filedata = string(data)
	return nil
}

// Delete File: Entry point delegate function
func deleteFile(ctx *fasthttp.RequestCtx) {
	filename := ctx.UserValue("filename").(string)
	recReqFileInfo := &FileInfo{Filename: filename}
	if err := recReqFileInfo.fetchFilePath(ctx); err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	// delete the File
	if err := os.Remove(recReqFileInfo.Filepath); err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
	} else {
		ctx.SetStatusCode(fasthttp.StatusOK)
	}
}

// Read Directory Contents: Entry point delegate function
func readDirContents(ctx *fasthttp.RequestCtx) {
	var recReqFileInfo FileInfo
	if err := json.Unmarshal(ctx.PostBody(), &recReqFileInfo); err != nil {
		log.Println(err)
	}

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusNotFound)
		return
	}

	filePath := filepath.Join(currDir, recReqFileInfo.Filepath)

	dirStats, httpStatusCode := readLocalDir(ctx, filePath)
	if fasthttp.StatusOK != httpStatusCode {
		fmt.Errorf("Read Directory contents returned http StatusCode was incorrect, got: %d.", httpStatusCode)
		ctx.SetStatusCode(httpStatusCode)
	} else {
		respJSONBytes, _ := json.Marshal(dirStats)
		ctx.SetStatusCode(fasthttp.StatusOK)
		ctx.SetContentType("application/json")
		ctx.SetBody(respJSONBytes)
	}
}

func readLocalDir(ctx *fasthttp.RequestCtx, dirname string) (dirStats *DirStats, respStatusCode int) {

	var filesPath []string
	var filesInfo []os.FileInfo
	wordsPerFile := make([]int, 0)
	wordsLenPerDir := make([]int, 0)
	alphaNumsPerFile := make([]int, 0)

	err := filepath.Walk(dirname, visitDir(&filesPath, &filesInfo))
	if err != nil {
		return nil, fasthttp.StatusNotFound
	}

	numFiles := len(filesInfo)
	var numBytes int64
	for i := 0; i < len(filesInfo); i++ {
		fileInfo := filesInfo[i]
		filePath := filesPath[i]

		numBytes += fileInfo.Size()

		wordsLen, charsCount, wordsCount := fetchFileStats(filePath)

		wordsPerFile = append(wordsPerFile, wordsCount)
		alphaNumsPerFile = append(alphaNumsPerFile, charsCount)
		wordsLenPerDir = append(wordsLenPerDir, wordsLen...)
	}

	avgAlphanumChar, stdDevAlphanumChar := calculateAvgAndStdDev(alphaNumsPerFile)
	avgWordLen, stdDevWordLen := calculateAvgAndStdDev(wordsLenPerDir)
	dirStats = NewDirStats(numFiles, numBytes, avgAlphanumChar, stdDevAlphanumChar, avgWordLen, stdDevWordLen)
	fmt.Printf("DirStats are:\n%+v\n", dirStats)

	return dirStats, fasthttp.StatusOK
}

func visitDir(filesPath *[]string, filesInfo *[]os.FileInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		*filesPath = append(*filesPath, path)
		*filesInfo = append(*filesInfo, info)
		return nil
	}
}

func fetchFileStats(filePath string) (wordsLen []int, charsCount, wordsCount int) {
	fileHandle, err := os.OpenFile(filePath, os.O_RDONLY|os.O_SYNC, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer fileHandle.Close()

	scanner := bufio.NewScanner(fileHandle)
	allAlphaNumsTokens := fetchAlphabeticTokens(scanner)

	return fetchTokensStats(allAlphaNumsTokens)
}

// It fetches only Alphabetic token
func (token Token) fetchAlphabeticToken() (alphaNumToken Token) {
	alphaNumToken = Token(alphaNumRegExp.FindString(string(token)))
	if nil == alphaNumToken || 0 == len(alphaNumToken) {
		alphaNumToken = nil
	}

	return alphaNumToken
}

// It returns all the Alphabetic Tokens
func fetchAlphabeticTokens(scanner *bufio.Scanner) (allAlphaNumsTokens []Token) {
	var alphaNumToken Token

	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		alphaNumToken = Token(scanner.Text()).fetchAlphabeticToken()
		allAlphaNumsTokens = append(allAlphaNumsTokens, alphaNumToken)
	}
	return allAlphaNumsTokens
}

// It fetches token array stats
func fetchTokensStats(allTokens []Token) (wordsLen []int, charsCount, wordsCount int) {
	var tokenChars int
	charsCount = 0
	wordsLen = make([]int, 0)
	for _, token := range allTokens {
		tokenChars = utf8.RuneCountInString(string(token))
		charsCount += tokenChars
		wordsLen = append(wordsLen, tokenChars)
	}
	wordsCount = len(allTokens)
	return wordsLen, charsCount, wordsCount
}

func calculateAvgAndStdDev(array []int) (avg, stddev float64) {
	var sum int
	total := len(array)
	for _, elem := range array {
		sum += elem
	}
	avg = float64(sum) / float64(total)

	for _, elem := range array {
		// The use of Pow math function func Pow(x, y float64) float64
		stddev += math.Pow(float64(elem)-avg, 2)
	}
	// The use of Sqrt math function func Sqrt(x float64) float64
	stddev = math.Sqrt(stddev / float64(total))
	return avg, stddev
}

func (fileInfo *FileInfo) fetchFilePath(ctx *fasthttp.RequestCtx) error {
	// fmt.Fprintf(ctx, "fileinfo, %s!\n", ctx.PostBody())

	if err := json.Unmarshal(ctx.PostBody(), &fileInfo); err != nil {
		log.Println(err)
		fileInfo.Filepath = ""
		return err
	}

	currDir, err := filepath.Abs("./")
	if err != nil {
		log.Println(err)
		fileInfo.Filepath = ""
		return err
	}

	fileInfo.Filepath = filepath.Join(currDir, fileInfo.Filepath+"/"+fileInfo.Filename)
	return nil
}

func NewDirStats(numFiles int, numBytes int64, avgAlphanumChar, stdDevAlphanumChar, avgWordLen, stdDevWordLen float64) *DirStats {
	return &DirStats{numFiles, numBytes, avgAlphanumChar, stdDevAlphanumChar, avgWordLen, stdDevWordLen}
}

func main() {

	// curl -v #verbose

	// curl -H "Content-Type: application/json" -X POST -d '{"filepath":"localfile","filedata":"CDA8777865C7CC3C"}' http://localhost:8081/api/file/newFile02
	// curl -H "Content-Type: application/json" -X GET -d '{"filepath":"localfile"}' http://localhost:8081/api/file/newFile02
	// curl -H "Content-Type: application/json" -X PUT -d '{"filepath":"localfile","filedata":"new Test 02a CDA8777865C7CC3C"}' http://localhost:8081/api/file/newFile02
	// curl -H "Content-Type: application/json" -X DELETE -d '{"filepath":"localfile"}' http://localhost:8081/api/file/newFile03
	// curl -H "Content-Type: application/json" -X GET -d '{"filepath":"localfile"}' http://localhost:8081/api/dir/stats

	router := fasthttprouter.New()
	router.POST("/api/file/:filename", createFile)
	router.GET("/api/file/:filename", readFile)
	router.PUT("/api/file/:filename", updateFile)
	router.DELETE("/api/file/:filename", deleteFile)
	router.GET("/api/dir/stats", readDirContents)

	log.Fatal(fasthttp.ListenAndServe(":8081", router.Handler))
}
