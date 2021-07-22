// Copyright 2017 The Hugo Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package i18n

import (
	"reflect"
	"strings"

	"github.com/gohugoio/hugo/common/hreflect"
	"github.com/gohugoio/hugo/common/loggers"
	"github.com/gohugoio/hugo/config"
	"github.com/gohugoio/hugo/helpers"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

type translateFunc func(translationID string, templateData interface{}) string

var (
	i18nWarningLogger = helpers.NewDistinctFeedbackLogger()
)

// Translator handles i18n translations.
type Translator struct {
	translateFuncs map[string]translateFunc
	cfg            config.Provider
	logger         loggers.Logger
}

// NewTranslator creates a new Translator for the given language bundle and configuration.
func NewTranslator(b *i18n.Bundle, cfg config.Provider, logger loggers.Logger) Translator {
	t := Translator{cfg: cfg, logger: logger, translateFuncs: make(map[string]translateFunc)}
	t.initFuncs(b)
	return t
}

// Func gets the translate func for the given language, or for the default
// configured language if not found.
func (t Translator) Func(lang string) translateFunc {
	if f, ok := t.translateFuncs[lang]; ok {
		return f
	}
	t.logger.Infof("Translation func for language %v not found, use default.", lang)
	if f, ok := t.translateFuncs[t.cfg.GetString("defaultContentLanguage")]; ok {
		return f
	}

	t.logger.Infoln("i18n not initialized; if you need string translations, check that you have a bundle in /i18n that matches the site language or the default language.")
	return func(translationID string, args interface{}) string {
		return ""
	}

}

func (t Translator) initFuncs(bndl *i18n.Bundle) {
	enableMissingTranslationPlaceholders := t.cfg.GetBool("enableMissingTranslationPlaceholders")
	for _, lang := range bndl.LanguageTags() {
		currentLang := lang
		currentLangStr := currentLang.String()
		// This may be pt-BR; make it case insensitive.
		currentLangKey := strings.ToLower(strings.TrimPrefix(currentLangStr, artificialLangTagPrefix))
		localizer := i18n.NewLocalizer(bndl, currentLangStr)
		t.translateFuncs[currentLangKey] = func(translationID string, templateData interface{}) string {

			var pluralCount interface{}

			if templateData != nil {
				tp := reflect.TypeOf(templateData)
				if hreflect.IsNumber(tp.Kind()) {
					pluralCount = templateData
					// This was how go-i18n worked in v1.
					templateData = map[string]interface{}{
						"Count": templateData,
					}

				}
			}

			translated, translatedLang, err := localizer.LocalizeWithTag(&i18n.LocalizeConfig{
				MessageID:    translationID,
				TemplateData: templateData,
				PluralCount:  pluralCount,
			})

			if err == nil && currentLang == translatedLang {
				return translated
			}

			if _, ok := err.(*i18n.MessageNotFoundErr); !ok {
				t.logger.Warnf("Failed to get translated string for language %q and ID %q: %s", currentLangStr, translationID, err)
			}

			if t.cfg.GetBool("logI18nWarnings") {
				i18nWarningLogger.Printf("i18n|MISSING_TRANSLATION|%s|%s", currentLangStr, translationID)
			}

			if enableMissingTranslationPlaceholders {
				return "[i18n] " + translationID
			}

			return translated
		}
	}
}
