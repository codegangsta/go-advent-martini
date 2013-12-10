package main

import (
	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"labix.org/v2/mgo"
)

type Wish struct {
	Name        string `form:"name"`
	Description string `form:"description"`
}

func DB() martini.Handler {
	session, err := mgo.Dial("mongodb://localhost")
	if err != nil {
		panic(err)
	}

	return func(c martini.Context) {
		s := session.Clone()
		c.Map(s.DB("advent"))

		defer s.Close()
		c.Next()
	}
}

func GetAll(db *mgo.Database) []Wish {
	var wishlist []Wish
	db.C("wishes").Find(nil).All(&wishlist)
	return wishlist
}

func main() {
	m := martini.Classic()
	m.Use(render.Renderer())
	m.Use(DB())

	m.Get("/wishes", func(r render.Render, db *mgo.Database) {
		r.HTML(200, "list", GetAll(db))
	})

	m.Post("/wishes", binding.Form(Wish{}), func(wish Wish, r render.Render, db *mgo.Database) {
		db.C("wishes").Insert(wish)
		r.HTML(200, "list", GetAll(db))
	})

	m.Run()
}
