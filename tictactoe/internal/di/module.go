package di

import "go.uber.org/fx"

var Module = fx.Module("tictactoe",
	fx.Provide(
		NewGameStorage,
		NewGameRepository,
		
		NewMinimax,
		NewGameService,

		NewGameHandler,
		NewUserService,
		NewAuthHandler,
		NewAuthMiddlewareHandler,
		NewRouter,
	),
	
	fx.Invoke(
		RegisterServer,
	),
)