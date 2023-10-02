package usecase

import (
	"context"
	"database/sql"
	"errors"
	"github.com/h2p2f/dedicated-vault/internal/client/config"
	"github.com/h2p2f/dedicated-vault/internal/client/models"
	"github.com/h2p2f/dedicated-vault/internal/client/usecase/mocks"
	pb "github.com/h2p2f/dedicated-vault/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestClientUseCase_CreateUser(t *testing.T) {

	tests := []struct {
		name                     string
		userName                 string
		password                 string
		passphrase               string
		createUserError          error
		registerToken            string
		registerError            error
		updateLastServerErr      error
		expectedConfigPassphrase string
		expectedConfigToken      string
		expectedConfigUser       string
		expectedStorageErr       error
		expectedTransporterErr   error
		expectedErr              error
	}{
		{
			name:                     "Successful create user",
			userName:                 "testuser",
			password:                 "testpassword",
			passphrase:               "testpassphrase",
			createUserError:          nil,
			registerToken:            "testtoken",
			registerError:            nil,
			updateLastServerErr:      nil,
			expectedConfigPassphrase: "testpassphrase",
			expectedConfigToken:      "testtoken",
			expectedConfigUser:       "testuser",
			expectedStorageErr:       nil,
			expectedTransporterErr:   nil,
			expectedErr:              nil,
		},
		{
			name:                     "Error creating user in storage",
			userName:                 "testuser",
			password:                 "testpassword",
			passphrase:               "testpassphrase",
			createUserError:          errors.New("storage error"),
			registerToken:            "",
			registerError:            nil,
			updateLastServerErr:      nil,
			expectedConfigPassphrase: "",
			expectedConfigToken:      "",
			expectedConfigUser:       "",
			expectedStorageErr:       errors.New("storage error"),
			expectedTransporterErr:   nil,
			expectedErr:              errors.New("storage error"),
		},
		{
			name:                     "Error registering user with transporter",
			userName:                 "testuser",
			password:                 "testpassword",
			passphrase:               "testpassphrase",
			createUserError:          nil,
			registerToken:            "",
			registerError:            errors.New("transporter error"),
			updateLastServerErr:      nil,
			expectedConfigPassphrase: "",
			expectedConfigToken:      "",
			expectedConfigUser:       "",
			expectedStorageErr:       nil,
			expectedTransporterErr:   errors.New("transporter error"),
			expectedErr:              errors.New("transporter error"),
		},
		{
			name:                     "Error updating last server updated time",
			userName:                 "testuser",
			password:                 "testpassword",
			passphrase:               "testpassphrase",
			createUserError:          nil,
			registerToken:            "testtoken",
			registerError:            nil,
			updateLastServerErr:      errors.New("storage error"),
			expectedConfigPassphrase: "testpassphrase",
			expectedConfigToken:      "testtoken",
			expectedConfigUser:       "testuser",
			expectedStorageErr:       errors.New("storage error"),
			expectedTransporterErr:   nil,
			expectedErr:              errors.New("storage error"),
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock objects
			mockStorage := mocks.NewStorager(t)
			mockStorage.On("CreateUser", tt.userName).Return(tt.createUserError)
			if tt.createUserError == nil && tt.expectedTransporterErr == nil {
				mockStorage.On("UpdateLastServerUpdated", tt.userName, int64(0)).Return(tt.updateLastServerErr)
			}
			mockTransport := mocks.NewTransporter(t)
			if tt.createUserError == nil {
				mockTransport.On("Register", context.Background(), &pb.User{
					Name:     tt.userName,
					Password: tt.password,
				}).Return(tt.registerToken, tt.registerError)
			}
			testConfig := config.NewClientConfig()
			// Create client use case
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}

			// Call function
			err := clientUseCase.CreateUser(context.Background(), tt.userName, tt.password, tt.passphrase)

			// Check output and error
			assert.Equal(t, tt.expectedConfigPassphrase, clientUseCase.Config.Passphrase)
			assert.Equal(t, tt.expectedConfigToken, clientUseCase.Config.Token)
			assert.Equal(t, tt.expectedConfigUser, clientUseCase.Config.User)
			assert.Equal(t, tt.expectedErr, err)
			if tt.expectedStorageErr == nil && tt.updateLastServerErr == nil {
				_, transportError := clientUseCase.Transporter.Register(context.Background(), &pb.User{
					Name:     tt.userName,
					Password: tt.password,
				})
				assert.Equal(t, tt.expectedTransporterErr, transportError)
			}
		})
	}
}

