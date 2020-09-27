package namespaces

import (
	"fmt"
	"net/http"
	"okapi/lib/wiki"
	"strings"

	"okapi/lib/task"
	"okapi/models"
)

// Worker namespaces unit processor
func Worker(ctx *task.Context) task.Worker {
	return func(id int, payload task.Payload) (string, map[string]interface{}, error) {
		message := "name:'%s', id:#%d"
		project := payload.(*models.Project)
		info := map[string]interface{}{
			"_db_name":   project.DBName,
			"_lang_name": project.LangName,
		}

		namespaces, status, err := wiki.Client(project.SiteURL).GetNamespaces()

		if err != nil {
			return "", info, err
		}

		if status != http.StatusOK {
			return "", info, fmt.Errorf("no namespaces for the project")
		}

		for _, ns := range namespaces {
			namespace := models.Namespace{
				ID:    ns.ID,
				Lang:  project.Lang,
				Title: ns.Name,
			}

			if namespace.ID == 0 && len(namespace.Title) <= 0 {
				namespace.Title = "Article"
			}

			err := models.DB().
				Model(&namespace).
				Column("created_at", "updated_at").
				Where("id = ? and lang = ?", ns.ID, project.Lang).
				Select()

			if err != nil {
				models.DB().Insert(&namespace)
			} else {
				models.DB().
					Model(&namespace).
					Where("id = ? and lang = ?", ns.ID, project.Lang).
					Update()
			}

			title := strings.ReplaceAll(ns.Name, " ", "_")

			if len(title) > 0 && ns.ID != 0 {
				resp, err := models.
					DB().
					Exec("update pages set ns_id = ? where title like ? and project_id = ?", ns.ID, title+":%", project.ID)

				if err != nil {
					ctx.Log.Error("update query failed failed", err.Error())
				} else {
					ctx.Log.Info(
						"project ns update",
						fmt.Sprintf("project: '%s', namespace: '%s', ns_id: '%d', rows: '%d'", project.DBName, title, ns.ID, resp.RowsAffected()))
				}
			}
		}

		return fmt.Sprintf(message, project.DBName, project.ID), info, err
	}
}
