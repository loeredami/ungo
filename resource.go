package ungo

import (
	"fmt"
	"os"
)

func PackResources(file_paths []string, output_path string) {
	WithTempFile(func(temp *os.File) {
		table := NewSmallMap[string, []byte](len(file_paths))
		for _, path := range file_paths {
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			table.Set(path, data)
		}

		header, data := func() ([][]byte, [][]byte) {
			var headers [][]byte
			var contents [][]byte
			table.ForEach(func(name string, content []byte) {
				headers = append(headers, []byte(name))
				contents = append(contents, content)
			})
			return headers, contents
		}()

		for i, h := range header {
			_, _ = temp.Write(h)
			_, _ = temp.Write(data[i])
		}

		var temp_data []byte
		_, _ = temp.Read(temp_data)
		os.WriteFile(output_path, temp_data, 0644)
	})
}

type Package struct {
	Files *SmallMap[string, []byte]
}

func LoadPackage(path string) (*Package, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	pkg := &Package{Files: NewSmallMap[string, []byte](len(data))}
	header, contents := func() ([][]byte, [][]byte) {
		var headers [][]byte
		var contents [][]byte
		for i := 0; i < len(data); {
			name_len := int(data[i])
			i++
			headers = append(headers, data[i:i+name_len])
			i += name_len
			contents = append(contents, data[i:])
			break
		}
		return headers, contents
	}()
	for i, h := range header {
		pkg.Files.Set(string(h), contents[i])
	}
	return pkg, nil
}

func (p *Package) Get(name string) ([]byte, bool) {
	return p.Files.Get(name)
}

type ResourceLoader struct {
	packages    []Lazy[*Package]
	loose_files SmallMap[string, Lazy[[]byte]]
}

func (r *ResourceLoader) Get(name string) Result[[]byte] {
	for _, pkg := range r.packages {
		if value, ok := pkg.Value().Get(name); ok {
			return VSuccess(value)
		}
	}

	if data, ok := r.loose_files.Get(name); ok {
		return VSuccess(data.Value())
	}

	if data, err := os.ReadFile(name); err == nil {
		r.loose_files.Set(name, NewLazy(func() []byte {
			data, _ := os.ReadFile(name)
			return data
		}))
		return VSuccess(data)
	}

	return VFail[[]byte](fmt.Errorf("resource not found %s", name))
}
