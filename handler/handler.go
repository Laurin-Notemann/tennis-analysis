package handler

type ResourceHandlers struct {
  UserHandler UserHandler
  TokenHandler RefreshTokenHandler
  AuthHandler AuthenticationHandler
  PlayerHandler PlayerHandler
  TeamHandler TeamHandler
}
