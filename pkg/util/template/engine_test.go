package template

import (
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEngineRender(t *testing.T) {
	Convey("Engine.Render", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := NewMockTemplateResolver(ctrl)
		engine := &Engine{Resolver: resolver}

		resolver.EXPECT().Resolve(gomock.Any(), gomock.Any()).Return(&Resolved{
			T: T{
				Type:   "page-a",
				IsHTML: true,
				Defines: []string{
					`
					{{- define "define-a" -}}
					define-a
					{{- end -}}
					`,
				},
			},
			Content: `<!DOCTYPE html>
<html>
<head><title>Hi</title></head>
<body>
{{ template "define-a" }}
{{ template "component-a" }}
<p>{{ template "greeting" (makemap "URL" .URL) }}</p>
</body>
</html>`,
			Translations: map[string]Translation{
				"greeting": Translation{
					LanguageTag: "en",
					Value:       `<a href="{URL}">Hi</a>`,
				},
			},
			ComponentContents: []string{
				`
				{{- define "component-a" -}}
				component-a
				{{- end -}}
				`,
			},
		}, nil)

		out, err := engine.Render(&RenderContext{
			ValidatorOptions: []ValidatorOption{
				AllowRangeNode(true),
				AllowTemplateNode(true),
				AllowDeclaration(true),
				MaxDepth(15),
			},
		}, "", map[string]interface{}{
			"URL": "http://www.example.com",
		})
		So(err, ShouldBeNil)
		So(out, ShouldEqual, `<!DOCTYPE html>
<html>
<head><title>Hi</title></head>
<body>
define-a
component-a
<p><a href="http://www.example.com">Hi</a></p>
</body>
</html>`)
	})
}

func TestEngineTranslation(t *testing.T) {
	Convey("Engine.Translation", t, func() {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		resolver := NewMockTemplateResolver(ctrl)
		engine := &Engine{Resolver: resolver}

		resolver.EXPECT().ResolveTranslations(gomock.Any(), "translation.json").Return(map[string]Translation{
			`"greeting"`: {
				LanguageTag: "en",
				Value:       `Hello, {Name}`,
			},
		}, nil)

		ctx := &RenderContext{
			ValidatorOptions: []ValidatorOption{
				AllowRangeNode(true),
				AllowTemplateNode(true),
				AllowDeclaration(true),
				MaxDepth(15),
			},
		}

		Convey("should render translation entry", func() {
			translations, err := engine.Translation(ctx, "translation.json")
			So(err, ShouldBeNil)

			out, err := translations.RenderText(
				`"greeting"`,
				map[string]interface{}{
					"Name": "<User>",
				},
			)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, `Hello, <User>`)

			out, err = translations.RenderHTML(
				`"greeting"`,
				map[string]interface{}{
					"Name": "<User>",
				},
			)
			So(err, ShouldBeNil)
			So(out, ShouldEqual, `Hello, &lt;User&gt;`)
		})

		Convey("should fail when rendering non-existing translation entry", func() {
			translations, err := engine.Translation(ctx, "translation.json")
			So(err, ShouldBeNil)

			_, err = translations.RenderText(
				"none",
				map[string]interface{}{
					"URL": "http://www.example.com",
				},
			)
			So(IsNotFound(err), ShouldBeTrue)
		})
	})
}
