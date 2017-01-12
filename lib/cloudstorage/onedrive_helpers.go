package cloudstorage

func (ocs *OneDriveCloudStorage) token() string {
	return ocs.Endpoint.Token.TokenType + " " + ocs.Endpoint.Token.AccessToken
}
