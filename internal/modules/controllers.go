package modules

import (
	"app/internal/infrastructure/responder"
	pC "app/internal/modules/pet/controller"
	sC "app/internal/modules/store/controller"
	uC "app/internal/modules/user/controller"
	"net/http"
	"os"
	customMiddleware "app/internal/infrastructure/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

type Controller struct {
	User  uC.UserControllerer
	Pet   pC.PetControllerer
	Store sC.StoreControllerer
}

func NewController(services *Service, respond responder.Responder) *Controller {
	return &Controller{
		User:  uC.NewUserController(services.User, respond),
		Pet:   pC.NewPetController(services.Pet, respond),
		Store: sC.NewStoreController(services.Store, respond),
	}
}

func (c *Controller) InitRoutesUser() http.Handler {
	r := chi.NewRouter()

	r.Get("/login", c.User.LoginUser)
	r.Post("/", c.User.CreateUser)
	r.Get("/login", c.User.LoginUser)
	r.Get("/logout", c.User.LogoutUser)
	r.Post("/createWithArray", c.User.CreateUsersWithArrayInput)
	r.Post("/createWithList", c.User.CreateUsersWithListInput)

	r.Route("/{username}", func(r chi.Router) {
		r.Use(middleware.Logger)
		r.Get("/", c.User.GetUserByName)
		r.Put("/", c.User.UpdateUser)
		r.Delete("/", c.User.DeleteUser)

	})

	return r
}

func (c *Controller) InitRoutesPet() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// подключаем авторизацию
		signKey := os.Getenv("SIGN_KEY")
		tokenAuth := jwtauth.New("HS256", []byte(signKey), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(customMiddleware.Authenticator)

		r.Post("/{petId}/uploadImage", c.Pet.UploadFile)
		r.Post("/", c.Pet.AddPet)
		r.Put("/", c.Pet.UpdatePet)
		r.Get("/findByStatus", c.Pet.FindPetsByStatus)
		r.Get("/findByTags", c.Pet.FindPetsByTags)
		r.Get("/{petId}", c.Pet.GetPetById)
		r.Post("/{petId}", c.Pet.UpdatePetWithForm)
		r.Delete("/{petId}", c.Pet.DeletePet)
	})

	return r
}

func (c *Controller) InitRoutesStore() http.Handler {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		// подключаем авторизацию
		signKey := os.Getenv("SIGN_KEY")
		tokenAuth := jwtauth.New("HS256", []byte(signKey), nil)
		r.Use(jwtauth.Verifier(tokenAuth))
		//r.Use(jwtauth.Authenticator)
		r.Use(customMiddleware.Authenticator)

		r.Get("/inventory", c.Store.GetInventory)
	})

	r.Post("/order", c.Store.PlaceOrder)
	r.Get("/order/{orderId}", c.Store.GetOrderById)
	r.Delete("/order/{orderId}", c.Store.DeleteOrder)

	return r
}
