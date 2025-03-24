package usecase

import (
	"testing"

	"codepair-sinarmas/models"
	"codepair-sinarmas/service/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func Test_OTPUsecase_RequestOtp(t *testing.T) {
	type testCase struct {
		name             string
		wantError        bool
		expectedResponse *models.OTPResponse
		request          *models.OTPRequest
		onGetUserByID    func(mock *mocks.MockUserRepository)
		onGetOtpByUserID func(mock *mocks.MockOtpRepository)
		onSaveOTP        func(mock *mocks.MockOtpRepository)
	}

	var testTable []testCase
	testTable = append(testTable, testCase{
		name:      "success with no otp log",
		wantError: false,
		request: &models.OTPRequest{
			UserID: 1,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(1)).Return(&models.User{
				UserID: 1,
				Name:   "John Doe",
			}, nil)
		},
		onGetOtpByUserID: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().GetOtpByUserID(int64(1)).Return(&models.OTPLog{}, nil)
		},
		onSaveOTP: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().SaveOTP(gomock.Any()).Return(&models.OTPLog{
				UserID:           1,
				OTPCode:          "123456",
				NotificationType: models.NotificationTypeEmail,
				Status:           "created",
			}, nil)
		},
		expectedResponse: &models.OTPResponse{
			UserID: 1,
			OTP:    "123456",
		},
	})

	testTable = append(testTable, testCase{
		name:      "user not found",
		wantError: true,
		request: &models.OTPRequest{
			UserID: 2,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(2)).Return(nil, gorm.ErrRecordNotFound)
		},
		onGetOtpByUserID: nil,
		onSaveOTP:        nil,
		expectedResponse: nil,
	})

	testTable = append(testTable, testCase{
		name:      "failed to get user by ID",
		wantError: true,
		request: &models.OTPRequest{
			UserID: 3,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(3)).Return(nil, assert.AnError)
		},
		onGetOtpByUserID: nil,
		onSaveOTP:        nil,
		expectedResponse: nil,
	})

	testTable = append(testTable, testCase{
		name:      "otp already created",
		wantError: false,
		request: &models.OTPRequest{
			UserID: 4,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(4)).Return(&models.User{
				UserID: 4,
				Name:   "Jane Doe",
			}, nil)
		},
		onGetOtpByUserID: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().GetOtpByUserID(int64(4)).Return(&models.OTPLog{
				UserID:  4,
				OTPCode: "654321",
				Status:  "created",
			}, nil)
		},
		onSaveOTP: nil,
		expectedResponse: &models.OTPResponse{
			UserID: 4,
			OTP:    "654321",
		},
	})

	testTable = append(testTable, testCase{
		name:      "failed to get otp by user ID",
		wantError: true,
		request: &models.OTPRequest{
			UserID: 5,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(5)).Return(&models.User{
				UserID: 5,
				Name:   "Alice",
			}, nil)
		},
		onGetOtpByUserID: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().GetOtpByUserID(int64(5)).Return(nil, assert.AnError)
		},
		onSaveOTP:        nil,
		expectedResponse: nil,
	})

	testTable = append(testTable, testCase{
		name:      "failed to save otp",
		wantError: true,
		request: &models.OTPRequest{
			UserID: 6,
		},
		onGetUserByID: func(mock *mocks.MockUserRepository) {
			mock.EXPECT().GetUserByID(int64(6)).Return(&models.User{
				UserID: 6,
				Name:   "Bob",
			}, nil)
		},
		onGetOtpByUserID: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().GetOtpByUserID(int64(6)).Return(&models.OTPLog{}, nil)
		},
		onSaveOTP: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().SaveOTP(gomock.Any()).Return(nil, assert.AnError)
		},
		expectedResponse: nil,
	})

	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			otpRepo := mocks.NewMockOtpRepository(mockCtrl)
			userRepo := mocks.NewMockUserRepository(mockCtrl)

			if tc.onGetUserByID != nil {
				tc.onGetUserByID(userRepo)
			}

			if tc.onGetOtpByUserID != nil {
				tc.onGetOtpByUserID(otpRepo)
			}

			if tc.onSaveOTP != nil {
				tc.onSaveOTP(otpRepo)
			}

			usecase := &OtpUsecase{
				otpRepo:  otpRepo,
				userRepo: userRepo,
			}

			resp, err := usecase.RequestOtp(tc.request)

			if tc.wantError {
				assert.NotNil(t, err)
				assert.Nil(t, resp)
			} else {
				assert.Nil(t, err)
				assert.Equal(t, tc.expectedResponse, resp)
			}
		})
	}
}
