// Copyright (c) 2023 FajarLaksono. All Rights Reserved.

package pageversion

// type Service struct {
// 	name                string
// 	realm               string
// 	buildDate           string
// 	version             string
// 	gitHash             string
// 	versionRolesSeeding string
// }

// func New(name, realm, buildDate, gitHash, serviceRootPath string) (*Service, error) {
// 	version, versionRolesSeeding, err := loadAndParseVersionFile(serviceRootPath)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "load and parse version file")
// 	}

// 	return &Service{
// 		name:                name,
// 		realm:               realm,
// 		buildDate:           buildDate,
// 		version:             version,
// 		gitHash:             gitHash,
// 		versionRolesSeeding: versionRolesSeeding,
// 	}, nil
// }

// func loadAndParseVersionFile(serviceRootPath string) (string, string, error) {
// 	f, err := os.Open(serviceRootPath + "/version.json")
// 	if err != nil {
// 		return "", "", errors.Wrap(err, "open version file")
// 	}

// 	defer func() {
// 		if err := f.Close(); err != nil {
// 			logrus.WithError(err).Error("unable to close file")
// 		}
// 	}()

// 	byteVal, err := io.ReadAll(f)
// 	if err != nil {
// 		return "", "", errors.Wrap(err, "io read all version file")
// 	}

// 	versionInfo := make(map[string]interface{})
// 	if err := json.Unmarshal(byteVal, &versionInfo); err != nil {
// 		return "", "", errors.Wrap(err, "json unmarshal version byte value")
// 	}

// 	version, ok := versionInfo["version"]
// 	if !ok {
// 		return "", "", fmt.Errorf("version information is not available on json file: %s", serviceRootPath+"/version.json")
// 	}
// 	versionRolesSeeding, ok := versionInfo["version-roles-seeding"]
// 	if !ok {
// 		return "", "", fmt.Errorf("version-roles-seeding information is not available on json file: %s", serviceRootPath+"/version.json")
// 	}

// 	return version.(string), versionRolesSeeding.(string), nil
// }

// func (s *Service) GetAllVersionInfo() map[string]string {
// 	return map[string]string{
// 		"name":                  s.name,
// 		"realm":                 s.realm,
// 		"buildDate":             s.buildDate,
// 		"version":               s.version,
// 		"gitHash":               s.gitHash,
// 		"version-roles-seeding": s.versionRolesSeeding,
// 	}
// }

// func (s *Service) GetVersionInfo() string {
// 	return s.version
// }
