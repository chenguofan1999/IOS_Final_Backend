package main

import (
	"fmt"
	_ "ios/model"
	"ios/routes"
	// "ios/routes"
)

func main() {
	router := routes.InitRouter()
	router.Run(":5009")
	// test1()
}

func test1() {
	// model.InsertUser("Alice", "678") // userID = 1
	// model.InsertUser("Bob", "678")   // userID = 2
	// model.InsertUser("Tom", "678")   // userID = 3
	// model.InsertUser("Joe", "678")   // userID = 4

	var err error
	checkError(err)

	//// ==========================================================================================================
	// err = model.InsertContent("trump is a loser", "indeed", "www.loser.com", "www.loser.com", 1000000000, 5)
	// checkError(err)
	// err = model.InsertContent("Trump Is A Loser", "Indeed", "www.loser.com", "www.loser.com", 1000000000, 2)
	// checkError(err)
	// err = model.InsertContent("TRUMP IS A LOSER", "INDEED", "www.loser.com", "www.loser.com", 1000000000, 3)
	// checkError(err)

	// fmt.Println(model.QueryUserIDWithName("Joe"))
	// fmt.Println(model.QueryUserIDWithName("Jon"))
	// fmt.Println(model.QueryMiniUserWithUserID(2))
	// fmt.Println(model.QueryMiniUserWithUserID(5))

	//// ==========================================================================================================
	// err = model.InsertComment(1, 2, "absolutely", 10000000001)
	// checkError(err)

	// fmt.Println(model.QueryCommentNumWithContentID(1)) // expected: 0
	// fmt.Println(model.QueryCommentNumWithContentID(2)) // expected: 1
	// fmt.Println(model.QueryCommentWithCommentID(1))    // expected: &{...}
	// fmt.Println(model.QueryCommentWithCommentID(2))    // expected: nil
	// fmt.Println(model.QueryCommentsWithContentID(1))   // expected: []
	// fmt.Println(model.QueryCommentsWithContentID(2))   // expected: [{...}]

	//// ==========================================================================================================
	// err = model.InsertReply(4, 1, "agree", 1909090900)
	// checkError(err)

	// fmt.Println(model.QueryReplyNumWithCommentID(2)) // expected: 0, no such comment
	// fmt.Println(model.QueryReplyNumWithCommentID(1)) // expected: 1, nil
	// fmt.Println(model.QueryReplyWithReplyID(1))      // expected: &{...}
	// fmt.Println(model.QueryReplyWithReplyID(2))      // expected: nil
	// fmt.Println(model.QueryRepliesWithCommentID(1))  // expected: [{...}]
	// fmt.Println(model.QueryRepliesWithCommentID(2))  // expected: []

	//// ==========================================================================================================
	// err = model.InsertLikeContent(1, 1)
	// checkError(err)
	// err = model.InsertLikeContent(1, 1) // expected error: alreadyliked
	// checkError(err)
	// err = model.InsertLikeContent(5, 1) // expected error: no such user
	// checkError(err)
	// err = model.InsertLikeContent(1, 5) // expected error: no such content
	// checkError(err)
	// err = model.InsertLikeComment(1, 1) // expected error: alreadyliked
	// checkError(err)
	// err = model.InsertLikeComment(5, 1) // expected error: no such user
	// checkError(err)
	// err = model.InsertLikeComment(1, 2) // expected error: no such comment
	// checkError(err)
	// err = model.InsertLikeReply(1, 1) // expected error: alreadyliked
	// checkError(err)
	// err = model.InsertLikeReply(5, 1) // expected error: no such user
	// checkError(err)
	// err = model.InsertLikeReply(1, 2) // expected error: no such reply
	// checkError(err)

	// fmt.Println(model.QueryLikeNumWithContentID(1))
	// fmt.Println(model.QueryLikeNumWithContentID(5))
	// fmt.Println(model.QueryLikeNumWithCommentID(1))
	// fmt.Println(model.QueryLikeNumWithCommentID(2))
	// fmt.Println(model.QueryLikeNumWithReplyID(1))
	// fmt.Println(model.QueryLikeNumWithReplyID(2))

	// fmt.Println(model.QueryHasLikedContent(1, 1))
	// fmt.Println(model.QueryHasLikedContent(1, 2))
	// fmt.Println(model.QueryHasLikedComment(1, 1))
	// fmt.Println(model.QueryHasLikedComment(1, 2))
	// fmt.Println(model.QueryHasLikedReply(1, 1))
	// fmt.Println(model.QueryHasLikedReply(1, 2))

	//// ==========================================================================================================
	// err = model.InsertFollow(1, 1)
	// checkError(err)
	// err = model.InsertFollow(1, 2)
	// checkError(err)
	// err = model.InsertFollow(1, 2)
	// checkError(err)
	// err = model.InsertFollow(2, 1)
	// checkError(err)
	// err = model.InsertFollow(2, 3)
	// checkError(err)
	// err = model.InsertFollow(2, 4)
	// checkError(err)
	// err = model.InsertFollow(3, 1)
	// checkError(err)
	// err = model.DeleteFollow(1, 3)
	// checkError(err)
	// err = model.DeleteFollow(1, 1)
	// checkError(err)
	// err = model.DeleteFollow(1, 2)
	// checkError(err)
	// err = model.DeleteFollow(1, 2)
	// checkError(err)

	// fmt.Println(model.QueryFollowersWithUserID(1))
	// fmt.Println(model.QueryFollowersWithUserID(2))
	// fmt.Println(model.QueryFollowingWithUserID(1))
	// fmt.Println(model.QueryFollowingWithUserID(2))

	//// ==========================================================================================================
	// err = model.InsertContentTag(1, "politics")
	// checkError(err)
	// err = model.InsertContentTag(1, "loser")
	// checkError(err)
	// err = model.InsertContentTag(1, "loser")
	// checkError(err)
	// err = model.InsertUserTag(1, "politics")
	// checkError(err)
	// err = model.InsertUserTag(1, "loser")
	// checkError(err)
	// err = model.InsertUserTag(1, "loser")
	// checkError(err)

	// fmt.Println(model.QueryTagsWithContentID(1))
	// fmt.Println(model.QueryTagsWithContentID(2))
	// fmt.Println(model.QueryTagsWithUserID(1))
	// fmt.Println(model.QueryTagsWithUserID(2))

}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
	}
}
