package cloudstorage

func (scs *SlackCloudStorage) token() string {
	return "Bearer " + scs.Endpoint.Token.AccessToken
}
