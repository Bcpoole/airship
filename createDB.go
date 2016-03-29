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
    createCrew_Roles := `create table Crew_Roles (Employee_ID integer references Crew(Employee_ID), Role string);`
    createCannons := `create table Cannons (Field_of_View string, Floor_Number integer references Floors(Floor_Number), Crew_Member integer references Crew(Employee_ID));`
    createCannon_Ammo := `create table Cannon_Ammo (Field_of_View string references Cannons(Field_of_View), Floor_Number integer references Floors(Floor_Number), Ammunition_Type string);`
    createFloors := `create table Floors (Floor_Number integer not null primary key);`
    createCrew_Assigned_Floors := `create table Crew_Assigned_Floors (Employee_ID integer references Crew(Employee_ID), Floor_Number integer references Floors(Floor_Number));`
    createGuest_Rooms := `create table Guest_Rooms (Room_Number integer not null primary key, Nighly_Rate float, Maximum_Occupancy integer, Floor_Number integer references Floors(Floor_Number));`
    createPassengers := `create table Passengers (Ticket_Number integer, Name string, Room_Number integer references Guest_Rooms(Room_Number));`
    
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
    _, err = db.Exec(createGuest_Rooms)
	if err != nil {
		log.Printf("%q: %s\n", err, createGuest_Rooms)
		return
	}
    _, err = db.Exec(createPassengers)
	if err != nil {
		log.Printf("%q: %s\n", err, createPassengers)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}    
    tx.Commit()
    
    fmt.Println("Done!")
    
    //insert rows
    fmt.Println("Inserting rows into... ")
    
    fmt.Println("\tCrew")
	_, err = db.Exec("insert into Crew(Employee_ID, Annual_Salary, Name, Mans_Cannon, Fights_Sky_Pirates) values " +
    "(1501, 40000.00, 'Paul Goodman', 0, 0), " +
    "(1502, 40000.00, 'Justin Malkai', 0, 1), " +
    "(1503, 65000.50, 'Planters Peanut', 1, 1), " +
    "(7707, 808000.88, 'Rick Sanchez', 1, 1), " +
    "(7709, 120000.00, 'Joesph Joestar', 0, 1), " +
    "(7710, 89000.00, 'Gunguir', 1, 0), " +
    "(404, 40440.40, 'Spinzaku Kururugi', 1, 1)")
	if err != nil {
		log.Fatal(err)
	}
    
    
    fmt.Println("\tCrew_Roles")
	_, err = db.Exec("insert into Crew_Roles(Employee_ID, Role) values " +
    "(1501, 'Janitor'), " +
    "(1502, 'Cook'), " +
    "(1503, 'Doctor'), " +
    "(7707, 'Scientist'), " +
    "(7707, 'Brewer'), " +
    "(7709, 'Yoga Instructor'), " +
    "(7709, 'Vampire Slayer'), " +
    "(7709, 'Captain'), " +
    "(7710, 'Gunner'), " +
    "(404, 'DDR Instructor'), " +
    "(404, 'CQC Guard')")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tCannons")
	_, err = db.Exec("insert into Cannons(Field_of_View, Floor_Number, Crew_Member) values " +
    "('Bow', 3, 7710), " +
    "('Port', 2, 404), " +
    "('Starboard', 2, 1503), " +
    "('Bow', 1, 7707)")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tCannon_Ammo")
	_, err = db.Exec("insert into Cannon_Ammo(Field_of_View, Floor_Number, Ammunition_Type) values " +
    "('Bow', 3, 'Heated Shot'), " +
    "('Bow', 1, 'Grapeshot'), " +
    "('Port', 2, 'Round Shot'), " +
    "('Starboard', 2, 'Round Shot'), " +
    "('Port', 2, 'Shell'), " +
    "('Starboard', 2, 'Shell'), " +
    "('Bow', 1, 'Spider Shot'), " +
    "('Bow', 3, 'Canister Shot'), " +
    "('Bow', 1, 'Shrapnel')")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tFloors")
	_, err = db.Exec("insert into Floors(Floor_Number) values " +
    "(1), " +
    "(2), " +
    "(3)")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tCrew_Assigned_Floors")
	_, err = db.Exec("insert into Crew_Assigned_Floors(Employee_ID, Floor_Number) values " +
    "(1501, 1), " +
    "(1501, 2), " +
    "(1501, 3), " +
    "(1502, 1), " +
    "(1503, 2), " +
    "(7707, 1), " +
    "(7707, 3), " +
    "(7709, 1), " +
    "(7710, 1), " +
    "(404, 1), " +
    "(404, 2)")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tGuest_Rooms")
	_, err = db.Exec("insert into Guest_Rooms(Room_Number, Nighly_Rate, Maximum_Occupancy, Floor_Number) values " +
    "(200, 179.99, 2, 2), " +
    "(201, 179.99, 2, 2), " +
    "(202, 179.99, 2, 2), " +
    "(203, 179.99, 2, 2), " +
    "(204, 239.99, 4, 2), " +
    "(205, 239.99, 4, 2), " +
    "(100, 279.99, 4, 1), " +
    "(101, 279.99, 4, 1), " +
    "(102, 299.99, 4, 1), " +
    "(103, 299.99, 4, 1), " +
    "(0, 399.99, 2, 1)")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("\tPassengers")
	_, err = db.Exec("insert into Passengers(Ticket_Number, Name, Room_Number) values " +
    "(1, 'William Middleton', 0), " +
    "(2, 'Catherine Middleton', 0), " +
    "(3, 'Harry Houdini', 203), " +
    "(4, 'Dio Brando', 103), " +
    "(5, 'John Walton, Sr.', 204), " +
    "(6, 'John Walton, Jr.', 204), " +
    "(7, 'Olivia Walton', 204), " +
    "(8, 'Jason Walton', 204), " +
    "(9, 'Zebulon Walton', 205), " +
    "(10, 'Esther Walton', 205), " +
    "(11, 'Cindy Walton', 205), " +
    "(12, 'Rose Walton', 205)")
	if err != nil {
		log.Fatal(err)
	}
    
    fmt.Println("Done!")
}