func TestClientUseCase_LoginUser(t *testing.T) {

	tests := []struct {
		name                    string
		userName                string
		password                string
		passphrase              string
		getUserIDError          error
		createUserError         error
		loginToken              string
		loginError              error
		configLastServerUpdated int64
		getLastServerUpdated    int64
		fullSyncError           error
		updateLastServerErr     error
		expectedConfigPass      string
		expectedConfigToken     string
		expectedConfigUser      string
		expectedStorageErr      error
		expectedTransporterErr  error
	}{
		{
			name:                    "Successful login user with blank user ID",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          sql.ErrNoRows,
			createUserError:         nil,
			loginToken:              "testtoken",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     nil,
			expectedConfigPass:      "testpassphrase",
			expectedConfigToken:     "testtoken",
			expectedConfigUser:      "testuser",
			expectedStorageErr:      nil,
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Successful login user with existing user ID",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          nil,
			createUserError:         nil,
			loginToken:              "testtoken",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     nil,
			expectedConfigPass:      "testpassphrase",
			expectedConfigToken:     "testtoken",
			expectedConfigUser:      "testuser",
			expectedStorageErr:      nil,
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Error logging in user with transporter",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          nil,
			createUserError:         nil,
			loginToken:              "",
			loginError:              errors.New("transporter error"),
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     nil,
			expectedConfigPass:      "",
			expectedConfigToken:     "",
			expectedConfigUser:      "",
			expectedStorageErr:      errors.New("transporter error"),
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Error creating user in storage",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          errors.New("storage error"),
			createUserError:         errors.New("storage error"),
			loginToken:              "",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     nil,
			expectedConfigPass:      "",
			expectedConfigToken:     "",
			expectedConfigUser:      "",
			expectedStorageErr:      errors.New("storage error"),
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Error getting last server updated time from storage",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          sql.ErrNoRows,
			createUserError:         nil,
			loginToken:              "testtoken",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     errors.New("storage error"),
			expectedConfigPass:      "testpassphrase",
			expectedConfigToken:     "testtoken",
			expectedConfigUser:      "testuser",
			expectedStorageErr:      errors.New("storage error"),
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Error during full sync",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          sql.ErrNoRows,
			createUserError:         nil,
			loginToken:              "testtoken",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    0,
			fullSyncError:           errors.New("sync error"),
			updateLastServerErr:     nil,
			expectedConfigPass:      "testpassphrase",
			expectedConfigToken:     "testtoken",
			expectedConfigUser:      "testuser",
			expectedStorageErr:      errors.New("sync error"),
			expectedTransporterErr:  nil,
		},
		{
			name:                    "Error updating last server updated time",
			userName:                "testuser",
			password:                "testpassword",
			passphrase:              "testpassphrase",
			getUserIDError:          nil,
			createUserError:         nil,
			loginToken:              "testtoken",
			loginError:              nil,
			configLastServerUpdated: 10,
			getLastServerUpdated:    10,
			fullSyncError:           nil,
			updateLastServerErr:     errors.New("storage error"),
			expectedConfigPass:      "testpassphrase",
			expectedConfigToken:     "testtoken",
			expectedConfigUser:      "testuser",
			expectedStorageErr:      errors.New("storage error"),
			expectedTransporterErr:  nil,
		},
	}

	// Run test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock objects
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.LastServerUpdated = tt.configLastServerUpdated

			mockTransport.On("Login", context.Background(), &pb.User{
				Name:     tt.userName,
				Password: tt.password,
			}).Return(tt.loginToken, tt.loginError)

			if tt.loginError == nil {
				mockStorage.On("GetUserID", tt.userName).Return(int64(0), tt.getUserIDError)
				if tt.getUserIDError != nil {
					if errors.Is(tt.getUserIDError, sql.ErrNoRows) {
						mockStorage.On("CreateUser", tt.userName).Return(tt.createUserError)
					}
				}
				if tt.createUserError == nil {
					mockStorage.On("GetLastServerUpdated", tt.userName).Return(tt.getLastServerUpdated, nil)
					if tt.getLastServerUpdated < clientUseCase.Config.LastServerUpdated {
						mockStorage.On("DeleteAllData", tt.userName).Return(nil)
						var protoData []*pb.SecretData
						mockTransport.On("ListSecrets", context.Background()).Return(protoData, tt.fullSyncError)
						if tt.fullSyncError == nil {
							mockStorage.On("UpdateLastServerUpdated", tt.userName, int64(0)).Return(tt.updateLastServerErr)
						}
					}
					if tt.fullSyncError == nil {
						mockStorage.On("UpdateLastServerUpdated", tt.userName, int64(10)).Return(tt.updateLastServerErr)
					}
				}
			}
			// Call function
			err := clientUseCase.LoginUser(context.Background(), tt.userName, tt.password, tt.passphrase)

			// Check output and error
			assert.Equal(t, tt.expectedConfigPass, clientUseCase.Config.Passphrase)
			assert.Equal(t, tt.expectedConfigToken, clientUseCase.Config.Token)
			assert.Equal(t, tt.expectedConfigUser, clientUseCase.Config.User)
			assert.Equal(t, tt.expectedStorageErr, err)
		})
	}
}

