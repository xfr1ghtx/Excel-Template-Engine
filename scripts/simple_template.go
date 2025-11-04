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
	if err := f.DeleteSheet("Sheet1"); err != nil {
		log.Fatalf("Error deleting default sheet: %v", err)
	}

	if err := f.SetColWidth(sheetName, "A", "A", 30); err != nil {
		log.Fatalf("Error setting column width A: %v", err)
	}
	if err := f.SetColWidth(sheetName, "B", "B", 30); err != nil {
		log.Fatalf("Error setting column width B: %v", err)
	}

	headerStyle, err := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})
	if err != nil {
		log.Fatalf("Error creating header style: %v", err)
	}

	if err := f.SetCellValue(sheetName, "A1", "АКТ ВЫПОЛНЕННЫХ РАБОТ"); err != nil {
		log.Fatalf("Error setting header title: %v", err)
	}
	if err := f.MergeCell(sheetName, "A1", "B1"); err != nil {
		log.Fatalf("Error merging header cells: %v", err)
	}
	if err := f.SetCellStyle(sheetName, "A1", "B1", headerStyle); err != nil {
		log.Fatalf("Error setting header style: %v", err)
	}
	if err := f.SetRowHeight(sheetName, 1, 30); err != nil {
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
		if err := f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row[0]); err != nil {
			log.Fatalf("Error setting cell A%d: %v", rowNum, err)
		}
		if err := f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row[1]); err != nil {
			log.Fatalf("Error setting cell B%d: %v", rowNum, err)
		}
	}

	if err := f.SaveAs("templates/act_template.xlsx"); err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	log.Println("Template created successfully!")
}
