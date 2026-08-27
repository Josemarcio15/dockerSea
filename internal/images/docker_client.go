package images

// Docker operations for images are split by flow in list.go, pull.go and
// transfer.go. This file is the canonical adapter boundary for the module;
// new Docker calls should be added there or in one of those flow files rather
// than in service.go.