func TestClientUseCase_ChangePassword(t *testing.T) {
	tests := []struct {
		name                string
		userName            string
		password            string
		newPassword         string
		userID              int64
		getUserIDError      error
		loginError          error
		expectedConfigUser  string
		expectedConfigToken string
		expectedError       error
	}{
		{
			name:                "Successful change password",
			userName:            "testuser",
			password:            "testpassword",
			newPassword:         "newtestpassword",
			userID:              1,
			getUserIDError:      nil,
			loginError:          nil,
			expectedConfigUser:  "testuser",
			expectedConfigToken: "testtoken",
			expectedError:       nil,
		},
		{
			name:                "Error getting user ID",
			userName:            "testuser",
			password:            "testpassword",
			newPassword:         "newtestpassword",
			userID:              0,
			getUserIDError:      errors.New("storage error"),
			loginError:          nil,
			expectedConfigUser:  "",
			expectedConfigToken: "",
			expectedError:       errors.New("storage error"),
		},
		{
			name:                "Error changing password with transporter",
			userName:            "testuser",
			password:            "testpassword",
			newPassword:         "newtestpassword",
			userID:              1,
			getUserIDError:      nil,
			loginError:          errors.New("transporter error"),
			expectedConfigUser:  "",
			expectedConfigToken: "",
			expectedError:       errors.New("transporter error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}

			mockStorage.On("GetUserID", tt.userName).Return(tt.userID, tt.getUserIDError)
			if tt.getUserIDError == nil {
				mockTransport.On("ChangePassword", context.Background(), &pb.User{
					Name:     tt.userName,
					Password: tt.password,
				}, tt.newPassword).Return(tt.expectedConfigToken, tt.loginError)
			}
			err := clientUseCase.ChangePassword(context.Background(), tt.userName, tt.password, tt.newPassword)

			assert.Equal(t, tt.expectedError, err)
			assert.Equal(t, tt.expectedConfigUser, clientUseCase.Config.User)
			assert.Equal(t, tt.expectedConfigToken, clientUseCase.Config.Token)

		})
	}
}

func TestClientUseCase_SaveData(t *testing.T) {

	data := models.Data{
		UUID:     "testuuid",
		Meta:     "testmeta",
		DataType: "testdatatype",
		Folder: models.Folder{
			Text: models.TextData{
				Text: "",
			},
		},
	}

	tests := []struct {
		name                  string
		user                  string
		token                 string
		createError           error
		saveSecretError       error
		updateLastServerError error
		expectedError         error
	}{
		{
			name:                  "Successful save data",
			user:                  "testuser",
			token:                 "testtoken",
			createError:           nil,
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         nil,
		},
		{
			name:                  "Error creating data in storage",
			user:                  "testuser",
			token:                 "testtoken",
			createError:           errors.New("storage error"),
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error saving secret with transporter",
			user:                  "testuser",
			token:                 "testtoken",
			createError:           nil,
			saveSecretError:       errors.New("transporter error"),
			updateLastServerError: nil,
			expectedError:         errors.New("transporter error"),
		},
		{
			name:                  "Error updating last server updated time",
			user:                  "testuser",
			token:                 "testtoken",
			createError:           nil,
			saveSecretError:       nil,
			updateLastServerError: errors.New("storage error"),
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error user not logged in",
			user:                  "",
			token:                 "",
			createError:           nil,
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         errors.New("user not logged in"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.Token = tt.token
			if tt.token != "" {
				clientUseCase.Config.User = tt.user
				clientUseCase.Config.Passphrase = "testpassphrase"
				mockStorage.On("CreateData", tt.user, mock.Anything).Return(tt.createError)
				if tt.createError == nil {
					mockTransport.On("SaveSecret", context.Background(), mock.Anything).Return(tt.saveSecretError)
					if tt.saveSecretError == nil {
						mockStorage.On("UpdateLastServerUpdated", tt.user, int64(0)).Return(tt.updateLastServerError)
					}
				}
			}
			err = clientUseCase.SaveData(context.Background(), data)

			assert.Equal(t, tt.expectedError, err)

		})
	}
}

