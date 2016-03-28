package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
    fmt.Printf("Creating database 'airship.db'... ")
    
    os.Remove("./airship.db")
	db, err := sql.Open("sqlite3", "./airship.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
    
    fmt.Println("Done!")
    
    //create tables
    createCrew := `create table Crew (Employee_ID integer not null primary key, Annual_Salary float, Name string, Mans_Cannon integer, Fights_Sky_Pirates integer);`
    createCrew_Roles := `create table Crew_Roles (Employee_ID integer, Role string);`
    createCannons := `create table Cannons (Field_of_View string, Floor_Number integer, Crew_Member string);`
    createCannon_Ammo := `create table Cannon_Ammo (Field_of_View string, Floor_Number integer, Ammunition_Type string);`
    createFloors := `create table Floors (Floor_Number integer not null primary key);`
    createCrew_Assigned_Floors := `create table Crew_Assigned_Floors (Employee_ID integer, Floor_Number integer);`
    createGuest_Room := `create table Guest_Room (Room_Number integer not null primary key, Nighly_Rate float, Maximum_Occupancy integer, Floor_Number integer);`
    createPassenger := `create table Passenger (Ticket_Number integer, Name string, Room_Number integer);`
    
    fmt.Printf("Creating tables... ")
	_, err = db.Exec(createCrew)
	if err != nil {
		log.Printf("%q: %s\n", err, createCrew)
		return
	}
    _, err = db.Exec(createCrew_Roles)
	if err != nil {
		log.Printf("%q: %s\n", err, createCrew_Roles)
		return
	}
    _, err = db.Exec(createCannons)
	if err != nil {
		log.Printf("%q: %s\n", err, createCannons)
		return
	}
    _, err = db.Exec(createCannon_Ammo)
	if err != nil {
		log.Printf("%q: %s\n", err, createCannon_Ammo)
		return
	}
    _, err = db.Exec(createFloors)
	if err != nil {
		log.Printf("%q: %s\n", err, createFloors)
		return
	}
    _, err = db.Exec(createCrew_Assigned_Floors)
	if err != nil {
		log.Printf("%q: %s\n", err, createCrew_Assigned_Floors)
		return
	}
    _, err = db.Exec(createGuest_Room)
	if err != nil {
		log.Printf("%q: %s\n", err, createGuest_Room)
		return
	}
    _, err = db.Exec(createPassenger)
	if err != nil {
		log.Printf("%q: %s\n", err, createPassenger)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}    
    tx.Commit()
    
    fmt.Println("Done!")
    
    //insert rows
    
    fmt.Printf("Inserting rows... ")
    
    fmt.Println("Done!")
    

    /*
	_, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	if err != nil {
		log.Fatal(err)
	}

	rows, err = db.Query("select id, name from foo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)
		fmt.Println(id, name)
	}
    */
}
