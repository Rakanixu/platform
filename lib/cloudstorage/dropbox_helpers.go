package cloudstorage

func (dcs *DropboxCloudStorage) token() string {
	return "Bearer " + dcs.Endpoint.Token.AccessToken
}
