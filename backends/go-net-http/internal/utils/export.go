package utils

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"strconv"
	"strings"
	"time"

	"stacks-base/backends/go-net-http/internal/repositories"
)

func BuildUsersCSV(users []repositories.User) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	rows := buildUserRows(users)
	if err := writer.WriteAll(rows); err != nil {
		return nil, fmt.Errorf("write csv: %w", err)
	}
	return buffer.Bytes(), nil
}

func BuildUsersXLSX(users []repositories.User) ([]byte, error) {
	rows := buildUserRows(users)

	buffer := &bytes.Buffer{}
	archive := zip.NewWriter(buffer)

	files := map[string]string{
		"[Content_Types].xml":        contentTypesXML,
		"_rels/.rels":                rootRelsXML,
		"xl/workbook.xml":            workbookXML,
		"xl/_rels/workbook.xml.rels": workbookRelsXML,
		"xl/styles.xml":              stylesXML,
		"xl/worksheets/sheet1.xml":   buildWorksheetXML(rows),
	}

	for name, content := range files {
		writer, err := archive.Create(name)
		if err != nil {
			return nil, fmt.Errorf("create xlsx entry %s: %w", name, err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			return nil, fmt.Errorf("write xlsx entry %s: %w", name, err)
		}
	}

	if err := archive.Close(); err != nil {
		return nil, fmt.Errorf("close xlsx archive: %w", err)
	}

	return buffer.Bytes(), nil
}

func BuildUsersPDF(users []repositories.User) ([]byte, error) {
	const (
		pageWidth     = 612
		pageHeight    = 792
		leftMargin    = 40
		topMargin     = 760
		lineHeight    = 14
		rowsPerPage   = 42
		bodyFontSize  = 9
		titleFontSize = 16
	)

	lines := []string{
		"Relatorio de usuarios",
		"Gerado em: " + time.Now().Format("2006-01-02 15:04:05"),
		"Total de registros: " + strconv.Itoa(len(users)),
		"",
		"Nome | Email | Papel | Status | Excluido | Ultimo login",
		strings.Repeat("-", 110),
	}

	for _, user := range users {
		lines = append(lines, truncatePDFLine(fmt.Sprintf(
			"%s | %s | %s | %s | %s | %s",
			user.Name,
			user.Email,
			user.Role,
			user.Status,
			formatNullableTime(user.DeletedAt),
			formatNullableTime(user.LastLoginAt),
		), 120))
	}

	if len(lines) == 0 {
		lines = append(lines, "Nenhum usuario encontrado.")
	}

	pageStreams := make([]string, 0, (len(lines)/rowsPerPage)+1)
	for start := 0; start < len(lines); start += rowsPerPage {
		end := start + rowsPerPage
		if end > len(lines) {
			end = len(lines)
		}

		var stream strings.Builder
		stream.WriteString("BT\n")
		stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", titleFontSize))
		stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", leftMargin, topMargin))
		stream.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(lines[start])))
		stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", bodyFontSize))

		y := topMargin - (lineHeight * 2)
		for _, line := range lines[start+1 : end] {
			stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", leftMargin, y))
			stream.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(line)))
			y -= lineHeight
		}
		stream.WriteString("ET\n")
		pageStreams = append(pageStreams, stream.String())
	}

	objects := []string{
		"<< /Type /Catalog /Pages 2 0 R >>",
		"",
		"<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>",
	}

	pageObjectNumbers := make([]int, 0, len(pageStreams))
	for index, stream := range pageStreams {
		pageObjectNumber := 4 + (index * 2)
		contentObjectNumber := pageObjectNumber + 1
		pageObjectNumbers = append(pageObjectNumbers, pageObjectNumber)

		objects = append(objects,
			fmt.Sprintf("<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %d %d] /Resources << /Font << /F1 3 0 R >> >> /Contents %d 0 R >>", pageWidth, pageHeight, contentObjectNumber),
			fmt.Sprintf("<< /Length %d >>\nstream\n%sendstream", len(stream), stream),
		)
	}

	kids := make([]string, 0, len(pageObjectNumbers))
	for _, objectNumber := range pageObjectNumbers {
		kids = append(kids, fmt.Sprintf("%d 0 R", objectNumber))
	}
	objects[1] = fmt.Sprintf("<< /Type /Pages /Count %d /Kids [%s] >>", len(pageObjectNumbers), strings.Join(kids, " "))

	var document bytes.Buffer
	document.WriteString("%PDF-1.4\n")

	offsets := make([]int, 0, len(objects)+1)
	offsets = append(offsets, 0)
	for index, object := range objects {
		offsets = append(offsets, document.Len())
		document.WriteString(fmt.Sprintf("%d 0 obj\n%s\nendobj\n", index+1, object))
	}

	xrefOffset := document.Len()
	document.WriteString(fmt.Sprintf("xref\n0 %d\n", len(objects)+1))
	document.WriteString("0000000000 65535 f \n")
	for _, offset := range offsets[1:] {
		document.WriteString(fmt.Sprintf("%010d 00000 n \n", offset))
	}
	document.WriteString(fmt.Sprintf("trailer << /Size %d /Root 1 0 R >>\n", len(objects)+1))
	document.WriteString(fmt.Sprintf("startxref\n%d\n%%%%EOF", xrefOffset))

	return document.Bytes(), nil
}

