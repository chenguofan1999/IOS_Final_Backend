package model

func init() {
	Connect()
	CreateUserTableIfNotExists()
	CreateContentTableIfNotExists()
	CreateCommentTableIfNotExists()
	CreateReplyTableIfNotExists()
	CreateFollowTableIfNotExists()
	CreateLikeContentTableIfNotExists()
	CreateLikeCommentTableIfNotExists()
	CreateLikeReplyTableIfNotExists()
	CreateUserTagsTableIfNotExists()
	CreateContentTagsTableIfNotExists()
	CreateHistoryTableIfNotExists()
}
