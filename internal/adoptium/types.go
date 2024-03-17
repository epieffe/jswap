package adoptium

type availableReleases struct {
	Releases          []int `json:"available_releases"`
	LTS               []int `json:"available_lts_releases"`
	MostRecentLTS     int   `json:"most_recent_lts"`
	MostRecentRelease int   `json:"most_recent_feature_release"`
}

type asset struct {
	ReleaseLink string `json:"release_link"`
	ReleaseName string `json:"release_name"`
	Vendor      string `json:"vendor"`
	Binary      struct {
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
	} `json:"binary"`
	Version struct {
		Build          int    `json:"build"`
		Major          int    `json:"major"`
		Minor          int    `json:"minor"`
		OpenjdkVersion string `json:"openjdk_version"`
		Optional       string `json:"optional"`
		Security       int    `json:"security"`
		Semver         string `json:"semver"`
	} `json:"version"`
}