func TestClientUseCase_ChangeData(t *testing.T) {

	data := models.Data{
		UUID:     "testuuid",
		Meta:     "testmeta",
		DataType: "testdatatype",
		Folder: models.Folder{
			Text: models.TextData{
				Text: "",
			},
		},
	}

	tests := []struct {
		name                  string
		user                  string
		token                 string
		changeError           error
		saveSecretError       error
		updateLastServerError error
		expectedError         error
	}{
		{
			name:                  "Successful change data",
			user:                  "testuser",
			token:                 "testtoken",
			changeError:           nil,
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         nil,
		},
		{
			name:                  "Error changing data in storage",
			user:                  "testuser",
			token:                 "testtoken",
			changeError:           errors.New("storage error"),
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error saving secret with transporter",
			user:                  "testuser",
			token:                 "testtoken",
			changeError:           nil,
			saveSecretError:       errors.New("transporter error"),
			updateLastServerError: nil,
			expectedError:         errors.New("transporter error"),
		},
		{
			name:                  "Error updating last server updated time",
			user:                  "testuser",
			token:                 "testtoken",
			changeError:           nil,
			saveSecretError:       nil,
			updateLastServerError: errors.New("storage error"),
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error user not logged in",
			user:                  "",
			token:                 "",
			changeError:           nil,
			saveSecretError:       nil,
			updateLastServerError: nil,
			expectedError:         errors.New("user not logged in"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.Token = tt.token
			if tt.token != "" {
				clientUseCase.Config.User = tt.user
				clientUseCase.Config.Passphrase = "testpassphrase"
				mockStorage.On("UpdateData", tt.user, mock.Anything).Return(tt.changeError)
				if tt.changeError == nil {
					mockTransport.On("ChangeSecret", context.Background(), mock.Anything).Return(tt.saveSecretError)
					if tt.saveSecretError == nil {
						mockStorage.On("UpdateLastServerUpdated", tt.user, int64(0)).Return(tt.updateLastServerError)
					}
				}
			}
			err = clientUseCase.ChangeData(context.Background(), data)

			assert.Equal(t, tt.expectedError, err)

		})
	}
}

func TestClientUseCase_DeleteData(t *testing.T) {

	data := models.Data{
		UUID:     "testuuid",
		Meta:     "testmeta",
		DataType: "testdatatype",
		Folder: models.Folder{
			Text: models.TextData{
				Text: "",
			},
		},
	}

	tests := []struct {
		name                  string
		user                  string
		token                 string
		deleteError           error
		deleteSecretError     error
		updateLastServerError error
		expectedError         error
	}{
		{
			name:                  "Successful delete data",
			user:                  "testuser",
			token:                 "testtoken",
			deleteError:           nil,
			deleteSecretError:     nil,
			updateLastServerError: nil,
			expectedError:         nil,
		},
		{
			name:                  "Error deleting data in storage",
			user:                  "testuser",
			token:                 "testtoken",
			deleteError:           errors.New("storage error"),
			deleteSecretError:     nil,
			updateLastServerError: nil,
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error deleting secret with transporter",
			user:                  "testuser",
			token:                 "testtoken",
			deleteError:           nil,
			deleteSecretError:     errors.New("transporter error"),
			updateLastServerError: nil,
			expectedError:         errors.New("transporter error"),
		},
		{
			name:                  "Error updating last server updated time",
			user:                  "testuser",
			token:                 "testtoken",
			deleteError:           nil,
			deleteSecretError:     nil,
			updateLastServerError: errors.New("storage error"),
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error user not logged in",
			user:                  "",
			token:                 "",
			deleteError:           nil,
			deleteSecretError:     nil,
			updateLastServerError: nil,
			expectedError:         errors.New("user not logged in"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.Token = tt.token
			if tt.token != "" {
				clientUseCase.Config.User = tt.user
				clientUseCase.Config.Passphrase = "testpassphrase"
				mockStorage.On("DeleteData", tt.user, mock.Anything).Return(tt.deleteError)
				if tt.deleteError == nil {
					mockTransport.On("DeleteSecret", context.Background(), mock.Anything).Return(tt.deleteSecretError)
					if tt.deleteSecretError == nil {
						mockStorage.On("UpdateLastServerUpdated", tt.user, int64(0)).Return(tt.updateLastServerError)
					}
				}
			}
			err = clientUseCase.DeleteData(context.Background(), data)

			assert.Equal(t, tt.expectedError, err)

		})
	}
}

