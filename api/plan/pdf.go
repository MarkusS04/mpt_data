package plan

import (
	"fmt"
	"io"
	api_helper "mpt_data/api/apihelper"
	"mpt_data/database/plan"
	generalmodel "mpt_data/models/general"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const packageName = "api.plan"

func getPlanPDF(w http.ResponseWriter, _ *http.Request, startDate time.Time, endDate time.Time) {
	const funcName = packageName + ".getPlanPDF"

	path, err := plan.GetOrCreatePDF(generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}

	pdfFile, err := os.Open(path)
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}
	defer pdfFile.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))

	// Kopiere den Inhalt der PDF-Datei in die Antwort
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		api_helper.InternalError(w, funcName, err.Error())
		return
	}
}
