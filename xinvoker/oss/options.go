package xoss

type options struct {
	Debug           bool   `mapStructure:"debug"`
	Mode            string `mapStructure:"mode"`
	Addr            string `mapStructure:"addr"`
	AccessKeyID     string `mapStructure:"accessKeyId"`
	AccessKeySecret string `mapStructure:"accessKeySecret"`
	CdnName         string `mapStructure:"cdnName"`
	OssBucket       string `mapStructure:"ossBucket"`
	FileBucket      string `mapStructure:"fileBucket"`
	IsDeleteSrcPath bool   `mapStructure:"isDeleteSrcPath"`
}

func newOssOptions() *options {
	return &options{
		Debug:           false,
		Mode:            "file",
		Addr:            "",
		AccessKeyID:     "",
		AccessKeySecret: "",
		CdnName:         "",
		OssBucket:       "",
		FileBucket:      ".",
		IsDeleteSrcPath: false,
	}
}
