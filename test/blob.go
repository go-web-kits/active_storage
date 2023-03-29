package test

import (
	"bytes"
	"time"

	. "github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/active_storage/services"
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActiveStorageBlob", func() {
	var (
		blob Blob
		p    *MonkeyPatches
	)

	BeforeEach(func() {
		Config.Service = services.ASTerra{}
		blob = Blob{Key: time.Now().String()}
		factory.Create(&blob)
	})

	AfterEach(func() {
		CleanData(Models...)
		Reset(&blob)
		p.Check()
	})

	Describe("FindBlobBy", func() {
		It("finds the blob by given conditions, and preload attachments", func() {
			result := FindBlobBy(dbx.EQ{"key": blob.Key}, true)
			Expect(result).To(HaveFound())
			Expect(result.Data.(*Blob)).To(BeTheSameRecordTo(blob))
		})
	})

	Describe("CreateBlob", func() {
		It("creates the blob", func() {
			result := CreateBlob("abc.sig", uint(1), "md5")
			Expect(result).To(HaveAffected())
			Expect(result.Data.(*Blob)).To(HaveAttributes(Blob{Filename: "abc.sig", ByteSize: 1}))
			Expect(result.Data.(*Blob).Key).NotTo(BeZero())
		})
	})

	Describe(".SignedId", func() {
		It("returns the blob id", func() {
			Expect(blob.SignedId()).To(Equal(blob.ID))
		})
	})

	Describe(".Size", func() {
		It("returns file size (MB) of the blob", func() {
			blob.ByteSize = 1234567
			Expect(blob.Size()).To(Equal(1.18))
		})
	})

	Describe(".MD5", func() {
		It("returns nil if checksum is blank", func() {
			Expect(blob.MD5()).To(BeNil())
		})

		It("decodes the checksum to MD5 and returns that", func() {
			blob.Checksum = "l7COINA+R0b76eiI44nHNA=="
			Expect(blob.MD5()).To(Equal("97b08e20d03e4746fbe9e888e389c734"))
		})
	})

	Describe("DirectUploadInfo", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("DirectUploadInfo").AndPerform(
				func(_ services.ASTerra, _ *Blob) (interface{}, error) { return struct{}{}, nil }).Once()
			Expect(blob.DirectUploadInfo()).To(Equal(struct{}{}))
		})
	})

	Describe("Upload", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("Upload").AndReturn(nil).Once()
			Expect(blob.Upload(bytes.NewReader([]byte("")))).To(Succeed())
		})
	})

	Describe("URLWithHeader", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("URLWithHeader").AndReturn("", nil, nil).Once()
			Expect(blob.URLWithHeader()).To(Equal(""))
		})
	})

	Describe("URL", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("URL").AndReturn("", nil).Once()
			Expect(blob.URL()).To(Equal(""))
		})
	})

	Describe("Download", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("Download").AndReturn(nil, nil).Once()
			Expect(blob.Download()).To(BeZero())
		})
	})

	Describe("Delete", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("Delete").AndReturn(nil).Once()
			Expect(blob.Delete()).To(Succeed())
		})
	})

	Describe("Sync", func() {
		It("calls the interface", func() {
			p = ExpectAnyInstanceLike(Config.Service).ToCall("Sync").AndReturn(nil).Once()
			Expect(blob.Sync()).To(Succeed())
		})
	})
})
