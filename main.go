package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
	
	"github.com/alexflint/go-arg"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/dustin/go-humanize"
)

const (
	defTemplate = "{{.title}} [GOG]"
	sanRegexStr = `[\/:*?"><|]`
	siteUrl = "https://www.gog.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/"+
		"537.36 (KHTML, like Gecko) Chrome/107.0.0.0 Safari/537.36"
)

var (
	jar, _ = cookiejar.New(nil)
	client = &http.Client{Transport: &Transport{}, Jar: jar}
)

var resolvePlatform = map[string]string{
	"windows": "1,2,4,8,4096,16384",
	"linux": "1024,2048,8192",
	"mac": "16,32",
}

var languages = []string{
	"en", "cz", "de", "es", "fr", "it",
	"hu", "nl", "pl", "pt", "br", "sv",
	"tr", "uk", "ru", "ar", "ko", "cn",
	"jp", "all",
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Referer", siteUrl+"/")
	req.Header.Add("Origin", siteUrl)
	return http.DefaultTransport.RoundTrip(req)
}


func (wc *WriteCounter) Write(p []byte) (int, error) {
	var speed int64 = 0
	n := len(p)
	wc.Downloaded += int64(n)
	percentage := float64(wc.Downloaded) / float64(wc.Total) * float64(100)
	wc.Percentage = int(percentage)
	toDivideBy := time.Now().UnixMilli() - wc.StartTime
	if toDivideBy != 0 {
		speed = int64(wc.Downloaded) / toDivideBy * 1000
	}
	fmt.Printf("\r%d%% @ %s/s, %s/%s ", wc.Percentage, humanize.Bytes(uint64(speed)),
		humanize.Bytes(uint64(wc.Downloaded)), wc.TotalStr)
	return n, nil
}

func handleErr(errText string, err error, _panic bool) {
	errString := errText + "\n" + err.Error()
	if _panic {
		panic(errString)
	}
	fmt.Println(errString)
}

func wasRunFromSrc() bool {
	buildPath := filepath.Join(os.TempDir(), "go-build")
	return strings.HasPrefix(os.Args[0], buildPath)
}

func getScriptDir() (string, error) {
	var (
		ok    bool
		err   error
		fname string
	)
	runFromSrc := wasRunFromSrc()
	if runFromSrc {
		_, fname, _, ok = runtime.Caller(0)
		if !ok {
			return "", errors.New("Failed to get script filename.")
		}
	} else {
		fname, err = os.Executable()
		if err != nil {
			return "", err
		}
	}
	return filepath.Dir(fname), nil
}

