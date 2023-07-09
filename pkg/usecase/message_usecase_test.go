package usecase

import (
	"context"
	"testing"

	"github.com/shunsukenagashima/chat-api/pkg/apperror"
	"github.com/shunsukenagashima/chat-api/pkg/clock"
	"github.com/shunsukenagashima/chat-api/pkg/domain/model"
	"github.com/shunsukenagashima/chat-api/pkg/domain/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAllMessagesByRoomID(t *testing.T) {
	mockMessageRepo := new(mocks.MessageRepository)
	clock := clock.FixedClocker{}
	mockMessages := []*model.Message{
		{
			MessageID: "1",
			RoomID:    "1",
			UserID:    "1",
			Content:   "Hello",
			CreatedAt: clock.Now(),
		},
		{
			MessageID: "2",
			RoomID:    "1",
			UserID:    "1",
			Content:   "World",
			CreatedAt: clock.Now(),
		},
		{
			MessageID: "3",
			RoomID:    "1",
			UserID:    "2",
			Content:   "Hello World",
			CreatedAt: clock.Now(),
		},
	}

	mockMessageRepo.On("GetMessagesByRoomID", mock.Anything, mockMessages[0].RoomID, mock.Anything, mock.Anything).Return(mockMessages, "4", nil)
	messageUsecase := NewMessageUsecase(mockMessageRepo)

	messages, _, err := messageUsecase.GetMessagesByRoomID(context.Background(), "1", "0", 10)

	assert.NoError(t, err)
	assert.NotEmpty(t, messages)
	assert.Equal(t, len(mockMessages), len(messages))
	for i, message := range messages {
		assert.Equal(t, mockMessages[i].MessageID, message.MessageID)
		assert.Equal(t, mockMessages[i].RoomID, message.RoomID)
		assert.Equal(t, mockMessages[i].UserID, message.UserID)
		assert.Equal(t, mockMessages[i].Content, message.Content)
		assert.Equal(t, mockMessages[i].CreatedAt, message.CreatedAt)
	}
	mockMessageRepo.AssertExpectations(t)
}

func TestCreateMessage(t *testing.T) {
	mockMessageRepo := new(mocks.MessageRepository)
	mockMessage := &model.Message{
		RoomID:  "1",
		UserID:  "1",
		Content: "Hello",
	}

	mockMessageRepo.On("Create", mock.Anything, mockMessage).Return(nil)
	messageUsecase := NewMessageUsecase(mockMessageRepo)

	err := messageUsecase.CreateMessage(context.Background(), mockMessage)

	assert.NoError(t, err)
	mockMessageRepo.AssertExpectations(t)
}

func TestUpdateMessage(t *testing.T) {
	clock := clock.FixedClocker{}
	mockMessage := &model.Message{
		MessageID: "1",
		RoomID:    "1",
		UserID:    "1",
		Content:   "Hello",
		CreatedAt: clock.Now(),
	}

	testCases := []struct {
		name          string
		roomId        string
		messageId     string
		newContent    string
		getByIdReturn *model.Message
		expectedErr   error
	}{
		{
			name:          "Success",
			roomId:        mockMessage.RoomID,
			messageId:     mockMessage.MessageID,
			newContent:    "Hello World",
			getByIdReturn: mockMessage,
			expectedErr:   nil,
		},
		{
			name:          "Invalid MessageID",
			roomId:        mockMessage.RoomID,
			messageId:     "2",
			newContent:    "Hello World",
			getByIdReturn: nil,
			expectedErr:   apperror.NewNotFoundErr("Message", "MessageID: 2"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockMessageRepo := new(mocks.MessageRepository)
			mockMessageRepo.On("GetByID", mock.Anything, tc.roomId, tc.messageId).Return(tc.getByIdReturn, tc.expectedErr)

			mockMessageRepo.On("Update", mock.Anything, tc.roomId, tc.messageId, tc.newContent).Return(nil)
			messageUsecase := NewMessageUsecase(mockMessageRepo)

			err := messageUsecase.UpdateMessage(context.Background(), tc.roomId, tc.messageId, tc.newContent)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				mockMessageRepo.AssertExpectations(t)
			}
		})
	}

}

func TestDeleteMessage(t *testing.T) {
	clock := clock.FixedClocker{}
	mockMessage := &model.Message{
		MessageID: "1",
		RoomID:    "1",
		UserID:    "1",
		Content:   "Hello",
		CreatedAt: clock.Now(),
	}

	testCases := []struct {
		name          string
		roomId        string
		messageId     string
		getByIdReturn *model.Message
		expectedErr   error
	}{
		{
			name:          "Success",
			roomId:        mockMessage.RoomID,
			messageId:     mockMessage.MessageID,
			getByIdReturn: mockMessage,
			expectedErr:   nil,
		},
		{
			name:          "Invalid MessageID",
			roomId:        mockMessage.RoomID,
			messageId:     "2",
			getByIdReturn: nil,
			expectedErr:   apperror.NewNotFoundErr("Message", "MessageID: 2"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockMessageRepo := new(mocks.MessageRepository)
			mockMessageRepo.On("GetByID", mock.Anything, tc.roomId, tc.messageId).Return(tc.getByIdReturn, tc.expectedErr)

			mockMessageRepo.On("Delete", mock.Anything, tc.roomId, tc.messageId).Return(nil)
			messageUsecase := NewMessageUsecase(mockMessageRepo)

			err := messageUsecase.DeleteMessage(context.Background(), tc.roomId, tc.messageId)

			if tc.expectedErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedErr, err)
			} else {
				assert.NoError(t, err)
				mockMessageRepo.AssertExpectations(t)

			}
		})
	}
}