func buildUserRows(users []repositories.User) [][]string {
	rows := [][]string{
		{"ID", "Nome", "Email", "Papel", "Status", "Excluido em", "Ultimo login", "Criado em", "Atualizado em"},
	}

	for _, user := range users {
		rows = append(rows, []string{
			user.ID,
			user.Name,
			user.Email,
			user.Role,
			user.Status,
			formatNullableTime(user.DeletedAt),
			formatNullableTime(user.LastLoginAt),
			user.CreatedAt.Format(time.RFC3339),
			user.UpdatedAt.Format(time.RFC3339),
		})
	}

	return rows
}

func buildWorksheetXML(rows [][]string) string {
	var builder strings.Builder
	builder.WriteString(`<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`)
	builder.WriteString(`<worksheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main"><sheetData>`)

	for rowIndex, row := range rows {
		builder.WriteString(fmt.Sprintf(`<row r="%d">`, rowIndex+1))
		for columnIndex, value := range row {
			cellReference := excelCellReference(columnIndex + 1)
			builder.WriteString(fmt.Sprintf(`<c r="%s%d" t="inlineStr"><is><t>%s</t></is></c>`, cellReference, rowIndex+1, xmlEscape(value)))
		}
		builder.WriteString(`</row>`)
	}

	builder.WriteString(`</sheetData></worksheet>`)
	return builder.String()
}

func excelCellReference(index int) string {
	result := ""
	for index > 0 {
		index--
		result = string(rune('A'+(index%26))) + result
		index /= 26
	}
	return result
}

func xmlEscape(value string) string {
	var buffer bytes.Buffer
	_ = xml.EscapeText(&buffer, []byte(value))
	return buffer.String()
}

func formatNullableTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}

func escapePDFText(value string) string {
	replacer := strings.NewReplacer(
		"\\", "\\\\",
		"(", "\\(",
		")", "\\)",
	)
	return replacer.Replace(value)
}

func truncatePDFLine(value string, max int) string {
	if len(value) <= max {
		return value
	}
	if max <= 3 {
		return value[:max]
	}
	return value[:max-3] + "..."
}

const contentTypesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/xl/workbook.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet.main+xml"/>
  <Override PartName="/xl/worksheets/sheet1.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.worksheet+xml"/>
  <Override PartName="/xl/styles.xml" ContentType="application/vnd.openxmlformats-officedocument.spreadsheetml.styles+xml"/>
</Types>`

const rootRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`

const workbookXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<workbook xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main" xmlns:r="http://schemas.openxmlformats.org/officeDocument/2006/relationships">
  <sheets>
    <sheet name="Users" sheetId="1" r:id="rId1"/>
  </sheets>
</workbook>`

const workbookRelsXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/worksheet" Target="worksheets/sheet1.xml"/>
  <Relationship Id="rId2" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/styles" Target="styles.xml"/>
</Relationships>`

const stylesXML = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<styleSheet xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
  <fonts count="1">
    <font>
      <sz val="11"/>
      <name val="Calibri"/>
    </font>
  </fonts>
  <fills count="1">
    <fill>
      <patternFill patternType="none"/>
    </fill>
  </fills>
  <borders count="1">
    <border/>
  </borders>
  <cellStyleXfs count="1">
    <xf numFmtId="0" fontId="0" fillId="0" borderId="0"/>
  </cellStyleXfs>
  <cellXfs count="1">
    <xf numFmtId="0" fontId="0" fillId="0" borderId="0" xfId="0"/>
  </cellXfs>
</styleSheet>`
