package store

import (
	"context"
	"fmt"
	"github.com/vitaliy-paliy/amica/pkg/models"
	"github.com/vitaliy-paliy/amica/pkg/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type UserStore struct {
	client *mongo.Client
}

func NewUserStore(client *mongo.Client) *UserStore {
	return &UserStore{client: client}
}

func (us *UserStore) Find(user *models.User) *models.CustomError {
	selectedFields := bson.M{"_id": 1, "username": 1, "phone_number": 1} // Return only selected fields.
	opts := options.FindOne().SetProjection(selectedFields)
	err := utils.Collection(us.client, "users").FindOne(context.TODO(), bson.D{{"_id", user.ID}}, opts).Decode(user)
	if err != nil {
		return &models.CustomError{Code: 404, Message: fmt.Sprintf("User with ID(%s) was not found.", user.ID.Hex())}
	}

	return nil
}

func (us *UserStore) IsRequestPending(request *models.FriendRequest) error {
	return utils.Collection(us.client, "friend-requests").FindOne(context.TODO(), request).Decode(request)
}

func (us *UserStore) ValidateRequest(request *models.FriendRequest) *models.CustomError {
	if request.SelfRequest() {
		return &models.CustomError{Code: 400, Message: "Invalid friend request."}
	}

	query := bson.M{"_id": request.Sender.ID, "friends": bson.M{"$in": []*models.User{request.Recipient}}}
	if err := utils.Collection(us.client, "users").FindOne(context.TODO(), query).Decode(&bson.D{}); err == nil {
		return &models.CustomError{Code: 409, Message: "Users are already friends."}
	}

	return nil
}

func (us *UserStore) SendFriendRequest(request *models.FriendRequest) (*models.FriendRequest, *models.CustomError) {
	if err := us.ValidateRequest(request); err != nil {
		return nil, err
	}

	invertedRequest := request.Invert()
	if err := us.IsRequestPending(invertedRequest); err == nil {
		return us.AcceptFriendRequest(invertedRequest)
	}

	request.Timestamp = time.Now()
	opts := options.Update().SetUpsert(true) // If a given friend request exists, update it.
	update := bson.D{{"$set", request}}

	res, err := utils.Collection(us.client, "friend-requests").UpdateOne(context.TODO(), request.ToUserDoc(), update, opts)
	if err != nil {
		return nil, &models.CustomError{Code: 500, Message: "An error occured."}
	}

	if res.UpsertedID == nil {
		return nil, &models.CustomError{Code: 400, Message: "This user is already present in your sent friend requests"}
	}

	request.ID = res.UpsertedID.(primitive.ObjectID)

	return request, nil
}

func (us *UserStore) AcceptFriendRequest(request *models.FriendRequest) (*models.FriendRequest, *models.CustomError) {
	if err := us.DeleteFriendRequest(request); err != nil {
		return nil, err
	}

	return request, us.bulkWriteFriendRequest("$addToSet", request)
}

func (us *UserStore) DeleteFriendRequest(request *models.FriendRequest) *models.CustomError {
	_, err := utils.Collection(us.client, "friend-requests").DeleteOne(context.TODO(), request)
	if err != nil {
		return &models.CustomError{Code: 500, Message: "An error occured."}
	}

	return nil
}

func (us *UserStore) RemoveFriend(request *models.FriendRequest) *models.CustomError {
	return us.bulkWriteFriendRequest("$pull", request)
}

func (us *UserStore) bulkWriteFriendRequest(operator string, request *models.FriendRequest) *models.CustomError {
	writeModels := make([]mongo.WriteModel, 2)
	users := *request.ToSlice()
	for idx, user := range users {
		model := bson.D{{operator, bson.M{"friends": user}}}
		writeModels[idx] = mongo.NewUpdateOneModel().SetFilter(bson.D{{"_id", users[idx^1].ID}}).SetUpdate(model).SetUpsert(true)
	}
	opts := options.BulkWrite().SetOrdered(false)

	_, err := utils.Collection(us.client, "users").BulkWrite(context.TODO(), writeModels, opts)
	if err != nil {
		return &models.CustomError{Code: 500, Message: "An error occured."}
	}

	return nil
}
