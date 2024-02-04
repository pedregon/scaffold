package scaffold

type (
	// Loader is an on-load callback.
	Loader func(name string, next func(string) error) error
)

// SkipLoader skips loading Plugin(s) by name.
func SkipLoader(plugins ...string) Loader {
	return func(name string, next func(string) error) error {
		for _, plg := range plugins {
			if plg == name {
				return nil
			}
		}
		return next(name)
	}
}