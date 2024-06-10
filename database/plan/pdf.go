// Package plan provides functions to create, retrive, update plans and select availabe and absent people for the plan
package plan

import (
	"fmt"
	"mpt_data/database/logging"
	"mpt_data/database/task"
	"mpt_data/helper/config"
	dbModel "mpt_data/models/dbmodel"
	generalmodel "mpt_data/models/general"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	"gorm.io/gorm"
)

type (
	// pdf represents the PDF document.
	// holds the file itself and asociated parameters
	pdf struct {
		file               *fpdf.Fpdf
		rowHeight          float64
		WidthPageAvailable float64
		widthDate          float64
		colorTextHeader    rgb
		colorBackHeader    rgb
		// 0 for even row, 1 for odd row
		colorBack [2]rgb
	}

	// rgb represents the RGB color.
	rgb struct {
		r, g, b int
	}

	// pdfDate holds the data that should be printed to pdf
	pdfData struct {
		tasks []dbModel.Task
		data  []planData
	}

	planData struct {
		meeting dbModel.Meeting
		tag     dbModel.Tag
		person  map[*dbModel.TaskDetail]dbModel.Person
	}
)

// GetOrCreatePDF generates a PDF file based on the provided period.
func GetOrCreatePDF(db *gorm.DB, period generalmodel.Period) (path string, err error) {
	const funcName = packageName + ".GetOrCreatePDF"
	var file dbModel.PDF
	if err :=
		db.
			Where("start_date = ?", period.StartDate).
			Where("end_date = ?", period.EndDate).
			First(&file).Error; err != nil {
		logging.LogError(funcName, err.Error())
	} else if !file.DataChanged && file.FilePath != "" {
		return file.FilePath, nil
	}

	pdf := getPDF()

	headline := pdf.printDateTitle(period)

	pdfData, err := getPdfData(db, period)
	if err != nil {
		logging.LogError(funcName, err.Error())
		return "", err
	}
	pdf.printTable(pdfData)

	pdfName := fmt.Sprintf("Dienerplan-%s.pdf", strings.ReplaceAll(headline, " ", ""))
	pdfFile := fmt.Sprintf("%s/%s", config.Config.PDF.Path, pdfName)

	if err :=
		pdf.file.OutputFileAndClose(pdfFile); err != nil {
		logging.LogError(funcName, err.Error())
		return "", err
	}

	if file.ID == 0 {
		db.Create(&dbModel.PDF{Name: pdfName, FilePath: pdfFile, Period: period, DataChanged: false})
	} else {
		file.DataChanged = false
		file.FilePath = pdfFile
		file.Name = pdfName

		db.Save(&file)
	}

	return pdfFile, nil
}

// getPDF initializes the PDF document.
func getPDF() *pdf {
	pdf := &pdf{
		file: fpdf.New(fpdf.OrientationPortrait, fpdf.UnitCentimeter, fpdf.PageSizeA4, ""),
	}
	pdf.file.SetMargins(2, 2, 2)
	pdf.file.SetFooterFunc(func() {
		pdf.file.SetY(-1.5)
		pdf.file.SetFont("Times", "I", 10)
		pdf.file.CellFormat(0, 1.0, fmt.Sprintf("Stand %s", time.Now().Format("02.01.2006")), "", 0, "C", false, 0, "")
	})
	pdf.file.AddPage()

	pdf.rowHeight = 0.57
	pdf.WidthPageAvailable = 17
	pdf.widthDate = 4.75

	pdf.colorTextHeader = rgb{r: 255, g: 255, b: 255}
	pdf.colorBackHeader = rgb{r: 68, g: 113, b: 196}
	pdf.colorBack[0] = rgb{r: 217, g: 226, b: 243}
	pdf.colorBack[1] = rgb{r: 255, g: 255, b: 255}

	return pdf
}

// getPdfData retrieves data needed for PDF generation.
func getPdfData(db *gorm.DB, period generalmodel.Period) (data pdfData, err error) {
	if data.tasks, err = task.GetTask(db); err != nil {
		return data, err
	}

	planFields, err := GetPlan(period)
	if err != nil {
		return data, err
	}

	// If planFields are present, organize data accordingly
	if len(planFields) > 0 {
		lastMeeting := planFields[0].MeetingID
		currentPlanData := planData{
			meeting: planFields[0].Meeting,
			person:  make(map[*dbModel.TaskDetail]dbModel.Person),
			tag:     planFields[0].Meeting.Tag,
		}

		for _, plan := range planFields {
			// If new meeting, create new planData entry, meetings are ordered by data
			if plan.MeetingID != lastMeeting {
				data.data = append(data.data, currentPlanData)
				currentPlanData = planData{
					meeting: plan.Meeting,
					person:  make(map[*dbModel.TaskDetail]dbModel.Person),
					tag:     plan.Meeting.Tag,
				}
			}
			currentPlanData.person[&plan.TaskDetail] = plan.Person
			lastMeeting = plan.MeetingID
		}
		data.data = append(data.data, currentPlanData)
	}

	return data, nil
}

