package utils

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"strings"
	"time"

	"stacks-base/backends/go-net-http/internal/repositories"
)

func BuildUsersCSV(users []repositories.User) ([]byte, error) {
	buffer := &bytes.Buffer{}
	writer := csv.NewWriter(buffer)
	writer.Comma = ';'
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
		pageWidth      = 612
		pageHeight     = 792
		leftMargin     = 36
		topMargin      = 756
		tableTop       = 692
		rowHeight      = 22
		headerFontSize = 8
		bodyFontSize   = 8
		titleFontSize  = 16
		metaFontSize   = 10
	)

	columns := []pdfTableColumn{
		{Title: "Nome", Width: 98},
		{Title: "Email", Width: 160},
		{Title: "Papel", Width: 52},
		{Title: "Status", Width: 52},
		{Title: "Excluido", Width: 78},
		{Title: "Ultimo login", Width: 100},
	}

	rows := make([][]string, 0, len(users))
	for _, user := range users {
		rows = append(rows, []string{
			user.Name,
			user.Email,
			user.Role,
			user.Status,
			formatNullableReportTime(user.DeletedAt),
			formatNullableReportTime(user.LastLoginAt),
		})
	}

	rowsPerPage := 24
	if len(rows) == 0 {
		rowsPerPage = 1
	}

	pageCount := 1
	if len(rows) > 0 {
		pageCount = (len(rows) + rowsPerPage - 1) / rowsPerPage
	}

	pageStreams := make([]string, 0, pageCount)
	for start := 0; start < len(rows) || (len(rows) == 0 && start == 0); start += rowsPerPage {
		end := start + rowsPerPage
		if end > len(rows) {
			end = len(rows)
		}

		pageRows := rows[start:end]
		pageStreams = append(pageStreams, buildUsersPDFPage(
			columns,
			pageRows,
			pdfPageMetadata{
				LeftMargin:     leftMargin,
				TopMargin:      topMargin,
				TableTop:       tableTop,
				RowHeight:      rowHeight,
				HeaderFontSize: headerFontSize,
				BodyFontSize:   bodyFontSize,
				TitleFontSize:  titleFontSize,
				MetaFontSize:   metaFontSize,
				GeneratedAt:    time.Now(),
				TotalRows:      len(users),
				PageLabel:      fmt.Sprintf("Pagina %d de %d", (start/rowsPerPage)+1, pageCount),
			},
		))
		if len(rows) == 0 {
			break
		}
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

type pdfTableColumn struct {
	Title string
	Width int
}

type pdfPageMetadata struct {
	LeftMargin     int
	TopMargin      int
	TableTop       int
	RowHeight      int
	HeaderFontSize int
	BodyFontSize   int
	TitleFontSize  int
	MetaFontSize   int
	GeneratedAt    time.Time
	TotalRows      int
	PageLabel      string
}

func buildUsersPDFPage(columns []pdfTableColumn, rows [][]string, meta pdfPageMetadata) string {
	var stream strings.Builder

	stream.WriteString("0.95 0.96 0.98 rg\n")
	stream.WriteString("0.80 0.84 0.90 RG\n")
	stream.WriteString(fmt.Sprintf("%d %d %d %d re B\n", meta.LeftMargin, meta.TableTop, totalPDFTableWidth(columns), meta.RowHeight))
	stream.WriteString("0.35 0.40 0.48 RG\n")

	stream.WriteString("BT\n")
	stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", meta.TitleFontSize))
	stream.WriteString("0.06 0.09 0.16 rg\n")
	stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", meta.LeftMargin, meta.TopMargin))
	stream.WriteString("(Relatorio de usuarios) Tj\n")

	stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", meta.MetaFontSize))
	stream.WriteString("0.30 0.36 0.45 rg\n")
	stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", meta.LeftMargin, meta.TopMargin-18))
	stream.WriteString(fmt.Sprintf("(Gerado em: %s) Tj\n", escapePDFText(meta.GeneratedAt.Format("2006-01-02 15:04:05"))))
	stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", meta.LeftMargin, meta.TopMargin-34))
	stream.WriteString(fmt.Sprintf("(Total de registros: %d) Tj\n", meta.TotalRows))
	stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", meta.LeftMargin+410, meta.TopMargin-34))
	stream.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(meta.PageLabel)))
	stream.WriteString("ET\n")

	writePDFTableHeader(&stream, columns, meta)

	if len(rows) == 0 {
		writePDFEmptyState(&stream, columns, meta)
		return stream.String()
	}

	y := meta.TableTop - meta.RowHeight
	for _, row := range rows {
		writePDFTableRow(&stream, columns, row, y, meta)
		y -= meta.RowHeight
	}

	return stream.String()
}

func writePDFTableHeader(stream *strings.Builder, columns []pdfTableColumn, meta pdfPageMetadata) {
	x := meta.LeftMargin
	stream.WriteString("BT\n")
	stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", meta.HeaderFontSize))
	stream.WriteString("0.18 0.23 0.30 rg\n")

	for _, column := range columns {
		stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", x+6, meta.TableTop+7))
		stream.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(column.Title)))
		x += column.Width
	}

	stream.WriteString("ET\n")
}

func writePDFTableRow(stream *strings.Builder, columns []pdfTableColumn, row []string, y int, meta pdfPageMetadata) {
	x := meta.LeftMargin
	stream.WriteString("0.88 0.90 0.94 RG\n")
	for _, column := range columns {
		stream.WriteString(fmt.Sprintf("%d %d %d %d re S\n", x, y, column.Width, meta.RowHeight))
		x += column.Width
	}

	x = meta.LeftMargin
	stream.WriteString("BT\n")
	stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", meta.BodyFontSize))
	stream.WriteString("0.06 0.09 0.16 rg\n")

	for index, column := range columns {
		value := ""
		if index < len(row) {
			value = truncatePDFCell(row[index], column.Width, meta.BodyFontSize)
		}
		stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", x+6, y+8))
		stream.WriteString(fmt.Sprintf("(%s) Tj\n", escapePDFText(value)))
		x += column.Width
	}

	stream.WriteString("ET\n")
}

func writePDFEmptyState(stream *strings.Builder, columns []pdfTableColumn, meta pdfPageMetadata) {
	stream.WriteString("0.88 0.90 0.94 RG\n")
	stream.WriteString(fmt.Sprintf("%d %d %d %d re S\n", meta.LeftMargin, meta.TableTop-meta.RowHeight, totalPDFTableWidth(columns), meta.RowHeight))
	stream.WriteString("BT\n")
	stream.WriteString(fmt.Sprintf("/F1 %d Tf\n", meta.BodyFontSize))
	stream.WriteString("0.30 0.36 0.45 rg\n")
	stream.WriteString(fmt.Sprintf("1 0 0 1 %d %d Tm\n", meta.LeftMargin+6, meta.TableTop-meta.RowHeight+8))
	stream.WriteString("(Nenhum usuario encontrado.) Tj\n")
	stream.WriteString("ET\n")
}

func totalPDFTableWidth(columns []pdfTableColumn) int {
	total := 0
	for _, column := range columns {
		total += column.Width
	}
	return total
}

func truncatePDFCell(value string, width int, fontSize int) string {
	maxChars := width / max(4, (fontSize/2)+1)
	if maxChars < 6 {
		maxChars = 6
	}
	return truncatePDFLine(value, maxChars)
}

func formatNullableReportTime(value *time.Time) string {
	if value == nil {
		return "-"
	}
	return value.Format("2006-01-02 15:04")
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
