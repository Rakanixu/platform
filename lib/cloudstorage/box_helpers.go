package cloudstorage

func (bcs *BoxCloudStorage) token() string {
	return "Bearer " + bcs.Endpoint.Token.AccessToken
}
