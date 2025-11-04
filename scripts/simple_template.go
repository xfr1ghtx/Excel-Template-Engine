package main

import (
	"fmt"
	"log"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalf("Error closing file: %v", err)
		}
	}()

	sheetName := "Акт"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		log.Fatalf("Error creating sheet: %v", err)
	}
	f.SetActiveSheet(index)
	err = f.DeleteSheet("Sheet1")
	if err != nil {
		log.Fatalf("Error deleting default sheet: %v", err)
	}

	err = f.SetColWidth(sheetName, "A", "A", 30)
	if err != nil {
		log.Fatalf("Error setting column width A: %v", err)
	}
	err = f.SetColWidth(sheetName, "B", "B", 30)
	if err != nil {
		log.Fatalf("Error setting column width B: %v", err)
	}

	var headerStyle int
	headerStyle, err = f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		log.Fatalf("Error creating header style: %v", err)
	}

	err = f.SetCellValue(sheetName, "A1", "АКТ ВЫПОЛНЕННЫХ РАБОТ")
	if err != nil {
		log.Fatalf("Error setting header title: %v", err)
	}
	err = f.MergeCell(sheetName, "A1", "B1")
	if err != nil {
		log.Fatalf("Error merging header cells: %v", err)
	}
	err = f.SetCellStyle(sheetName, "A1", "B1", headerStyle)
	if err != nil {
		log.Fatalf("Error setting header style: %v", err)
	}
	err = f.SetRowHeight(sheetName, 1, 30)
	if err != nil {
		log.Fatalf("Error setting header row height: %v", err)
	}

	data := [][]string{
		{"", ""},
		{"Номер договора:", "{{contractNumber}}"},
		{"Дата договора:", "{{contractDate}}"},
		{"Заказчик:", "{{customer}}"},
		{"Подрядчик:", "{{contractor}}"},
		{"Объект:", "{{objectName}}"},
		{"", ""},
		{"Общая стоимость:", "{{totalCost}}"},
		{"Стоимость инспекции:", "{{totalCostInspection}}"},
		{"Стоимость рассмотрения:", "{{totalCostConsiderations}}"},
		{"ID позиций:", "{{positionIds}}"},
		{"", ""},
		{"Дата создания:", "{{createdAt}}"},
	}

	for i, row := range data {
		rowNum := i + 2
		err = f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row[0])
		if err != nil {
			log.Fatalf("Error setting cell A%d: %v", rowNum, err)
		}
		err = f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row[1])
		if err != nil {
			log.Fatalf("Error setting cell B%d: %v", rowNum, err)
		}
	}

	err = f.SaveAs("templates/act_template.xlsx")
	if err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	log.Println("Template created successfully!")
}
