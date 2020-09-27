package scan

import (
	"fmt"
	"math"
	"net/http"
	"time"

	"okapi/lib/ores"
	"okapi/lib/task"
	"okapi/lib/wiki"
	"okapi/models"
)

// Worker processor for one title and store him into database as page
func Worker(id int, payload task.Payload) (string, map[string]interface{}, error) {
	message := "page title: '%s', page id: #%d"
	page := payload.(*models.Page)
	info := map[string]interface{}{
		"_title":  page.Title,
		"_status": http.StatusOK,
		"_id":     page.ID,
	}

	wikiPage, status, err := wiki.Client(page.SiteURL).GetMeta(page.Title)

	if err != nil {
		return "", info, err
	}

	if status != http.StatusOK {
		info["_status"] = status
		return "", info, fmt.Errorf(message+", status code: %d", page.Title, page.ID, status)
	}

	if wikiPage.Redirect {
		info["_status"] = 301
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, "page is a redirect!")
	}

	models.DB().
		Model(page).
		Where("title = ? and project_id = ?", page.Title, page.ProjectID).
		Select()

	page.Title = wikiPage.Title
	page.NsID = wikiPage.Namespace
	current := time.Now()
	diff := current.Sub(wikiPage.Timestamp)

	if int(math.Floor(diff.Hours())) < page.Project.TimeDelay {
		err = scoreRevision(page, wikiPage.Revision, current)
	} else {
		page.SetRevision(wikiPage.Revision)
	}

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	err = models.Save(page)

	if err != nil {
		return "", info, fmt.Errorf(message+", %s", page.Title, page.ID, err)
	}

	info["_id"] = page.ID

	return fmt.Sprintf(message, page.Title, page.ID), info, nil
}

func scoreRevision(page *models.Page, revision int, currentTime time.Time) error {
	threshold := page.Project.GetThreshold(ores.Damaging)

	if threshold == nil {
		return fmt.Errorf("%s threshold model does not exist", ores.Damaging)
	}

	score, scoreErr := ores.Damaging.ScoreOne(page.Project.DBName, revision)

	if scoreErr == nil && score.Probability.False >= *threshold {
		page.SetRevision(page.Revision)
		page.SetScore(page.Revision, ores.Damaging, *score)
		return nil
	}

	revisions, status, err := wiki.Client(page.SiteURL).GetRevisionsHistory(page.Title, 10)

	if err != nil {
		return err
	}

	if status != http.StatusOK {
		return fmt.Errorf("response failed with status code -> '%d'", status)
	}

	revs := []int{}

	for _, rev := range revisions {
		revs = append(revs, rev.RevID)
	}

	if scoreErr == nil {
		scores, err := ores.ScoreMany(page.Project.DBName, ores.Damaging, revs)

		if err == nil {
			for _, rev := range revs {
				if score, ok := scores[rev]; ok {
					if score.Probability.False >= *threshold {
						page.SetRevision(rev)
						page.SetScore(rev, ores.Damaging, *score)
						return nil
					}
				}
			}
		}
	}

	for _, rev := range revisions {
		diff := currentTime.Sub(rev.Timestamp)

		if int(math.Floor(diff.Hours())) >= page.Project.TimeDelay {
			page.SetRevision(rev.RevID)
			return nil
		}
	}

	return fmt.Errorf("all scoring methods failed")
}