// printTable prints the table with data.
func (pdf *pdf) printTable(data pdfData) {
	pdf.file.SetFont("Times", "", 12)

	for _, task := range data.tasks {
		pdf.printTaskHeader(task)

		for i, row := range data.data {
			pdf.setTextColor(rgb{0, 0, 0})
			pdf.setFillColor(pdf.colorBack[i%2])

			width := (pdf.WidthPageAvailable - pdf.widthDate) / float64(len(task.TaskDetails))
			pdf.file.SetFont("Times", "B", 12)
			pdf.writeCell(pdf.widthDate/2, row.meeting.Date.Format("02.01."))
			pdf.file.SetFont("Times", "", 12)
			pdf.writeCell(pdf.widthDate/2, getWeekdayName(row.meeting.Date.Weekday(), German))

			if row.tag.ID != 0 {
				pdf.writeCell(pdf.WidthPageAvailable-pdf.widthDate, row.tag.Descr)
				pdf.file.Ln(-1)
				continue
			}

			for i := 0; i < len(task.TaskDetails); i++ {
				if person, ok := getEntryByAttributeValue(row.person, task.TaskDetails[i].ID); ok {
					pdf.writeCell(width, fmt.Sprintf("%s %s", person.GivenName, person.LastName))
				} else {
					pdf.writeCell(width, "")
				}
			}
			pdf.file.Ln(-1)
		}

		pdf.file.Ln(1)
	}
}

// printDateTitle prints the title with date.
func (pdf *pdf) printDateTitle(period generalmodel.Period) string {
	pdf.file.SetFont("Times", "B", 15)
	var headline string
	if period.StartDate.Month() == period.EndDate.Month() && period.StartDate.Year() == period.EndDate.Year() {
		headline = fmt.Sprintf("%s %s",
			getMonthName(period.StartDate.Month(), German),
			period.StartDate.Format(" 2006"),
		)
	} else if period.StartDate.Month() != period.EndDate.Month() && period.StartDate.Year() == period.EndDate.Year() {
		headline = fmt.Sprintf("%s-%s %s",
			getMonthName(period.StartDate.Month(), German),
			getMonthName(period.EndDate.Month(), German),
			period.StartDate.Format("2006"),
		)
	} else {
		headline = fmt.Sprintf("%s %s - %s %s",
			getMonthName(period.StartDate.Month(), German),
			period.StartDate.Format("2006"),
			getMonthName(period.EndDate.Month(), German),
			period.EndDate.Format("2006"),
		)
	}
	pdf.file.CellFormat(0, pdf.rowHeight, pdf.file.UnicodeTranslatorFromDescriptor("")(headline), "", 1, "C", false, 0, "")
	pdf.file.Ln(1)

	return headline
}

// setFillColor sets the fill color of the PDF.
func (pdf *pdf) setFillColor(rgb rgb) { pdf.file.SetFillColor(rgb.r, rgb.g, rgb.b) }

// setTextColor sets the text color of the PDF.
func (pdf *pdf) setTextColor(rgb rgb) { pdf.file.SetTextColor(rgb.r, rgb.g, rgb.b) }

/*
write text to table cell

	border is set
	text is centered
	background is painted
*/
func (pdf *pdf) writeCell(width float64, text string) {
	pdf.file.CellFormat(width, pdf.rowHeight, pdf.file.UnicodeTranslatorFromDescriptor("")(text), "1", 0, "C", true, 0, "")
}

// printTaskHeader prints the header for a task.
func (pdf *pdf) printTaskHeader(headline dbModel.Task) {
	pdf.setTextColor(pdf.colorTextHeader)
	pdf.setFillColor(pdf.colorBackHeader)

	pdf.file.CellFormat(pdf.WidthPageAvailable, pdf.rowHeight, pdf.file.UnicodeTranslatorFromDescriptor("")(headline.Descr), "1", 1, "C", true, 0, "")
	pdf.writeCell(pdf.widthDate, "Datum")

	width := (pdf.WidthPageAvailable - pdf.widthDate) / float64(len(headline.TaskDetails))
	for _, secondaryHeader := range headline.TaskDetails {
		pdf.file.CellFormat(width, pdf.rowHeight, pdf.file.UnicodeTranslatorFromDescriptor("")(secondaryHeader.Descr), "1", 0, "C", true, 0, "")
	}
	pdf.file.Ln(-1)
}

func getEntryByAttributeValue(taskDetailsMap map[*dbModel.TaskDetail]dbModel.Person, attributeValue uint) (dbModel.Person, bool) {
	for taskDetail, person := range taskDetailsMap {
		if taskDetail.ID == attributeValue {
			return person, true
		}
	}
	return dbModel.Person{}, false
}
