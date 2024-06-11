package plan

import (
	"fmt"
	"io"
	"mpt_data/api/apihelper"
	"mpt_data/api/middleware"
	"mpt_data/database/plan"
	generalmodel "mpt_data/models/general"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const packageName = "api.plan"

func getPlanPDF(w http.ResponseWriter, r *http.Request, startDate time.Time, endDate time.Time) {
	const funcName = packageName + ".getPlanPDF"

	tx := middleware.GetTx(r.Context())
	path, err := plan.GetOrCreatePDF(tx, generalmodel.Period{StartDate: startDate, EndDate: endDate})
	if err != nil {
		apihelper.InternalError(w, err)
		return
	}

	pdfFile, err := os.Open(path)
	if err != nil {
		apihelper.InternalError(w, err)
		return
	}
	defer pdfFile.Close()

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(path)))

	// Kopiere den Inhalt der PDF-Datei in die Antwort
	_, err = io.Copy(w, pdfFile)
	if err != nil {
		apihelper.InternalError(w, err)
		return
	}
}