func TestClientUseCase_GetDataByType(t *testing.T) {
	tests := []struct {
		name          string
		user          string
		token         string
		dataType      string
		getDataError  error
		expectedError error
	}{
		{
			name:          "Successful get data by type",
			user:          "testuser",
			token:         "testtoken",
			dataType:      "testdatatype",
			getDataError:  nil,
			expectedError: nil,
		},
		{
			name:          "Error getting data from storage",
			user:          "testuser",
			token:         "testtoken",
			dataType:      "testdatatype",
			getDataError:  errors.New("storage error"),
			expectedError: errors.New("storage error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.Token = tt.token
			if tt.token != "" {
				clientUseCase.Config.User = tt.user
				clientUseCase.Config.Passphrase = "testpassphrase"
				mockStorage.On("GetData", tt.user).Return([]models.StoredData{}, tt.getDataError)
			}
			_, err = clientUseCase.GetDataByType(tt.dataType)

			assert.Equal(t, tt.expectedError, err)

		})
	}
}

func TestClientUseCase_FullSync(t *testing.T) {
	tests := []struct {
		name                  string
		user                  string
		token                 string
		deleteDataError       error
		listSecretsError      error
		updateLastServerError error
		expectedError         error
	}{
		{
			name:                  "Successful full sync",
			user:                  "testuser",
			token:                 "testtoken",
			deleteDataError:       nil,
			listSecretsError:      nil,
			updateLastServerError: nil,
			expectedError:         nil,
		},
		{
			name:                  "Error deleting data in storage",
			user:                  "testuser",
			token:                 "testtoken",
			deleteDataError:       errors.New("storage error"),
			listSecretsError:      nil,
			updateLastServerError: nil,
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error listing secrets with transporter",
			user:                  "testuser",
			token:                 "testtoken",
			deleteDataError:       nil,
			listSecretsError:      errors.New("transporter error"),
			updateLastServerError: nil,
			expectedError:         errors.New("transporter error"),
		},
		{
			name:                  "Error updating last server updated time",
			user:                  "testuser",
			token:                 "testtoken",
			deleteDataError:       nil,
			listSecretsError:      nil,
			updateLastServerError: errors.New("storage error"),
			expectedError:         errors.New("storage error"),
		},
		{
			name:                  "Error user not logged in",
			user:                  "",
			token:                 "",
			deleteDataError:       nil,
			listSecretsError:      nil,
			updateLastServerError: nil,
			expectedError:         errors.New("user not logged in"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			mockStorage := mocks.NewStorager(t)
			mockTransport := mocks.NewTransporter(t)
			testConfig := config.NewClientConfig()
			clientUseCase := &ClientUseCase{
				Config:      testConfig,
				Storage:     mockStorage,
				Transporter: mockTransport,
			}
			clientUseCase.Config.Token = tt.token
			if tt.token != "" {
				clientUseCase.Config.User = tt.user
				mockStorage.On("DeleteAllData", tt.user).Return(tt.deleteDataError)
				if tt.deleteDataError == nil {
					var protoData []*pb.SecretData
					mockTransport.On("ListSecrets", context.Background()).Return(protoData, tt.listSecretsError)
					if tt.listSecretsError == nil {
						mockStorage.On("UpdateLastServerUpdated", tt.user, int64(0)).Return(tt.updateLastServerError)
					}
				}
			}
			err = clientUseCase.FullSync(context.Background())
			assert.Equal(t, tt.expectedError, err)

		})
	}
}
