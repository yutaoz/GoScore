package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Game struct {
	Home        string
	Visitor     string
	Arena       string
	HomeScore   int
	AwayScore   int
	Attendance  int
	BoxScoreUrl string
	Day         string
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getMonths(url string) []string {
	request, err := http.NewRequest("GET", url, nil)
	check(err)

	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}
	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}

	resContent := string(responseBody)
	sidx := strings.Index(resContent, "filter\">")
	eidx := strings.Index(resContent, "<div id=\"all_schedule")
	resContent = resContent[sidx+len("filter\">") : eidx]
	results := strings.Split(resContent, "<a href=")
	results = results[1:]
	fmt.Println(results[1:])
	fmt.Println(len(results))

	for i, val := range results {
		endidx := strings.Index(val, "\">")
		monthurl := val[1:endidx]
		results[i] = monthurl
		//fmt.Println(monthurl)
	}
	fmt.Println(results)
	return results
}

func getData(baseurl string) []string {
	//baseurl := "https://www.basketball-reference.com/leagues/NBA_1984_games-november.html"
	request, error := http.NewRequest("GET", baseurl, nil)

	if error != nil {
		fmt.Println(error)
	}

	request.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}
	response, error := client.Do(request)

	if error != nil {
		fmt.Println(error)
	}

	responseBody, error := io.ReadAll(response.Body)

	if error != nil {
		fmt.Println(error)
	}

	resContent := string(responseBody)
	sIndex := strings.Index(resContent, "Schedule Table")
	eIndex := strings.Index(resContent, "nonempty_tables_num")
	extractedData := resContent[sIndex+len("Schedule Table") : eIndex]
	sIndex = strings.Index(extractedData, "tbody")
	eIndex = strings.Index(extractedData, "/table")
	extractedData = extractedData[sIndex+len("tbody") : eIndex]

	dataFields := strings.Split(extractedData, "<tr >")
	dataFields = dataFields[1:]
	fmt.Println(len(dataFields))

	// clean up memory after execution
	defer response.Body.Close()

	return dataFields
}

func parseInfo(game string) Game {
	// get Day
	sIndex := strings.Index(game, "year=")
	eIndex := strings.Index(game[sIndex:], "</a>")
	eIndex += sIndex
	day := game[sIndex+len("year=0000> ") : eIndex]
	day = strings.ReplaceAll(day, ",", "")
	//fmt.Println(day)

	// get visitor team
	sIndex = strings.Index(game, "html\">")
	eIndex = strings.Index(game[sIndex:], "</a>")
	eIndex += sIndex
	visitor := game[sIndex+len("html\">") : eIndex]
	//fmt.Println(visitor)

	// get home team
	homeIndex := strings.Index(game, "home_team_name")
	sIndex = strings.Index(game[homeIndex:], "html\">")
	sIndex += homeIndex
	eIndex = strings.Index(game[sIndex:], "</a>")
	eIndex += sIndex
	home := game[sIndex+len("html\">") : eIndex]
	//fmt.Println(home)

	// get visitor points
	sIndex = strings.Index(game, "visitor_pts")
	eIndex = strings.Index(game[sIndex:], "</td>")
	eIndex += sIndex
	visitorPts, _ := strconv.Atoi(game[sIndex+len("visitor_pts\" >") : eIndex])
	//fmt.Println(visitorPts)

	// get home points
	sIndex = strings.Index(game, "home_pts")
	eIndex = strings.Index(game[sIndex:], "</td>")
	eIndex += sIndex
	homePts, _ := strconv.Atoi(game[sIndex+len("home_pts\" >") : eIndex])
	//fmt.Println(homePts)

	// get arena
	sIndex = strings.Index(game, "arena_name")
	eIndex = strings.Index(game[sIndex:], "</td>")
	eIndex += sIndex
	arena := game[sIndex+len("arena_name\" >") : eIndex]
	//fmt.Println(arena)

	// get attendance
	sIndex = strings.Index(game, "attendance")
	eIndex = strings.Index(game[sIndex:], "</td>")
	eIndex += sIndex
	attendance := game[sIndex+len("attendance\" >") : eIndex]
	attendance = strings.ReplaceAll(attendance, ",", "")
	s, err := strconv.Atoi(attendance)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(s)

	// get boxscore url
	bScoreIndex := strings.Index(game, "box_score_text")
	sIndex = strings.Index(game[bScoreIndex:], "href=\"")
	sIndex += bScoreIndex
	eIndex = strings.Index(game[sIndex:], ">")
	eIndex += sIndex
	bscoreurl := game[sIndex+len("href=\"") : eIndex-1]
	//fmt.Println(bscoreurl)

	return Game{home, visitor, arena, homePts, visitorPts, s, bscoreurl, day}

}

func writeData(data []string, month int, year int) {
	fpath := "data/" + strconv.Itoa(year) + "/" + strconv.Itoa(month) + ".txt"
	err := os.MkdirAll(filepath.Dir(fpath), 0755)
	check(err)
	f, err := os.Create(fpath)
	check(err)
	defer f.Close()
	for _, game := range data {
		gamedata := parseInfo(game)

		_, err := f.WriteString(gamedata.Day + "," + gamedata.Home + "," + gamedata.Visitor + "," + strconv.Itoa(gamedata.HomeScore) +
			"," + strconv.Itoa(gamedata.AwayScore) + "," + gamedata.Arena + "," + strconv.Itoa(gamedata.Attendance) + "," + gamedata.BoxScoreUrl + "\n")
		check(err)
	}
}

func main() {
	//rawdata := getData()
	//writeData(rawdata, 11, 1984)
	for i := 1955; i <= 2024; i++ {
		url := "https://www.basketball-reference.com/leagues/NBA_" + strconv.Itoa(i) + "_games.html"
		ms := getMonths(url)
		month := 1
		year := i
		for _, m := range ms {
			burl := "https://www.basketball-reference.com" + m
			rawdata := getData(burl)
			if strings.Index(m, "january") != -1 {
				month = 1
			} else if strings.Index(m, "february") != -1 {
				month = 2
			} else if strings.Index(m, "march") != -1 {
				month = 3
			} else if strings.Index(m, "april") != -1 {
				month = 4
			} else if strings.Index(m, "may") != -1 {
				month = 5
			} else if strings.Index(m, "june") != -1 {
				month = 6
			} else if strings.Index(m, "july") != -1 {
				month = 7
			} else if strings.Index(m, "august") != -1 {
				month = 8
			} else if strings.Index(m, "september") != -1 {
				month = 9
			} else if strings.Index(m, "october") != -1 {
				month = 10
			} else if strings.Index(m, "november") != -1 {
				month = 11
			} else if strings.Index(m, "december") != -1 {
				month = 12
			}

			writeData(rawdata, month, year)
		}
		fmt.Println("Sleeping...")
		time.Sleep(60 * time.Second)
		fmt.Println("Awake!")
	}

}
