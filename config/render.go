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
	pages, _ := filepath.Glob(filepath.Join(baseDir, "views/**/index.tmpl"))

	fmt.Println("\n🔧 Chargement des templates...")

	for _, page := range pages {
		relPath, _ := filepath.Rel(filepath.Join(baseDir, "views"), page)
		name := strings.TrimSuffix(filepath.ToSlash(relPath), filepath.Ext(relPath)) // ex: home/index

		files := append([]string{}, layouts...)
		files = append(files, partials...)
		files = append(files, contents...)
		files = append(files, page)

		fmt.Printf("\n📄 Enregistrement du template : %s\n", name)
		for _, f := range files {
			fmt.Printf("  └── %s\n", f)
		}

		tpl := template.New(name)

		tpl.Funcs(template.FuncMap{}) // Ajoute tes fonctions ici si besoin

		tpl, err := tpl.ParseFiles(files...)
		if err != nil {
			fmt.Printf("❌ ERREUR DE PARSING [%s]: %v\n", name, err)
			continue
		}

		for _, tmpl := range tpl.Templates() {
			fmt.Println("  →", tmpl.Name())
		}

		r.Add(name, tpl)
	}


	return r
}
