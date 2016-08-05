package loglight

type packageNameFilter struct {
  packageNames map[string]bool
  isWhitelist  bool
}

func NewPackageNameFilter(packageNames []string, isWhitelist bool) PackageFilter {

  filter := &packageNameFilter{isWhitelist: isWhitelist, packageNames: make(map[string]bool)}

  for _, name := range packageNames {
    filter.packageNames[name] = true
  }

  return filter
}

func(fiter packageNameFilter) Filter(packageName string) bool {
  return fiter.packageNames[packageName] && fiter.isWhitelist
}

func(packageFilter NullPackageFilter) Filter(packageName string) bool {
  return true
}
