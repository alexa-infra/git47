package web

import (
	"github.com/alexa-infra/git47/internal/core"
	"net/http"
)

type summaryViewData struct {
	core.SummaryData
	RequestContext *requestContext
}

func summaryHandler(w http.ResponseWriter, r *http.Request) {
	ctx, _ := getRequestContext(r)

	summary, err := core.GetSummary(ctx.Ref)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := summaryViewData{
		SummaryData:    summary,
		RequestContext: ctx,
	}

	err = renderTemplate(w, "git-summary.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
