package util

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nicksnyder/go-i18n/i18n"
)

const DefaultLocale = "en"

var locales = make(map[string]string)

// FindDir looks for the given directory in nearby ancestors, falling back to `./` if not found.
func FindDir(dir string) (string, bool) {
	for _, parent := range []string{".", "..", "../.."} {
		foundDir, err := filepath.Abs(filepath.Join(parent, dir))
		if err != nil {
			continue
		} else if _, err := os.Stat(foundDir); err == nil {
			return foundDir, true
		}
	}
	return "./", false
}

// InitTranslations TBD
func InitTranslations(dir string) error {
	i18nDirectory, found := FindDir(dir)
	if !found {
		return fmt.Errorf("unable to find i18n directory")
	}

	files, _ := os.ReadDir(i18nDirectory)
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".json" {
			filename := f.Name()
			locales[strings.Split(filename, ".")[0]] = filepath.Join(i18nDirectory, filename)
			if err := i18n.LoadTranslationFile(filepath.Join(i18nDirectory, filename)); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetTranslationFunc TBD
func GetTranslationFuncFromReq(HTTPRequest *http.Request) i18n.TranslateFunc {
	return HTTPRequest.Context().Value("T").(i18n.TranslateFunc)
}

// GetTranslations TBD
func GetTranslationFunc(locale string) i18n.TranslateFunc {
	return TFuncWithFallback(locale)
}

// TFuncWithFallback TBD
func TFuncWithFallback(pref string) i18n.TranslateFunc {
	t, _ := i18n.Tfunc(pref)
	return func(translationID string, args ...interface{}) string {
		if translated := t(translationID, args...); translated != translationID {
			return translated
		}
		t, _ := i18n.Tfunc(DefaultLocale)
		return t(translationID, args...)
	}
}

func Translate(T i18n.TranslateFunc, message string) string {
	if T == nil {
		return message
	}
	return T(message)
}
