package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const qtdMonitoramentos = 3
const delayInSeconds = 5

func showMenu() {
	fmt.Println("1 - Start Monitoring")
	fmt.Println("2 - Show Logs")
	fmt.Println("0 - Exit")
}

func readCommand() int {
	var command int
	fmt.Scan(&command)
	return command
}

func registerLog(site string, status bool) {
	file, err := os.OpenFile("log.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println("Error when opening log file:", err)
	}

	file.WriteString(site + " - " + time.Now().Format("02/01/2006 15:04:05") + " - online: " + strconv.FormatBool(status) + "\n")

	file.Close()
}

func printLogs() {
	file, err := ioutil.ReadFile("log.txt")

	if err != nil {
		fmt.Println("Error when opening log file:", err)
	}

	fmt.Println(string(file))
}

func testSite(site string) {
	res, err := http.Get(site)
	if err != nil {
		fmt.Println("Error when checking website: ", err)
	} else if res.StatusCode == 200 {
		fmt.Println("Site:", site, "was loaded with success!")
		registerLog(site, true)
	} else {
		fmt.Println("Site:", site, "was loaded with error:", res.StatusCode)
		registerLog(site, false)
	}
}

func startMonitoring() {
	fmt.Println("Monitoring...")
	sites := readSitesFromFile()

	for i := 0; i < qtdMonitoramentos; i++ {
		for _, site := range sites {
			testSite(site)
		}
		fmt.Println("Waiting", delayInSeconds, "seconds...")
		time.Sleep(delayInSeconds * time.Second)
	}
}

func readSitesFromFile() []string {
	var sites []string

	file, err := os.Open("sites.txt")

	if err != nil {
		fmt.Println("Error when opening file:", err)
		os.Exit(-1)
	}

	reader := bufio.NewReader(file)

	for {
		line, err := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error when reading file:", err)
			break
		}
		fmt.Println(line)
		sites = append(sites, line)
	}

	defer file.Close()

	return sites
}

func main() {
	for {
		showMenu()
		command := readCommand()

		switch command {
		case 1:
			startMonitoring()
		case 2:
			fmt.Println("Showing Logs...")
			printLogs()
		case 0:
			fmt.Println("Exiting...")
			os.Exit(0)
		default:
			fmt.Println("Invalid command")
			os.Exit(-1)
		}
	}
}
