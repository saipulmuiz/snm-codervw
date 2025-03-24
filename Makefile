generate-mocks:
	# repositories
	@mockgen -destination=./service/repository/mocks/mock_user_repository.go -package=mocks codepair-sinarmas/service UserRepository
	@mockgen -destination=./service/repository/mocks/mock_otp_repository.go -package=mocks codepair-sinarmas/service OtpRepository