func readConfig() (*Config, error) {
	data, err := os.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	var obj Config
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func parseArgs() *Args {
	var args Args
	arg.MustParse(&args)
	return &args
}

func makeDirs(path string) error {
	err := os.MkdirAll(path, 0755)
	return err
}

func checkLang(userLang string) bool {
	for _, lang := range languages {
		if lang == userLang {
			return true
		}
	}
	return false
}

func parseCfg() (*Config, error) {
	cfg, err := readConfig()
	if err != nil {
		return nil, err
	}

	args := parseArgs()
	query := strings.TrimSpace(args.Query)
	if query != "" && len(query) < 3 {
		return nil, errors.New("query must be at least two characters")
	}

	cfg.Query = args.Query
	if args.Platform != "" {
		cfg.Platform = args.Platform
	}
	if args.Language != "" {
		cfg.Language = args.Language
	}
	if args.FolderTemplate != "" {
		cfg.FolderTemplate = args.FolderTemplate
	}
	if args.Goodies {
		cfg.Goodies = args.Goodies 
	}

	if strings.TrimSpace(cfg.Platform) == "" || strings.TrimSpace(cfg.Language) == "" {
		return nil, errors.New("platform and language are required")
	}

	lang := strings.ToLower(cfg.Language)
	if !checkLang(lang) {
		return nil, errors.New("invalid language: " + cfg.Language)
	}
	args.Language = lang

	platform := strings.ToLower(cfg.Platform)
	if platform == "win" {
		platform = "windows"
	} else if platform == "osx" {
		platform = "mac"
	}
	platformIds, ok := resolvePlatform[platform]
	if !ok {
		return nil, errors.New("invalid platform: " + cfg.Platform)
	}
	cfg.Platform = platform
	cfg.PlatformIDs = platformIds
	if args.OutPath != "" {
		cfg.OutPath = args.OutPath
	}
	if cfg.OutPath == "" {
		cfg.OutPath = "GOG downloads"
	}
	if strings.TrimSpace(cfg.FolderTemplate) == "" {
		cfg.FolderTemplate = defTemplate
	}

	return cfg, nil
}

func readCookies() ([]*Cookie, error) {
	data, err := os.ReadFile("cookies.json")
	if err != nil {
		return nil, err
	}

	var obj []*Cookie
	err = json.Unmarshal(data, &obj)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func setCookies(_cookies []*Cookie) error {
	var cookies []*http.Cookie
	for _, _cookie := range _cookies {
		name := _cookie.Name
		value := _cookie.Value
		cookie := &http.Cookie{
			Domain: _cookie.Domain,
			Name:   name,
			Path:   _cookie.Path,
			Secure: _cookie.Secure,
			Value:  value,
		}
		cookies = append(cookies, cookie)
	}

	u, err := url.Parse(siteUrl)
	if err != nil {
		return err
	}

	client.Jar.SetCookies(u, cookies)
	return nil
}

func checkCookies() (bool, error) {
	req, err := client.Get(siteUrl+"/userData.json")
	if err != nil {
		return false, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return false, errors.New(req.Status)
	}
	var obj UserData
    err = json.NewDecoder(req.Body).Decode(&obj)
    if err != nil {
    	return false, err
    }
    ok := obj.IsLoggedIn
    if ok {
    	fmt.Println("Signed in as " + obj.Username + ".\n")
    }
    return ok, nil
}

func search(queryStr, platformIds, lang string) ([]Product, error) {
	req, err := http.NewRequest(
		http.MethodGet, siteUrl+"/account/getFilteredProducts", nil)
	if err != nil {
		return nil, err
	}

	pageNum := 1
	query := url.Values{}
	query.Set("hiddenFlag", "0")
	if lang != "all" {
		query.Set("language", lang)
	}
	query.Set("mediaType", "1")
	query.Set("sortBy", "date_purchased")
	if queryStr != "" {
		query.Set("search", queryStr)
	}
	query.Set("system", platformIds)
	query.Set("totalPages", "1")
	var products []Product

	for {
		query.Set("page", strconv.Itoa(pageNum))
		req.URL.RawQuery = query.Encode()

		do, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if do.StatusCode != http.StatusOK {
			do.Body.Close()
			return nil, errors.New(do.Status)
		}

		var obj Search
		err = json.NewDecoder(do.Body).Decode(&obj)
		do.Body.Close()
		if err != nil {
			return nil, err
		}

		if obj.TotalPages == 0 {
			break
		}

		products = append(products, obj.Products...)
		if pageNum == obj.TotalPages {
			break
		}

		pageNum++
		time.Sleep(time.Second*1)
	}
	return products, nil
}

func getUserGameIdx(products []Product, prodLen int) (int, error) {
	var (
		idx int
		opts []string
	)
	for _, p := range products {
		opts = append(opts, p.Title)
	}
	prompt := &survey.Select{Options: opts}
	err := survey.AskOne(
		prompt, &idx, survey.WithValidator(survey.Required),
		survey.WithPageSize(10))
	if err != nil {
		return 0, err
	}
	return idx, nil
}

func selectGameId(products []Product, queryStr string) (int, error) {
	prodLen := len(products)
	if prodLen == 1 {
		return products[0].ID, nil
	}
	if queryStr != "" {
		fmt.Println("Your search yielded more than one result.")
	}
	idx, err := getUserGameIdx(products, prodLen)
	if err != nil {
		return 0, err
	}
	return products[idx].ID, nil
}

func getGameMeta(id int) (*GameMeta, error) {
	req, err := client.Get(
		siteUrl + "/account/gameDetails/" + strconv.Itoa(id) + ".json")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()
	if req.StatusCode != http.StatusOK {
		return nil, errors.New(req.Status)
	}

	var obj GameMeta
    err = json.NewDecoder(req.Body).Decode(&obj)
    if err != nil {
    	return nil, err
    }
    return &obj, nil
}

func parseDownloads(meta *GameMeta, platform string, goodies bool) ([]*Download, error) {
	var parsedDloads []*Download
	// Shambles, get structs working.
	downloads := meta.Downloads[0][1].(map[string]interface{})[platform].([]interface{})
	for _, _d := range downloads {
		d := _d.(map[string]interface{})
		ver := "<no ver>"
		if d["version"] != nil {
			ver = d["version"].(string)
		}
		parsedDload := &Download{
			ManualURL: siteUrl + d["manualUrl"].(string),
			Name: 	   d["name"].(string), 
			Version:   ver,
			Date:      d["date"].(string),
			Size:      d["size"].(string),
		}
		parsedDloads = append(parsedDloads, parsedDload)
	}

	if goodies {
		for _, e := range meta.Extras {
			e.ManualURL = siteUrl + e.ManualURL
			parsedDloads = append(parsedDloads, e)
		}		
	}

	return parsedDloads, nil
}

func getLongestNameLen(downloads []*Download) int {
	var longest int
	for _, d := range downloads {
		curLen := len(d.Name)
		if curLen > longest {
			longest = curLen
		}
	}
	return longest
}

func getUserDloadIndexes(downloads []*Download) ([]int, error) {
	var (
		indexes []int
		opts []string
	)
	longestNameLen := getLongestNameLen(downloads)

	for _, d := range downloads {
		ver := d.Version
		if ver == "" {
			ver = "<no ver>"
		} 
		spaces := strings.Repeat(" ", longestNameLen-len(d.Name))
		opts = append(opts, d.Name + spaces + " - " + ver + ", " + d.Size)
	}

	prompt := &survey.MultiSelect{Options: opts}
	err := survey.AskOne(
		prompt, &indexes, survey.WithValidator(survey.Required),
		survey.WithPageSize(10))
	if err != nil {
		return nil, err
	}
	return indexes, nil
}

func selectDownloads(downloads []*Download) ([]*Download, error) {
	// prodLen := len(products)
	// if prodLen == 1 {
	// 	return products[0].ID, nil
	// }
	var selectedDloads []*Download
	indexes, err := getUserDloadIndexes(downloads)
	if err != nil {
		return nil, err
	}
	for _, idx := range indexes {
		selectedDloads = append(selectedDloads, downloads[idx])
	}
	return selectedDloads, nil
}

func fileExists(path string) (bool, int64, error) {
	f, err := os.Stat(path)
	if err == nil {
		return !f.IsDir(), f.Size(), nil
	} else if os.IsNotExist(err) {
		return false, 0, nil
	}
	return false, 0, err
}

func sanitise(fname string) string {
	regex := regexp.MustCompile(sanRegexStr)
	fname = regex.ReplaceAllString(fname, "_")
	return fname
}

func getFname(itemUrl string) (string, error) {
	req, err := client.Head(itemUrl)
	if err != nil {
		return "", err
	}
	if req.StatusCode != http.StatusOK {
		return "", errors.New(req.Status)
	}

	fname := path.Base(req.Request.URL.String())
	fname, err = url.PathUnescape(fname)
	if err != nil {
		return "", err
	}
	return fname, nil
}

func getBase(fname string) string {
	ext := filepath.Ext(fname)
	if ext == ".gz" {
		ext = ".tar.gz"
	}
	base := fname[:len(fname)-len(ext)]
	return base
}

func parseTempMeta(title string) map[string]string {
	parsed := map[string]string{
		"title": title,
		"titlePeriods": strings.ReplaceAll(title, " ", "."),
	}
	return parsed
}

func parseTemplate(text string, meta map[string]string) string {
	var buffer bytes.Buffer
	for {
		err := template.Must(template.New("").Parse(text)).Execute(&buffer, meta)
		if err == nil {
			break
		}
		fmt.Println("Failed to parse template. Default will be used instead.")
		text = defTemplate
		buffer.Reset()
	}
	return html.UnescapeString(buffer.String())
}

func downloadItem(download *Download, outPath string) error {
	var startByte int64
	itemUrl := download.ManualURL

	fname, err := getFname(itemUrl)
	if err != nil {
		return err
	}

	outPath = filepath.Join(outPath, fname)
	exists, _, err := fileExists(outPath)
	if err != nil {
		return err
	}
	if exists {
		fmt.Println("Item already exists locally.")
		return nil
	}

	base := getBase(outPath)
	incompPath := filepath.Join(base + ".incomplete")

	exists, size, err := fileExists(incompPath)
	if err != nil {
		return err
	}
	if exists {
		startByte = size
		fmt.Println("Incomplete item exists locally. Resuming...")
	}

	req, err := http.NewRequest(http.MethodGet, download.ManualURL, nil)
	if err != nil {
		return err
	}
	req.Header.Add(
		"Range", "bytes=" + strconv.FormatInt(startByte, 10) + "-")

	do, err := client.Do(req)
	if err != nil {
		return err
	}
	defer do.Body.Close()
	if do.StatusCode != http.StatusOK && do.StatusCode != http.StatusPartialContent {
		return errors.New(do.Status)
	}

	f, err := os.OpenFile(incompPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0755)
	if err != nil {
		return err
	}

	totalBytes := do.ContentLength + startByte
	counter := &WriteCounter{
		Total:     totalBytes,
		TotalStr:  humanize.Bytes(uint64(totalBytes)),
		StartTime: time.Now().UnixMilli(),
		Downloaded: startByte,
	}

	_, err = io.Copy(f, io.TeeReader(do.Body, counter))
	f.Close()
	err = os.Rename(incompPath, outPath)
	fmt.Println("")
	return err
}

func init() {
	fmt.Println(`
 _____ _____ _____    ____                _           _         
|   __|     |   __|  |    \ ___ _ _ _ ___| |___ ___ _| |___ ___ 
|  |  |  |  |  |  |  |  |  | . | | | |   | | . | .'| . | -_|  _|
|_____|_____|_____|  |____/|___|_____|_|_|_|___|__,|___|___|_|
`)
}

func main() {
	scriptDir, err := getScriptDir()
	if err != nil {
		panic(err)
	}

	err = os.Chdir(scriptDir)
	if err != nil {
		panic(err)
	}

	cfg, err := parseCfg()
	if err != nil {
		handleErr("failed to parse config/args", err, true)
	}

	err = makeDirs(cfg.OutPath)
	if err != nil {
		handleErr("failed to make output folder(s)", err, true)
	}

	cookies, err := readCookies()
	if err != nil {
		handleErr("failed to read cookies", err, true)
	}

	err = setCookies(cookies)
	if err != nil {
		handleErr("failed to set cookies", err, true)
	}

	ok, err := checkCookies()
	if err != nil {
		handleErr("failed to check cookies", err, true)
	}
	if !ok {
		panic("bad cookies")
	}

	products, err := search(cfg.Query, cfg.PlatformIDs, cfg.Language)
	if err != nil {
		panic(err)
	}
	if len(products) == 0 {
		fmt.Println("No search results.")
		os.Exit(1)
	}

	id, err := selectGameId(products, cfg.Query)
	if err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
		}
		handleErr("failed to select game id", err, true)
	}

	gameMeta, err := getGameMeta(id)
	if err != nil {
		handleErr("failed to get game meta", err, true)
	}
	fmt.Println("--" + gameMeta.Title + "--")

	downloads, err := parseDownloads(gameMeta, cfg.Platform, cfg.Goodies)
	if err != nil {
		handleErr("failed to parse items", err, true)
	}

	downloads, err = selectDownloads(downloads)
	if err != nil {
		if err == terminal.InterruptErr {
			os.Exit(0)
		}
		handleErr("failed to select items", err, true)
	}

	itemTotal := len(downloads)
	templateMeta := parseTempMeta(gameMeta.Title)
	template := parseTemplate(cfg.FolderTemplate, templateMeta)

	outPath := filepath.Join(cfg.OutPath, sanitise(template))
	err = makeDirs(outPath)
	if err != nil {
		handleErr("failed to make game folder", err, true)
	}

	for i, item := range downloads {
		fmt.Printf("Item %d of %d:\n", i+1, itemTotal)
		fmt.Println(item.Name)

		err = downloadItem(item, outPath)
		if err != nil {
			handleErr("failed to download item", err, false)
		}
	}
}