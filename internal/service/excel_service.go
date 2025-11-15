// Package service 包含了应用的业务逻辑
package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"questflow/internal/model"

	"github.com/xuri/excelize/v2"
)

// ExcelService 定义了导出 Excel 服务的接口
type ExcelService interface {
	ExportSubmissionsToExcel(formDef *formDefinition, submissions []model.Submission) (*bytes.Buffer, error)
}

// excelServiceImpl 是 ExcelService 的实现
type excelServiceImpl struct{}

// NewExcelService 创建一个新的 ExcelService 实例
func NewExcelService() ExcelService {
	return &excelServiceImpl{}
}

// ExportSubmissionsToExcel 将提交数据导出为 Excel 文件流
func (s *excelServiceImpl) ExportSubmissionsToExcel(formDef *formDefinition, submissions []model.Submission) (*bytes.Buffer, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "提交数据"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return nil, err
	}

	// --- 1. 构建并写入表头 ---
	headers := []string{"提交序号", "提交时间"}
	// questionID 到 Excel 列字母的映射，例如: {"q1": "C", "q2": "D"}
	questionIDToCol := make(map[string]string)

	for i, q := range formDef.Questions {
		headers = append(headers, q.Title)
		// Excel 列从 A 开始，前两列是固定的，所以题目从第3列（即索引2）开始
		colName, _ := excelize.ColumnNumberToName(i + 3)
		questionIDToCol[q.ID] = colName
	}

	// 写入表头
	if err := f.SetSheetRow(sheetName, "A1", &headers); err != nil {
		return nil, err
	}

	// --- 2. 遍历提交数据并写入每一行 ---
	for i, sub := range submissions {
		rowNum := i + 2 // 数据从第二行开始
		// a. 写入固定列：提交序号和提交时间
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", rowNum), i+1)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", rowNum), sub.CreatedAt.Format("2006-01-02 15:04:05"))

		// b. 解析答案并写入对应的题目列
		var answers map[string]interface{}
		// 在 GORM 中，datatypes.JSON 实际上是 []byte 类型
		if err := json.Unmarshal(sub.Data, &answers); err != nil {
			// 如果解析失败，跳过此条记录的答案部分
			continue
		}

		for qID, ans := range answers {
			colName, ok := questionIDToCol[qID]
			if !ok {
				continue // 如果答案中的 qID 在表单定义中找不到，则跳过
			}

			cell := fmt.Sprintf("%s%d", colName, rowNum)

			// c. 根据答案类型进行格式化
			// 在这里，我们需要将选项ID转换为可读的文本
			answerText := formatAnswer(ans, qID, formDef)
			f.SetCellValue(sheetName, cell, answerText)
		}
	}

	// --- 3. 调整样式（可选，但能提升体验）---
	// a. 设置表头样式
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"#E0E0E0"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	f.SetCellStyle(sheetName, "A1", fmt.Sprintf("%s1", headers[len(headers)-1]), headerStyle)

	// b. 自动调整列宽
	cols, _ := f.GetCols(sheetName)
	for i, col := range cols {
		colName, _ := excelize.ColumnNumberToName(i + 1)
		maxwidth := 0
		for _, cellValue := range col {
			cellWidth := len(cellValue)
			if cellWidth > maxwidth {
				maxwidth = cellWidth
			}
		}
		// 加一点余量
		f.SetColWidth(sheetName, colName, colName, float64(maxwidth+5))
	}
	// 特别设置提交时间列的宽度
	f.SetColWidth(sheetName, "B", "B", 20)

	f.SetActiveSheet(index)
	// 删除默认创建的 "Sheet1"
	f.DeleteSheet("Sheet1")

	// 将文件写入内存缓冲区
	buffer, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

// formatAnswer 是一个辅助函数，用于将不同类型的答案格式化为可读的字符串
func formatAnswer(ans interface{}, qID string, formDef *formDefinition) string {
	var question *questionDefinition
	for i := range formDef.Questions {
		if formDef.Questions[i].ID == qID {
			question = &formDef.Questions[i]
			break
		}
	}

	if question == nil {
		return "未知问题"
	}

	// 创建选项ID到文本的映射以便快速查找
	optionMap := make(map[string]string)
	for _, opt := range question.Options {
		optionMap[opt.ID] = opt.Text
	}

	switch question.Type {
	case "single_choice", "judgment":
		if optID, ok := ans.(string); ok {
			if text, exists := optionMap[optID]; exists {
				return text
			}
			return optID // 回退显示ID
		}
	case "multi_choice":
		if optIDs, ok := ans.([]interface{}); ok {
			var result string
			for i, optIDRaw := range optIDs {
				if optID, ok := optIDRaw.(string); ok {
					if text, exists := optionMap[optID]; exists {
						if i > 0 {
							result += ", "
						}
						result += text
					}
				}
			}
			return result
		}
	case "text_input":
		if text, ok := ans.(string); ok {
			return text
		}
	}

	// 如果所有类型都不匹配，则返回原始值的字符串表示
	return fmt.Sprintf("%v", ans)
}
