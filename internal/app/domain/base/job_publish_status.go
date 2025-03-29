package base

type JobPublishStatus int

const (
	JobPublishStatusUnpublished JobPublishStatus = iota
	JobPublishStatusPublished
)

func (s JobPublishStatus) String() string {
	return [...]string{"unpublished", "published"}[s]
}
