package services_test

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/go-web-kits/active_storage"
	. "github.com/go-web-kits/active_storage/services"
	service "github.com/go-web-kits/active_storage/services/terra"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/let"
	"github.com/pkg/errors"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Terra", func() {
	var (
		blob  Blob
		terra ASTerra
		p     *MonkeyPatches
	)

	BeforeEach(func() {
		blob = Blob{
			Key:      "key",
			Filename: "abc.txt",
		}
	})

	AfterEach(func() {
		p.Check()
	})

	Describe("Upload", func() { // <<---
		BeforeEach(func() {
			IsExpectedToCall(service.DirectUpload).
				AndPerform(func(_ io.Reader, _, _ string, _ int64, _ []time.Duration) (string, error) { return "file_id", nil })
		})

		It("returns the upload error", func() {
			Expect(terra.Upload(&blob, bytes.NewReader([]byte{}), "checksum")).To(Succeed())
		})

		It("changes key of the blob", func() {
			Expect(blob.Key).To(Equal("key"))
			_ = terra.Upload(&blob, bytes.NewReader([]byte{}), "checksum")
			Expect(blob.Key).To(Equal("file_id"))
		})
	})

	Describe("Download", func() { // <<---
		BeforeEach(func() {
			ExpectAnyInstanceLike(&http.Client{}).ToCall("Do").AndReturn(
				&http.Response{Body: ioutil.NopCloser(bytes.NewReader([]byte("body")))}, nil)
		})

		It("does successfully", func() {
			IsExpectedToCall(service.RequestCachedURL).AndReturn(service.DownloadInfo{}, nil)
			body, err := terra.Download(blob)

			Expect(err).NotTo(HaveOccurred())
			Expect(body).To(Equal([]byte("body")))
		})

		When("error occurs while getURL", func() {
			It("returns error", func() {
				IsExpectedToCall(service.RequestCachedURL).AndReturn(service.DownloadInfo{}, errors.New(""))
				_, err := terra.Download(blob)

				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("URL", func() { // <<---
		It("returns the url", func() {
			IsExpectedToCall(service.RequestCachedURL).AndReturn(service.DownloadInfo{URL: "url"}, nil)

			Expect(terra.URL(blob)).To(Equal("url"))
		})
	})

	Describe("DirectUploadInfo", func() { // <<---
		When("RequestSTS failed", func() {
			It("returns the error", func() {
				IsExpectedToCall(service.RequestSTS).AndReturn(&service.STS{}, errors.New(""))

				_, err := terra.DirectUploadInfo(&blob)
				Expect(err).To(HaveOccurred())
			})
		})

		BeforeEach(func() {
			IsExpectedToCall(service.RequestSTS).AndReturn(&service.STS{FileID: "file_id"}, nil)
		})

		When("the update fails", func() {
			It("updates the blob's key but returns the GORM error", func() {
				let.UpdateBy().Fail()
				_, err := terra.DirectUploadInfo(&blob)
				Expect(err).To(HaveOccurred())
			})
		})

		It("does successfully", func() {
			let.UpdateBy().Succeed()
			Expect(blob.Key).To(Equal("key"))
			sts, err := terra.DirectUploadInfo(&blob)
			Expect(err).NotTo(HaveOccurred())
			Expect(sts.(*service.STS).FileID).To(Equal("file_id"))
			Expect(blob.Key).To(Equal("file_id"))
		})
	})

	Describe("Sync", func() { // <<---
		BeforeEach(func() {
			service.ExtranetBuckets = []string{"foo"}
		})

		It("does successfully", func() {
			p = IsExpectedToCall(service.RequestSync).AndReturn(nil).AtLeastOnce()
			Expect(terra.Sync(&blob)).To(Succeed())
		})
	})
})
