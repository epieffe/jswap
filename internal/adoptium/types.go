package adoptium

type releases struct {
	Releases []string `json:"releases"`
}

type release struct {
	Binaries      []binary `json:"binaries"`
	DownloadCount int      `json:"download_count"`
	ID            string   `json:"id"`
	ReleaseLink   string   `json:"release_link"`
	ReleaseName   string   `json:"release_name"`
	ReleaseType   string   `json:"release_type"`
	Timestamp     string   `json:"timestamp"`
	UpdatedAt     string   `json:"updated_at"`
	Vendor        string   `json:"vendor"`
	VersionData   version  `json:"version_data"`
}

type asset struct {
	ReleaseLink string  `json:"release_link"`
	ReleaseName string  `json:"release_name"`
	Vendor      string  `json:"vendor"`
	Binary      binary  `json:"binary"`
	Version     version `json:"version"`
}

type binary struct {
	Architecture  string `json:"architecture"`
	DownloadCount int    `json:"download_count"`
	HeapSize      string `json:"heap_size"`
	ImageType     string `json:"image_type"`
	JvmImpl       string `json:"jvm_impl"`
	Os            string `json:"os"`
	Package       struct {
		Checksum      string `json:"checksum"`
		ChecksumLink  string `json:"checksum_link"`
		DownloadCount int    `json:"download_count"`
		Link          string `json:"link"`
		MetadataLink  string `json:"metadata_link"`
		Name          string `json:"name"`
		SignatureLink string `json:"signature_link"`
		Size          int    `json:"size"`
	} `json:"package"`
	Project   string `json:"project"`
	ScmRef    string `json:"scm_ref"`
	UpdatedAt string `json:"updated_at"`
}

type version struct {
	Build          int    `json:"build"`
	Major          int    `json:"major"`
	Minor          int    `json:"minor"`
	OpenjdkVersion string `json:"openjdk_version"`
	Optional       string `json:"optional"`
	Security       int    `json:"security"`
	Semver         string `json:"semver"`
}
