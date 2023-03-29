package test

import (
	"io/ioutil"
	"mime/multipart"
	"time"

	. "github.com/go-web-kits/active_storage"
	"github.com/go-web-kits/dbx"
	. "github.com/go-web-kits/testx"
	"github.com/go-web-kits/testx/factory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ActiveStorageAttachment", func() {
	var (
		post       Post
		blob       Blob
		attachment Attachment
		opt        dbx.Opt
		p          *MonkeyPatches
	)

	BeforeEach(func() {
		factory.Create(&post)
		blob = Blob{Key: time.Now().String()}
		factory.Create(&blob)
	})

	AfterEach(func() {
		CleanData(Models...)
		Reset(&post, &blob, &attachment)
		p.Check()
	})

	Describe("Attach", func() {
		BeforeEach(func() {
			opt = dbx.Opt{RelatedWith: post, AssocField: "Picture", Preload: "Blob"}

		})

		When("passing blob ID", func() {
			It("attaches the blob which has the specified id", func() {
				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())

				Expect(Attach(blob.ID, Attachment{OwnerID: post.ID, OwnerType: "Post"})).To(HaveAffected())

				Expect(dbx.Find(&attachment, nil, opt)).To(HaveFound())
				Expect(attachment.Blob).To(BeTheSameRecordTo(blob))
			})

			It("can not attach the blob that does not exist", func() {
				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())

				Expect(Attach(blob.ID+1, Attachment{OwnerID: post.ID, OwnerType: "Post"})).NotTo(HaveAffected())

				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())
			})
		})

		When("passing blob", func() {
			It("attaches the blob", func() {
				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())

				Expect(Attach(blob, Attachment{OwnerID: post.ID, OwnerType: "Post"})).To(HaveAffected())

				Expect(dbx.Find(&attachment, nil, opt)).To(HaveFound())
				Expect(attachment.Blob).To(BeTheSameRecordTo(blob))
			})
		})

		When("passing `*multipart.FileHeader` attachable", func() {
			BeforeEach(func() {
				p = ExpectAnyInstanceLike(&blob).ToCall("Upload").AndReturn(nil)
				p.IsExpectedToCall(ioutil.ReadAll).AndReturn([]byte("hello"), nil)
			})

			It("uploads the file, creates a blob, then attaches the blob", func() {
				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())

				fh := multipart.FileHeader{Filename: "abc.sig", Size: int64(1234)}
				Expect(Attach(&fh, Attachment{OwnerID: post.ID, OwnerType: "Post"})).To(HaveAffected())

				Expect(dbx.Find(&attachment, nil, opt)).To(HaveFound())
				Expect(attachment.Blob).To(HaveAttributes(Blob{Filename: "abc.sig", ByteSize: 1234}))
				Expect(attachment.Blob.Key).NotTo(BeZero())
			})
		})

		When("passing `OpenedFileHeader` attachable", func() {
			BeforeEach(func() {
				p = ExpectAnyInstanceLike(&blob).ToCall("Upload").AndReturn(nil)
				p.IsExpectedToCall(ioutil.ReadAll).AndReturn([]byte("hello"), nil)
			})

			It("uploads the file, creates a blob, then attaches the blob", func() {
				Expect(dbx.Find(&attachment, nil, opt)).NotTo(HaveFound())

				file := OpenedFileHeader{&multipart.FileHeader{Filename: "abc.sig", Size: int64(1234)}, []byte("abc"), "md5"}
				Expect(Attach(&file, Attachment{OwnerID: post.ID, OwnerType: "Post"})).To(HaveAffected())

				Expect(dbx.Find(&attachment, nil, opt)).To(HaveFound())
				Expect(attachment.Blob).To(HaveAttributes(Blob{Filename: "abc.sig", ByteSize: 1234}))
				Expect(attachment.Blob.Key).NotTo(BeZero())
			})
		})
	})
})
