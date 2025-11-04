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
	f.DeleteSheet("Sheet1")

	f.SetColWidth(sheetName, "A", "A", 30)
	f.SetColWidth(sheetName, "B", "B", 30)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Bold: true, Size: 14},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
	})

	f.SetCellValue(sheetName, "A1", "АКТ ВЫПОЛНЕННЫХ РАБОТ")
	f.MergeCell(sheetName, "A1", "B1")
	f.SetCellStyle(sheetName, "A1", "B1", headerStyle)
	f.SetRowHeight(sheetName, 1, 30)

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
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), row[0])
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), row[1])
	}

	if err := f.SaveAs("templates/act_template.xlsx"); err != nil {
		log.Fatalf("Error saving file: %v", err)
	}

	log.Println("Template created successfully!")
}
