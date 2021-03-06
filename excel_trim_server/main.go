package main

import (
	"fmt"
	"github.com/tealeg/xlsx"
	"gopkg.in/ini.v1"
	"io/ioutil"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

type XLSXSheetInfo struct {
	SheetName string
	Table     [][]string
}

func NewXLSXSheetInfo(sheetName string) *XLSXSheetInfo {
	return &XLSXSheetInfo{
		SheetName: sheetName,
	}
}

type XLSXInfo struct {
	XLSXName string
	Sheets   []*XLSXSheetInfo
}

func NewXLSXInfo(XLSXName string) *XLSXInfo {
	return &XLSXInfo{XLSXName: XLSXName}
}

func ReadOneXLSX(dir string, excelFileName string, outPutDir string) *XLSXInfo {
	retInfo := NewXLSXInfo(excelFileName)
	//fmt.Println("ReadOneXLSX dir:", dir, " excelFileName:", excelFileName, " outPutDir:", outPutDir)
	fmt.Println("开始读取:", excelFileName)
	xlFile, err := xlsx.OpenFile(dir + excelFileName)
	if err != nil {
		fmt.Printf("open failed: %s\n", err)
		return nil
	}
	for _, sheet := range xlFile.Sheets {
		sheetInfo := NewXLSXSheetInfo(sheet.Name)
		//fmt.Printf("Sheet Name: %s\n", sheet.Name)
		for _, row := range sheet.Rows {
			var rowInfo []string
			for _, cell := range row.Cells {
				text := cell.String()
				rowInfo = append(rowInfo, text)
				//fmt.Printf("%s\t", text)
			}
			sheetInfo.Table = append(sheetInfo.Table, rowInfo)
			//fmt.Println()
		}
		retInfo.Sheets = append(retInfo.Sheets, sheetInfo)
		//fmt.Println()
	}

	return retInfo
}

func WriteOneXLSX(dir string, outPutDir string, info *XLSXInfo) {
	//fmt.Println("WriteOneXLSX dir:", dir, " outPutDir:", outPutDir)
	fmt.Println("开始写入:", info.XLSXName)
	fromXLSXName := info.XLSXName
	pos := strings.Index(fromXLSXName, "DT_")
	if pos == -1 {
		fmt.Println("WriteOneXLSX DT_ not found. fromXLSXName:", fromXLSXName)
		return
	}

	XLSXName := fromXLSXName[pos:]

	sheetName := XLSXName

	sheetName = strings.TrimSuffix(sheetName, ".xlsx")

	DTConfigsTrimdPrefixPostfix = append(DTConfigsTrimdPrefixPostfix, sheetName)

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	var err error

	file = xlsx.NewFile()

	for _, sheetInfo := range info.Sheets {
		sheet, err = file.AddSheet(sheetName)
		if err != nil {
			fmt.Printf(err.Error())
		}

		for _, rowInfo := range sheetInfo.Table {
			row = sheet.AddRow()
			for _, colInfo := range rowInfo {
				cell = row.AddCell()
				cell.Value = colInfo
			}
		}
	}

	err = file.Save(outPutDir + XLSXName)
	if err != nil {
		fmt.Printf(err.Error())
	}
}

var indent = 0
func WriteGoLine(file *os.File, str string){
	for i := 0; i < indent; i++ {
		file.WriteString("\t")
	}
	file.WriteString(str + "\n")
}

func WriteGoFiles(DTConfigs []string){
	iniFileName := "uniqs.ini"
	cfg, err := ini.Load(iniFileName)
	if err != nil {
		log.Fatal("ini.Load failed iniFileName:", iniFileName)
	}
	outPutDir := cfg.Section("").Key("dtconfig_go_output_dir").String()
	if !strings.HasSuffix(outPutDir, "\\") && !strings.HasSuffix(outPutDir, "/") {
		outPutDir += "/"
	}
	fmt.Println("dtconfig_go_output_dir:", outPutDir)
	err = os.MkdirAll(outPutDir, os.ModePerm)
	if err != nil {
		log.Error("os.MkdirAll failed. err:", err)
	}

	fileName := ""
	var file *os.File

	// export.go
	if 1 == 1 {
		indent = 0
		fileName = "export.go"
		fileName = outPutDir + fileName
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Error("os.Open failed. fileName:", fileName)
			return
		}
		defer file.Close()

		WriteGoLine(file, "package DataTables")
		WriteGoLine(file, "")
		for _, c := range DTConfigs {
			cfg := c
			cfgData := cfg + "_Data"
			WriteGoLine(file, "func Get" + cfg + "() *" + cfgData + " {")
			indent++
			WriteGoLine(file, "return inst.Get(\"" + cfg + "\").(*" + cfgData + ")")
			indent--
			WriteGoLine(file, "}")
			WriteGoLine(file, "")
		}
	}

	// init_generated.go
	if 1 == 1 {
		indent = 0
		fileName = "init_generated.go"
		fileName = outPutDir + fileName
		file, err = os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
		if err != nil {
			log.Error("os.Open failed. fileName:", fileName)
			return
		}
		defer file.Close()

		WriteGoLine(file, "package DataTables")
		WriteGoLine(file, "")
		WriteGoLine(file, "func InitGenerated() {")
		indent++
		for _, c := range DTConfigs {
			cfg := c
			cfgData := cfg + "_Data"
			WriteGoLine(file, "inst.Register(\"" + cfg + "\", &" + cfgData + "{})")
		}
		indent--
		WriteGoLine(file, "}")
		WriteGoLine(file, "")
	}
}

var DTConfigsTrimdPrefixPostfix []string

func main() {
	iniFileName := "uniqs.ini"
	cfg, err := ini.Load(iniFileName)
	if err != nil {
		log.Fatal("ini.Load failed iniFileName:", iniFileName)
	}
	xlsxDir := cfg.Section("").Key("xlsx_dir").String()
	if !strings.HasSuffix(xlsxDir, "\\") && !strings.HasSuffix(xlsxDir, "/") {
		xlsxDir += "/"
	}
	fmt.Println("xlsx_dir:", xlsxDir)
	outPutDir := cfg.Section("").Key("output_dir").String()
	if !strings.HasSuffix(outPutDir, "\\") && !strings.HasSuffix(outPutDir, "/") {
		outPutDir += "/"
	}
	fmt.Println("output_dir:", outPutDir)

	files, err := ioutil.ReadDir(xlsxDir)
	if err != nil {
		log.Fatal("ioutil.ReadDir failed", err)
	}
	err = os.MkdirAll(outPutDir, os.ModePerm)
	if err != nil {
		log.Fatal("os.MkdirAll failed. err:", err)
	}

	for _, f := range files {
		fName := f.Name()
		//fmt.Println(fName)
		if f.IsDir() {
			continue
		}
		if strings.HasPrefix(fName, "~$") {
			continue
		}
		if !strings.HasSuffix(fName, ".xlsx") {
			continue
		}

		XLSXInfo := ReadOneXLSX(xlsxDir, fName, outPutDir)
		if XLSXInfo != nil {
			WriteOneXLSX(xlsxDir, outPutDir, XLSXInfo)
		}
	}

	WriteGoFiles(DTConfigsTrimdPrefixPostfix)
}
