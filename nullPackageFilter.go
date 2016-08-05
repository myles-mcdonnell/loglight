package loglight

type PackageFilter interface {
  Filter(packageName string) bool
}

type NullPackageFilter struct {}
