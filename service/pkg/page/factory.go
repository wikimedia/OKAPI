package page

import (
	"fmt"
	"okapi-data-service/models"
	schema "okapi-data-service/schema/v3"
	"strings"
	"time"

	"github.com/protsack-stephan/mediawiki-api-client"
)

const wikidataURL = "http://www.wikidata.org/entity/"

// Factory create schema.org article
type Factory struct {
	Project   *models.Project
	Language  *models.Language
	Namespace *models.Namespace
}

// Create change mediawiki page into schema.org article
func (f *Factory) Create(data *mediawiki.PageData, html string) *schema.Page {
	dateModified := time.Now()

	page := &schema.Page{
		Name:         data.Title,
		Identifier:   data.PageID,
		DateModified: &dateModified,
		Version: &schema.Version{
			Identifier: data.LastRevID,
		},
		URL: data.CanonicalURL,
		InLanguage: &schema.Language{
			Name:       f.Language.LocalName,
			Identifier: f.Language.Code,
		},
		IsPartOf: &schema.Project{
			Name:       f.Project.SiteName,
			Identifier: f.Project.DbName,
		},
		Namespace: &schema.Namespace{
			Identifier: f.Namespace.ID,
			Name:       f.Namespace.Title,
		},
		ArticleBody: &schema.ArticleBody{
			HTML: html,
		},
		License: []*schema.License{
			schema.NewLicense(),
		},
	}

	// Custom license for wikinews projects.
	if f.Project.SiteCode == "wikinews" {
		page.License = []*schema.License{
			{
				Name:       "Attribution 2.5 Generic",
				Identifier: "CC-BY-2.5",
				URL:        "https://creativecommons.org/licenses/by/2.5/",
			},
		}
	}

	if len(data.Revisions) > 0 {
		page.DateModified = &data.Revisions[0].Timestamp
		page.ArticleBody.Wikitext = data.Revisions[0].Slots.Main.Content
		page.Version = &schema.Version{
			Identifier:      data.LastRevID,
			Comment:         data.Revisions[0].Comment,
			Tags:            data.Revisions[0].Tags,
			IsMinorEdit:     data.Revisions[0].Minor,
			IsFlaggedStable: data.Flagged.StableRevID == data.LastRevID,
			Editor: &schema.Editor{
				Identifier:  data.Revisions[0].UserID,
				Name:        data.Revisions[0].User,
				IsAnonymous: data.Revisions[0].UserID == 0,
			},
		}
	}

	if len(data.Pageprops.WikibaseItem) != 0 {
		page.MainEntity = &schema.Entity{
			Identifier: data.Pageprops.WikibaseItem,
			URL:        fmt.Sprintf("%s%s", wikidataURL, data.Pageprops.WikibaseItem),
		}
	}

	if len(data.Protection) > 0 {
		for _, protection := range data.Protection {
			page.Protection = append(page.Protection, &schema.Protection{
				Type:   protection.Type,
				Level:  protection.Level,
				Expiry: protection.Expiry,
			})
		}
	}

	if len(data.WbEntityUsage) > 0 {
		for id, aspects := range data.WbEntityUsage {
			page.AdditionalEntities = append(page.AdditionalEntities, &schema.Entity{
				Identifier: id,
				URL:        fmt.Sprintf("%s%s", wikidataURL, id),
				Aspects:    aspects.Aspects,
			})
		}
	}

	if len(data.Categories) > 0 {
		for _, category := range data.Categories {
			if !category.Hidden {
				page.Categories = append(page.Categories, &schema.Page{
					Name: category.Title,
					URL:  fmt.Sprintf("%s/wiki/%s", f.Project.SiteURL, strings.ReplaceAll(category.Title, " ", "_")),
				})
			}
		}
	}

	if len(data.Templates) > 0 {
		for _, template := range data.Templates {
			page.Templates = append(page.Templates, &schema.Page{
				Name: template.Title,
				URL:  fmt.Sprintf("%s/wiki/%s", f.Project.SiteURL, strings.ReplaceAll(template.Title, " ", "_")),
			})
		}
	}

	if len(data.Redirects) > 0 {
		for _, redirect := range data.Redirects {
			page.Redirects = append(page.Redirects, &schema.Page{
				Name: redirect.Title,
				URL:  fmt.Sprintf("%s/wiki/%s", f.Project.SiteURL, strings.ReplaceAll(redirect.Title, " ", "_")),
			})
		}
	}

	return page
}
