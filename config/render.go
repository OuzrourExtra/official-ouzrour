package config

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-contrib/multitemplate"
)

func LoadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	baseDir, _ := os.Getwd()

	layouts, _ := filepath.Glob(filepath.Join(baseDir, "views/layouts/*.tmpl"))
	partials, _ := filepath.Glob(filepath.Join(baseDir, "views/partials/*.tmpl"))
	contents, _ := filepath.Glob(filepath.Join(baseDir, "views/**/content.tmpl"))

	// 🔥 ICI : on récupère TOUS les .tmpl SAUF layouts/ et partials/
	allPages, _ := filepath.Glob(filepath.Join(baseDir, "views/**/*.tmpl"))

	fmt.Println("\n🔧 Chargement des templates...")

	for _, page := range allPages {
		// ⚠️ On ignore les layouts et les partials (déjà inclus)
		if strings.Contains(page, "/layouts/") || strings.Contains(page, "\\layouts\\") {
			continue
		}
		if strings.Contains(page, "/partials/") || strings.Contains(page, "\\partials\\") {
			continue
		}
		if strings.HasSuffix(page, "content.tmpl") {
			continue
		}

		// ex: views/login.tmpl → "login"
		// ex: views/tasks/tasks.tmpl → "tasks/tasks"
		relPath, _ := filepath.Rel(filepath.Join(baseDir, "views"), page)
		name := strings.TrimSuffix(filepath.ToSlash(relPath), filepath.Ext(relPath))

		files := append([]string{}, layouts...)
		files = append(files, partials...)
		files = append(files, contents...)
		files = append(files, page)

		fmt.Printf("\n📄 Enregistrement du template : %s\n", name)
		for _, f := range files {
			fmt.Printf("  └── %s\n", f)
		}

		tpl := template.New(name).Funcs(template.FuncMap{})
		tpl, err := tpl.ParseFiles(files...)
		if err != nil {
			fmt.Printf("❌ ERREUR DE PARSING [%s]: %v\n", name, err)
			continue
		}

		for _, t := range tpl.Templates() {
			fmt.Println("  →", t.Name())
		}

		r.Add(name, tpl)
	}

	fmt.Println("\n✅ Tous les templates ont été chargés.\n")
	return r
}
