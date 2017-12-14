package main

import (
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	"os"
	"encoding/json"
)

type printData struct {
	serialNumber string
	weight string
	ipAddress string

}

type Settings struct {
	FinalPrintCodeFolder   string
	InterimPrintCodeFolder string

}

var settings Settings

const homeFolder = `C:\LabelPrint`
const settingsFile = homeFolder + `\settings.json`

func main() {

	_, err := os.Stat(homeFolder)
	if err != nil {
		fmt.Println("Creating home folder.")
		os.Mkdir(homeFolder, 0777)
	}

	args := os.Args

	if len(args) > 1 {
		switch args[1] {
		case "help":
			fmt.Println("enter help info")
		case "config":
			getConfig()
		case "-v":
			fmt.Println("enter version info")
		default:
			fmt.Println("enter default info")
		}
	}

	_, err = os.Stat(settingsFile)
	if err != nil {
		fmt.Println(`Could not locate settings.json.
Run setup using argument config. ex "LabelPrintcontrol config"`)
		os.Exit(1)
	}

	settingsData, err := ioutil.ReadFile(settingsFile)

	json.Unmarshal(settingsData, &settings)

	pData := printData{"QR121231234", "212.76", "10.50.201.31"}
	printLabel(pData)
}

func getConfig(){

	settings.FinalPrintCodeFolder = getUserInput("Enter final label print code folder: ")
	settings.InterimPrintCodeFolder = getUserInput("Enter interim label print code folder: ")

	_, err := os.Stat(settings.FinalPrintCodeFolder)
	if err != nil {
		fmt.Println("Creating final print folder.")
		os.MkdirAll(settings.FinalPrintCodeFolder, 0777)
	}

	_, err = os.Stat(settings.InterimPrintCodeFolder)
	if err != nil {
		fmt.Println("creating interim label print folder.")
		os.MkdirAll(settings.InterimPrintCodeFolder, 0777)
	}

	settingData, err := json.Marshal(&settings)
	if err != nil {
		fmt.Println("Could not create settings file.")
	}

	err = ioutil.WriteFile(settingsFile, settingData, 0777)
	if err != nil {
		fmt.Println("could not save settings.json")
	}
}

func getUserInput(prompt string) string{
	fmt.Println("")
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)

	return input
}

func printLabel(pd printData){

	printCode := getPrintCode(pd.serialNumber[:2])
	replaceVariables(&printCode, pd.serialNumber, pd.weight)

	fmt.Println(printCode)
}

func getPrintCode(alphaCode string) string {

	fileName := settings.FinalPrintCodeFolder + `\` + alphaCode + ".txt"

	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
	}

	return string(data)

}

func replaceVariables(pc *string, sn, w string){

	m := time.Now().Month()
	month := fmt.Sprintf("%02d", m)
	y := time.Now().Year()
	year := fmt.Sprintf("%d", y)

	*pc = strings.Replace(*pc, "_SerialNumber_", sn, -1)
	*pc = strings.Replace(*pc, "_Weight_", w, -1)
	*pc = strings.Replace(*pc, "_Month_", month, -1)
	*pc = strings.Replace(*pc, "_Year_", year, -1)
}
