package info

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type PkgInfo struct {
	Name    string
	Version string //first version in conandata.yml
	URLs    []string
}

func (p *PkgInfo) String() string {
	return fmt.Sprintf("Package: %s\nVersion: %s\nURLs: %v\n", p.Name, p.Version, p.URLs)
}

func ReadPackageInfoWithReturn(packageName, directory string) (*PkgInfo, error) {
	pkgInfo := &PkgInfo{
		Name: packageName,
	}
	packageDir := filepath.Join(directory, packageName)

	if _, err := os.Stat(packageDir); os.IsNotExist(err) {
		return nil, fmt.Errorf("package %s does not exist", packageName)
	}

	versions, err := os.ReadDir(packageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read version directories: %v", err)
	}

	if len(versions) == 0 {
		return nil, fmt.Errorf("no versions found for package %s", packageName)
	}

	var firstVersion os.DirEntry
	for _, v := range versions {
		if v.IsDir() {
			firstVersion = v
			break
		}
	}

	if firstVersion == nil {
		return nil, fmt.Errorf("no version directories found for package %s", packageName)
	}

	versionName := firstVersion.Name()

	// Read data.path file
	dataPathFile := filepath.Join(packageDir, versionName, "data.path")
	filePath, err := readFilePath(dataPathFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read data.path: %v %s", err, dataPathFile)
	}

	// Read content of the specified file
	content, err := readFileContent(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v %s", err, filePath)
	}

	// Parse YAML content
	var node yaml.Node
	err = yaml.Unmarshal([]byte(content), &node)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML content: %v %s", err, filePath)
	}

	var root *yaml.Node
	if node.Kind == yaml.DocumentNode && len(node.Content) > 0 {
		root = node.Content[0]
	}

	sources, err := getValueByKey(root, "sources")
	if err != nil {
		return nil, fmt.Errorf("failed to get sources: %v", err)
	}

	pkgInfo.Version = sources.Content[0].Value
	versionMap := sources.Content[1]
	if versionMap.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("sources first version is not a mapping node: %v", versionMap)
	}

	urlNode, err := getValueByKey(versionMap, "url")
	if err != nil {
		return nil, fmt.Errorf("failed to get url: %v", err)
	}

	switch urlNode.Kind {
	case yaml.SequenceNode:
		for _, url := range urlNode.Content {
			pkgInfo.URLs = append(pkgInfo.URLs, url.Value)
		}
	case yaml.ScalarNode:
		pkgInfo.URLs = append(pkgInfo.URLs, urlNode.Value)
	default:
		return nil, fmt.Errorf("unsupported URL format: %v", urlNode.Kind)
	}
	return pkgInfo, nil
}

func readFilePath(pathFile string) (string, error) {
	file, err := os.Open(pathFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return "", fmt.Errorf("data.path is empty")
}

func readFileContent(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getValueByKey(node *yaml.Node, key string) (*yaml.Node, error) {
	if node.Kind != yaml.MappingNode {
		return nil, fmt.Errorf("node is not a mapping node: %v", node)
	}

	for i := 0; i < len(node.Content); i += 2 {
		if node.Content[i].Value == key {
			return node.Content[i+1], nil
		}
	}

	return nil, fmt.Errorf("key not found: %s", key)
}
