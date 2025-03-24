package usecase

import (
	"testing"

	"codepair-sinarmas/models"
	"codepair-sinarmas/service/repository/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
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
			mock.EXPECT().GetUserByID(1).Return(&models.User{
				UserID: 1,
				Name:   "John Doe",
			}, nil)
		},
		onGetOtpByUserID: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().GetOtpByUserID(1).Return(&models.OTPLog{}, nil)
		},
		onSaveOTP: func(mock *mocks.MockOtpRepository) {
			mock.EXPECT().SaveOTP(&models.OTPLog{
				UserID:  1,
				OTPCode: "123456",
				Status:  "created",
			}).Return(&models.OTPLog{
				UserID:  1,
				OTPCode: "123456",
				Status:  "created",
			}, nil)
		},
		expectedResponse: &models.OTPResponse{
			UserID: 1,
			OTP:    "123456",
		},
	})

	// testTable = append(testTable, testCase{
	// 	name:      "email already registered",
	// 	wantError: true,
	// 	request: &models.RegisterUser{
	// 		Name:     "Jane Doe",
	// 		Email:    "jane@example.com",
	// 		Password: "password123",
	// 	},
	// 	onGetUserByEmail: func(mock *mocks.MockUserRepository) {
	// 		mock.EXPECT().GetUserByEmail("jane@example.com").Return(&models.User{UserID: 2}, nil)
	// 	},
	// 	expectedResponse: nil,
	// })

	// testTable = append(testTable, testCase{
	// 	name:      "error checking user by email",
	// 	wantError: true,
	// 	request: &models.RegisterUser{
	// 		Name:     "Error User",
	// 		Email:    "error@example.com",
	// 		Password: "password123",
	// 	},
	// 	onGetUserByEmail: func(mock *mocks.MockUserRepository) {
	// 		mock.EXPECT().GetUserByEmail("error@example.com").Return(nil, errors.New("database error"))
	// 	},
	// 	expectedResponse: nil,
	// })

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
