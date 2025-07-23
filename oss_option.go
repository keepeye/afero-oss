package ossfs

import "github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"

type OSSOptionFunc func(c *oss.Config)

// Specify the endpoint.
func OSSWithEndpoint(endpoint string) OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithEndpoint(endpoint)
	}
}

// Use an internal endpoint.
func OSSWithUseInternalEndpoint() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithUseInternalEndpoint(true)
	}
}

// Use an OSS-accelerated endpoint.
func OSSWithUseAccelerateEndpoint() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithUseAccelerateEndpoint(true)
	}
}

// Use a dual-stack endpoint.
func OSSWithUseDualStackEndpoint() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithUseDualStackEndpoint(true)
	}
}

// Access OSS by using a custom domain name.
func OSSWithUseCName(domain string) OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithEndpoint(domain).WithUseCName(true)
	}
}

// Use path request style.
func OSSWithUsePathStyle() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithUsePathStyle(true)
	}
}

// Specifies that CRC-64 is disabled during object upload.
func OSSWithDisableUploadCRC64Check() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithDisableUploadCRC64Check(true)
	}
}

// Specifies that CRC-64 is disabled during object download.
func OSSWithDisableDownloadCRC64Check() OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithDisableDownloadCRC64Check(true)
	}
}

// Specifies that additional headers to be signed. It's valid in V4 signature.
func OSSWithAdditionalHeaders(headers []string) OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithAdditionalHeaders(headers)
	}
}

// Specifies user identifier appended to the User-Agent header.
func OSSWithUserAgent(ua string) OSSOptionFunc {
	return func(c *oss.Config) {
		c.WithUserAgent(ua)
	}
}
