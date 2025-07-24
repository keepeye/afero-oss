package ossfs

import (
	"testing"

	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/stretchr/testify/assert"
)

func getNewOssCfgForTest() *oss.Config {
	return oss.LoadDefaultConfig()
}

func TestOSSWithEndpoint(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	endpoint := "testEndpoint"
	t.Run("test OSSWithEndpoint", func(t *testing.T) {
		OSSWithEndpoint(endpoint)(ossCfg)
		assert.Equal(t, &endpoint, ossCfg.Endpoint)
	})
}

func TestOSSWithUseInternalEndpoint(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithUseInternalEndpoint", func(t *testing.T) {
		assert.Zero(t, ossCfg.UseInternalEndpoint)
		OSSWithUseInternalEndpoint()(ossCfg)
		assert.Equal(t, true, *ossCfg.UseInternalEndpoint)
	})
}

func TestOSSWithUseAccelerateEndpoint(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithUseAccelerateEndpoint", func(t *testing.T) {
		assert.Zero(t, ossCfg.UseAccelerateEndpoint)
		OSSWithUseAccelerateEndpoint()(ossCfg)
		assert.Equal(t, true, *ossCfg.UseAccelerateEndpoint)
	})
}

func TestOSSWithUseDualStackEndpoint(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithUseDualStackEndpoint", func(t *testing.T) {
		assert.Zero(t, ossCfg.UseDualStackEndpoint)
		OSSWithUseDualStackEndpoint()(ossCfg)
		assert.Equal(t, true, *ossCfg.UseDualStackEndpoint)
	})
}

func TestOSSWithUseCName(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	domain := "testDomain"
	t.Run("test OSSWithUseCName", func(t *testing.T) {
		assert.Zero(t, ossCfg.UseCName)
		OSSWithUseCName(domain)(ossCfg)
		assert.Equal(t, true, *ossCfg.UseCName)
		assert.Equal(t, domain, *ossCfg.Endpoint)
	})
}

func TestOSSWithUsePathStyle(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithUsePathStyle", func(t *testing.T) {
		assert.Zero(t, ossCfg.UsePathStyle)
		OSSWithUsePathStyle()(ossCfg)
		assert.Equal(t, true, *ossCfg.UsePathStyle)
	})
}

func TestOSSWithDisableUploadCRC64Check(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithDisableUploadCRC64Check", func(t *testing.T) {
		assert.Zero(t, ossCfg.DisableUploadCRC64Check)
		OSSWithDisableUploadCRC64Check()(ossCfg)
		assert.Equal(t, true, *ossCfg.DisableUploadCRC64Check)
	})
}

func TestOSSWithDisableDownloadCRC64Check(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	t.Run("test OSSWithDisableDownloadCRC64Check", func(t *testing.T) {
		assert.Zero(t, ossCfg.DisableDownloadCRC64Check)
		OSSWithDisableDownloadCRC64Check()(ossCfg)
		assert.Equal(t, true, *ossCfg.DisableDownloadCRC64Check)
	})
}

func TestOSSWithAdditionalHeaders(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	headers := []string{
		"x-test: ok",
	}
	t.Run("test OSSWithAdditionalHeaders", func(t *testing.T) {
		assert.Zero(t, ossCfg.AdditionalHeaders)
		OSSWithAdditionalHeaders(headers)(ossCfg)
		assert.Equal(t, headers, ossCfg.AdditionalHeaders)
	})
}

func TestOSSWithUserAgent(t *testing.T) {
	ossCfg := getNewOssCfgForTest()
	ua := "testUa"
	t.Run("test OSSWithUserAgent", func(t *testing.T) {
		assert.Zero(t, ossCfg.UserAgent)
		OSSWithUserAgent(ua)(ossCfg)
		assert.Equal(t, ua, *ossCfg.UserAgent)
	})
}
