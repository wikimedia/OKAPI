package page

import (
	"fmt"
	"okapi-data-service/models"
	"okapi-data-service/schema/v3"
	"strings"
	"testing"
	"time"

	"github.com/protsack-stephan/mediawiki-api-client"
	"github.com/stretchr/testify/assert"
)

var factoryTestProject = &models.Project{
	DbName:   "afwikibooks",
	SiteName: "Wikibooks",
	SiteCode: "wikibooks",
	SiteURL:  "en.wikipedia.org",
}
var factoryTestLanguage = &models.Language{
	Code:      "af",
	LocalName: "Afrikaans",
}
var factoryTestNamespace = &models.Namespace{
	ID:    0,
	Title: "Article",
}
var factoryTestTitle = "Earth"
var factoryTestHTML = "...html goes here..."
var factoryTestWikitext = "...wikitext goes here..."
var factoryTestRevID = 122
var factoryTestRevQID = "Q2"

func TestFactory(t *testing.T) {
	assert := assert.New(t)

	fact := new(Factory)
	fact.Project = factoryTestProject
	fact.Language = factoryTestLanguage
	fact.Namespace = factoryTestNamespace

	data := &mediawiki.PageData{
		Revisions: []mediawiki.PageDataRevision{
			{},
		},
	}
	data.Title = factoryTestTitle
	data.LastRevID = factoryTestRevID
	data.Pageprops.WikibaseItem = factoryTestRevQID
	data.Revisions[0].Slots.Main.Content = factoryTestWikitext
	data.Revisions[0].Timestamp = time.Now().UTC()
	data.Categories = []mediawiki.PageDataCategory{
		{
			Ns:     schema.NamespaceCategory,
			Title:  "example",
			Hidden: false,
		},
		{
			Ns:     schema.NamespaceCategory,
			Title:  "example hidden category",
			Hidden: true,
		},
	}

	page := fact.Create(data, factoryTestHTML)
	assert.Equal(factoryTestTitle, page.Name)
	assert.Equal(factoryTestHTML, page.ArticleBody.HTML)
	assert.Equal(factoryTestWikitext, page.ArticleBody.Wikitext)
	assert.Equal(data.Revisions[0].Timestamp, *page.DateModified)
	assert.Equal(factoryTestRevID, page.Version.Identifier)
	assert.Equal(factoryTestProject.DbName, page.IsPartOf.Identifier)
	assert.Equal(factoryTestProject.SiteName, page.IsPartOf.Name)
	assert.Equal(factoryTestLanguage.Code, page.InLanguage.Identifier)
	assert.Equal(factoryTestLanguage.LocalName, page.InLanguage.Name)
	assert.Equal(factoryTestNamespace.ID, page.Namespace.Identifier)
	assert.Equal(factoryTestNamespace.Title, page.Namespace.Name)
	assert.Equal(factoryTestRevQID, page.MainEntity.Identifier)

	for _, category := range data.Categories {
		ctg := &schema.Page{
			Name: category.Title,
			URL:  fmt.Sprintf("%s/wiki/%s", factoryTestProject.SiteURL, strings.ReplaceAll(category.Title, " ", "_")),
		}

		if category.Hidden {
			assert.NotContains(page.Categories, ctg)
		} else {
			assert.Contains(page.Categories, ctg)
		}
	}

	for _, license := range page.License {
		assert.Equal(schema.LicenseIdentifier, license.Identifier)
		assert.Equal(schema.LicenseName, license.Name)
		assert.Equal(schema.LicenseURL, license.URL)
	}
}

// TestWikinewsLicense tests the license for wikinews projects.
func TestWikinewsLicense(t *testing.T) {
	assert := assert.New(t)

	fact := new(Factory)
	fact.Project = &models.Project{
		DbName:   "arwikinews",
		SiteName: "ويكي_الأخبار",
		SiteCode: "wikinews",
	}
	fact.Language = &models.Language{
		Code: "ar",
	}
	fact.Namespace = factoryTestNamespace

	data := &mediawiki.PageData{
		Revisions: []mediawiki.PageDataRevision{
			{},
		},
	}

	page := fact.Create(data, factoryTestHTML)

	for _, license := range page.License {
		assert.Equal("CC-BY-2.5", license.Identifier)
		assert.Equal("Attribution 2.5 Generic", license.Name)
		assert.Equal("https://creativecommons.org/licenses/by/2.5/", license.URL)
	}
}
