package postgres_test

import (
	"strconv"
	"testing"

	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"

	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
)

func TestSubscription_NoSettings(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "My new post", Description: "with this description"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(2)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newCommentSubscribers.Result[1].ID).Equals(aryaStark.ID)

	Expect(changeStatusSubscribers.Result).HasLen(2)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(changeStatusSubscribers.Result[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(nil)
	subscribed, err := users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(jonSnow)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_RemoveSubscriber(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	err = bus.Dispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: aryaStark})
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(1)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(1)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminSubmitted(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	newPostSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewPost}
	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventChangeStatus}
	err = bus.Dispatch(jonSnowCtx, newPostSubscribers, newCommentSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newPostSubscribers.Result).HasLen(1)
	Expect(newPostSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(newCommentSubscribers.Result).HasLen(1)
	Expect(newCommentSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(1)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_AdminUnsubscribed(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	bus.Dispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: aryaStark})
	bus.Dispatch(aryaStarkCtx, &cmd.RemoveSubscriber{Post: newPost.Result, User: jonSnow})

	newCommentSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	err = bus.Dispatch(aryaStarkCtx, newCommentSubscribers)
	Expect(err).IsNil()
	Expect(newCommentSubscribers.Result).HasLen(0)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_DisabledEmail(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err = users.UpdateSettings(map[string]string{
		models.NotificationEventNewComment.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	newCommentWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	newCommentEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelEmail, Event: models.NotificationEventNewComment}
	changeStatusSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newCommentWebSubscribers, newCommentEmailSubscribers, changeStatusSubscribers)
	Expect(err).IsNil()

	Expect(newCommentWebSubscribers.Result).HasLen(2)
	Expect(newCommentWebSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newCommentWebSubscribers.Result[1].ID).Equals(aryaStark.ID)

	Expect(newCommentEmailSubscribers.Result).HasLen(1)
	Expect(newCommentEmailSubscribers.Result[0].ID).Equals(jonSnow.ID)

	Expect(changeStatusSubscribers.Result).HasLen(2)
	Expect(changeStatusSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(changeStatusSubscribers.Result[1].ID).Equals(aryaStark.ID)

	users.SetCurrentUser(jonSnow)
	subscribed, err := users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(aryaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsTrue()

	users.SetCurrentUser(sansaStark)
	subscribed, err = users.HasSubscribedTo(newPost.Result.ID)
	Expect(err).IsNil()
	Expect(subscribed).IsFalse()
}

func TestSubscription_VisitorEnabledNewPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	users.SetCurrentTenant(demoTenant)
	users.SetCurrentUser(aryaStark)

	err = users.UpdateSettings(map[string]string{
		models.NotificationEventNewPost.UserSettingsKeyName: strconv.Itoa(int(models.NotificationChannelEmail | models.NotificationChannelWeb)),
	})
	Expect(err).IsNil()

	newPostWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewPost}
	newPostEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelEmail, Event: models.NotificationEventNewPost}
	err = bus.Dispatch(aryaStarkCtx, newPostWebSubscribers, newPostEmailSubscribers)
	Expect(err).IsNil()

	Expect(newPostWebSubscribers.Result).HasLen(2)
	Expect(newPostWebSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newPostWebSubscribers.Result[1].ID).Equals(aryaStark.ID)

	Expect(newPostEmailSubscribers.Result).HasLen(2)
	Expect(newPostEmailSubscribers.Result[0].ID).Equals(jonSnow.ID)
	Expect(newPostEmailSubscribers.Result[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_DisabledEverything(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(jonSnowCtx, newPost)
	Expect(err).IsNil()

	users.SetCurrentTenant(demoTenant)
	disableAll := map[string]string{
		models.NotificationEventNewPost.UserSettingsKeyName:      "0",
		models.NotificationEventNewComment.UserSettingsKeyName:   "0",
		models.NotificationEventChangeStatus.UserSettingsKeyName: "0",
	}
	users.SetCurrentUser(jonSnow)
	Expect(users.UpdateSettings(disableAll)).IsNil()
	users.SetCurrentUser(aryaStark)
	Expect(users.UpdateSettings(disableAll)).IsNil()

	newPostWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewPost}
	newPostEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelEmail, Event: models.NotificationEventNewPost}
	newCommentWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	newCommentEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelEmail, Event: models.NotificationEventNewComment}
	changeStatusWebSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventChangeStatus}
	changeStatusEmailSubscribers := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelEmail, Event: models.NotificationEventChangeStatus}
	err = bus.Dispatch(aryaStarkCtx, newPostWebSubscribers, newPostEmailSubscribers, newCommentWebSubscribers, newCommentEmailSubscribers, changeStatusWebSubscribers, changeStatusEmailSubscribers)
	Expect(err).IsNil()

	Expect(newPostWebSubscribers.Result).HasLen(0)
	Expect(newPostEmailSubscribers.Result).HasLen(0)
	Expect(newCommentWebSubscribers.Result).HasLen(0)
	Expect(newCommentEmailSubscribers.Result).HasLen(0)
	Expect(changeStatusWebSubscribers.Result).HasLen(0)
	Expect(changeStatusEmailSubscribers.Result).HasLen(0)
}

func TestSubscription_DeletedPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	err := bus.Dispatch(aryaStarkCtx, newPost)
	Expect(err).IsNil()

	bus.Dispatch(aryaStarkCtx, &cmd.SetPostResponse{Post: newPost.Result, Text: "Invalid Post!", Status: models.PostDeleted})

	q := &query.GetActiveSubscribers{Number: newPost.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	err = bus.Dispatch(aryaStarkCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(2)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
	Expect(q.Result[1].ID).Equals(aryaStark.ID)
}

func TestSubscription_SubscribedToDifferentPost(t *testing.T) {
	SetupDatabaseTest(t)
	defer TeardownDatabaseTest()

	newPost1 := &cmd.AddNewPost{Title: "Post #1", Description: "Description #1"}
	newPost2 := &cmd.AddNewPost{Title: "Post #2", Description: "Description #2"}
	err := bus.Dispatch(jonSnowCtx, newPost1, newPost2)
	Expect(err).IsNil()

	err = bus.Dispatch(jonSnowCtx, &cmd.AddSubscriber{Post: newPost2.Result, User: aryaStark})
	Expect(err).IsNil()

	q := &query.GetActiveSubscribers{Number: newPost1.Result.Number, Channel: models.NotificationChannelWeb, Event: models.NotificationEventNewComment}
	err = bus.Dispatch(jonSnowCtx, q)
	Expect(err).IsNil()
	Expect(q.Result).HasLen(1)
	Expect(q.Result[0].ID).Equals(jonSnow.ID)
}
