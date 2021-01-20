# IOS_Final_Backend

[YourTube](https://github.com/chenguofan1999/YourTube) 的后端，已运行于http://159.75.1.231:5009

**API:**

```cpp
用户服务
    注册与登录
    POST   /signup           注册
    POST   /login            登录

    用户信息
    GET    /users/{username} 获取某用户详细信息
    GET    /user             获取自己的用户信息
    PUT    /user/bio         更新自己的简介
    PUT    /user/avatar      更新自己的头像
    POST   /user/avatar      更新自己的头像（为了适配AFNetworking增加的冗余接口）
    GET    /user/tags        为自己增加关注的 tag
    POST   /user/tags        为自己增加关注的 tag
    DELETE /user/tags        为自己删除关注的 tag

    关注
    GET    /users/{username}/followers   获取某用户关注者
    GET    /users/{username}/following   获取某用户关注的人
    GET    /user/following/{username}    检查是否已关注某用户
    PUT    /user/following/{username}    关注某用户
    DELETE /user/following/{username}    取消关注某用户

Content 服务
    GET    /contents              获取内容集(query参数在下方)
    POST   /contents              发布内容
    GET    /contents/{contentID}  获取某条内容的详细信息
    DELETE /contents/{contentID}  删除内容，仅能删除自己发出的内容
    POST   /content/tags          为内容增加 tag , 仅能为自己发的内容增加标签
    DELETE /content/tags          为内容删除 tag , 仅能为自己发的内容删除标签

Comment 服务
    GET    /comments             获取评论集(query参数在下方)
    POST   /comments             发布评论
    DELETE /comments/{commentID} 删除评论, 仅能删除自己发的评论

Reply 服务
    GET    /replies           获取回复集(query参数在下方)
    POST   /replies           发布回复
    DELETE /replies/{replyID} 删除回复, 仅能删除自己发的回复

Like
    PUT    /like/content/{contentID}   喜欢某条内容
    DELETE /like/content/{contentID}   取消喜欢某条内容
    PUT    /like/comment/{commentID}   喜欢某条评论
    DELETE /like/comment/{commentID}   取消喜欢某条评论
    PUT    /like/reply/{replyID}       喜欢某条回复
    DELETE /like/reply/{replyID}       取消喜欢某条回复


// GET /contents 查询多条内容, 由 query 参数决定查询的模式:
// 这些参数是互斥的，即一条请求中只能 query 其中之一：
// 1. tag : 获取带有某个标签的内容 (tag = {tagName})
// 2. user : 获取指定用户的内容 (user = {userName})
// 3. search : 搜索标题含特定字符串的内容 (tag = {searchStr})
// 4. follow : 当前用户关注的所有用户的内容 (follow = true)
// 5. self : 当前用户自己发的内容 (self = true)
// 6. history : 获取自己的观看记录 (history = true)
// 7. allTags : 获取自己关注的全部tag的内容 (allTags = true)
// 8. likedBy : 获取某人喜欢的全部内容 (likedBy = {username})
// 9. 如果以上参数都没有，则为请求不经过筛选的公共内容
// 以下参数与上面的参数兼容
// 1. orderBy : viewNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
// 3. num : 指定条数, 默认 30

// GET /comments 查询评论集, 必要参数:
// 1. contentID: 获取某个 content 的所有评论
// 可选参数:
// 1. orderBy : likeNum / time ，默认 time
// 2. order : asc / desc ，默认 desc

// GET /replies 查询回复集, 必要参数:
// 1. commentID: 获取某个 comment 的所有回复
// 可选参数:
// 1. orderBy : likeNum / time ，默认 time
// 2. order : asc / desc ，默认 desc
```

