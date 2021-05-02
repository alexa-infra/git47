package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"net/http"
)

type summaryViewData struct {
	core.SummaryData
	*RequestContext
}

// GitSummary returns handler which renders summary page of a repository
func GitSummary(w http.ResponseWriter, r *http.Request) {
	ctx, _ := GetRequestContext(r)

	summary, err := core.GetSummary(ctx.Ref)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := summaryViewData{
		SummaryData:    summary,
		RequestContext: ctx,
	}

	err = RenderTemplate(w, "git-summary.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
