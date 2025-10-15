package changelog

// Changelog represents the final, curated data before Markdown generation.
type Changelog struct {
    Version  string
    Date     string
    Sections []ChangelogSection
}

// ChangelogSection contains a list of items for a specific category.
type ChangelogSection struct {
    Title string
    Items []string
}


