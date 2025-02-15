schema_version = 1

project {
  license        = "Apache-2.0"
  copyright_holder = "Arthur Diniz <arthurbdiniz@gmail.com>"
  copyright_year = 2025

  header_ignore = [
    # examples used within documentation (prose)
    "examples/**",

    # GitHub issue template configuration
    ".github/ISSUE_TEMPLATE/*.yml",

    # golangci-lint tooling configuration
    ".golangci.yml",

    # GoReleaser tooling configuration
    ".goreleaser.yml",
  ]
}
