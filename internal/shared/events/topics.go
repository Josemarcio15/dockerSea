package events

const (
	ContainerEvent        = "docker:container:event"
	ImagePullProgress     = "docker:image:pull:progress"
	ImagePullComplete     = "docker:image:pull:complete"
	ImageTransferProgress = "docker:image:transfer:progress"
	ImageTransferComplete = "docker:image:transfer:complete"
)
