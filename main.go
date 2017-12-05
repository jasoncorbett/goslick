package main

import (
	"github.com/jasoncorbett/goslick/com_slickqa"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"fmt"
)

func main() {
	session, err := mgo.Dial("localhost")
	projects := session.DB("slick").C("projects")
	//err = projects.Insert(&com_slickqa.Project{Name: "bar"})
	//if err != nil {
        //        fmt.Println(err)
        //}

	result := com_slickqa.Project{}
	err = projects.Find(bson.M{"name": "bar"}).One(&result)
	if err != nil {
                fmt.Println(err)
	}
	//result.Id = fmt.Sprintf("%x", result.Id)

	fmt.Println(result)
}
