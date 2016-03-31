package main

import (
    "database/sql"
	"html/template"
	"io/ioutil"
	"net/http"
	"regexp"
    "fmt"
    _ "github.com/mattn/go-sqlite3"
    "log"
    "strconv"
)

type Page struct {
	Title string
	Body  []byte
    
    QueryResults queryResults
    
    InsertMessage string
}

type queryResults struct {
    ColumnHeaders []string
    
    CrewResults []Crew
    CrewRoleResults []Crew_Role
    CrewFloorResults []Crew_Assigned_Floor
    CannonResults []Cannon
    CannonAmmoResults []Cannon_Ammo
    FloorResults []Floor
    GuestRoomResults []Guest_Room
    PassengerResults []Passenger
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func basicHandler(w http.ResponseWriter, r *http.Request, title string) {
	p := &Page{Title: title, Body: nil}
	renderTemplate(w, title, p)
}

func viewTableHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
    
	if err != nil {
		http.Redirect(w, r, "/view-table/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view-table", p)
}

func queryTablesHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
    
	if err != nil {
		http.Redirect(w, r, "/query-tables/", http.StatusFound)
		return
	}
	renderTemplate(w, "query-tables", p)
}

var templates = template.Must(template.ParseFiles("home.html",
"view-tables.html", "query-tables.html", "insert-into-tables.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(home|view-tables|query-tables|insert-into-tables)/([a-zA-Z0-9]*)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/" {
            http.Redirect(w, r, "/home/", 301)
        }
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        if (m[2] != "") {
            fn(w, r, m[2]) 
        } else {
            fn(w, r, m[1])
        }      
	}
}

//view-tables
func getTable(w http.ResponseWriter, r *http.Request) {
    if r.Method == "POST" {
        r.ParseForm()
        
        title := r.Form["tables"][0]
        
        db, err := sql.Open("sqlite3", "./airship.db")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()
        rows, err := db.Query("select * from " + title)
        
        columns, _ := rows.Columns()
        
        CrewResults := make([]Crew, 0)
        CrewRoleResults := make([]Crew_Role, 0)
        CrewFloorResults := make([]Crew_Assigned_Floor, 0)
        CannonResults := make([]Cannon, 0)
        CannonAmmoResults := make([]Cannon_Ammo, 0)
        FloorResults := make([]Floor, 0)
        GuestRoomResults := make([]Guest_Room, 0)
        PassengerResults := make([]Passenger, 0)
        
        defer rows.Close()
        for rows.Next() {
            switch title {
                case "crew":
                    var Employee_ID, Mans_Cannon, Fights_Sky_Pirates int32
                    var Name string
                    var Annual_Salary float64
                    
                    rows.Scan(&Employee_ID, &Annual_Salary, &Name, &Mans_Cannon, &Fights_Sky_Pirates)
                    res := Crew{Employee_ID, Mans_Cannon, Fights_Sky_Pirates, Name, Annual_Salary}
                    CrewResults = append(CrewResults, res)
                case "crew_roles":
                    var Employee_ID int32
                    var Role string
                    
                    rows.Scan(&Employee_ID, &Role)
                    res := Crew_Role{Employee_ID, Role}
                    CrewRoleResults = append(CrewRoleResults, res)
                case "crew_assigned_floors":
                    var Employee_ID, Floor_Number int32
                    
                    rows.Scan(&Employee_ID, &Floor_Number)
                    res := Crew_Assigned_Floor{Employee_ID, Floor_Number}
                    CrewFloorResults = append(CrewFloorResults, res)
                case "cannons":
                    var Field_of_View string
                    var Floor_Number, Crew_Member int32
                    
                    rows.Scan(&Field_of_View, &Floor_Number, &Crew_Member)
                    res := Cannon{Field_of_View, Floor_Number, Crew_Member}
                    CannonResults = append(CannonResults, res)
                case "cannon_ammo":
                    var Field_of_View, Ammunition_Type string
                    var Floor_Number int32
                    
                    rows.Scan(&Field_of_View, &Floor_Number, &Ammunition_Type)
                    res := Cannon_Ammo{Field_of_View, Floor_Number, Ammunition_Type}
                    CannonAmmoResults = append(CannonAmmoResults, res)
                case "floors":
                    var Floor_Number int32
                    
                    rows.Scan(&Floor_Number)
                    res := Floor{Floor_Number}
                    FloorResults = append(FloorResults, res)
                case "guest_rooms":
                    var Room_Number, Maximum_Occupancy, Floor_Number int32
                    var Nightly_Rate float64
                    
                    rows.Scan(&Room_Number, &Nightly_Rate, &Maximum_Occupancy, &Floor_Number)
                    res := Guest_Room{Room_Number, Nightly_Rate, Maximum_Occupancy, Floor_Number}
                    GuestRoomResults = append(GuestRoomResults, res)
                case "passengers":
                    var Ticket_Number, Room_Number int32
                    var Name string
                    
                    rows.Scan(&Ticket_Number, &Name, &Room_Number)
                    res := Passenger{Ticket_Number, Name, Room_Number}
                    PassengerResults = append(PassengerResults, res)
                default:
                    panic("table " + title + " does not exist!")
            }
        }
        
        results := queryResults { columns, CrewResults, CrewRoleResults, CrewFloorResults, CannonResults, CannonAmmoResults, FloorResults, GuestRoomResults, PassengerResults }
        
        p := &Page{Title: title, QueryResults: results}
        renderTemplate(w, "view-tables", p)
    } else if r.Method == "GET" {
        
    }
}

//insert-into-tables
func insertIntoTable(w http.ResponseWriter, r *http.Request) {
    fmt.Println(r.Method)
    if r.Method == "POST" {
        r.ParseForm()
        
        roomNum := r.Form["roomNum"][0]
        passName := r.Form["passName"][0]
        
        db, err := sql.Open("sqlite3", "file:airship.db?cache=shared&mode=rwc")
        if err != nil {
            log.Fatal(err)
        }
        defer db.Close()
        
        InsertMessage := ""
        
        fmt.Println(passName)
        
        PassengerResults := make([]Passenger, 0)
        passengerRows, err := db.Query("select * from passengers where room_number = " + roomNum)
        defer passengerRows.Close()
        for passengerRows.Next() {
            var Ticket_Number, Room_Number int32
            var Name string
            
            passengerRows.Scan(&Ticket_Number, &Name, &Room_Number)
            res := Passenger{Ticket_Number, Name, Room_Number}
            PassengerResults = append(PassengerResults, res)
        }
        
        PassCResults := make([]Passenger, 0)
        passCRows, err := db.Query("select * from passengers")
        defer passCRows.Close()
        for passCRows.Next() {
            var Ticket_Number, Room_Number int32
            var Name string
            
            passCRows.Scan(&Ticket_Number, &Name, &Room_Number)
            res := Passenger{Ticket_Number, Name, Room_Number}
            PassCResults = append(PassCResults, res)
        }
        passengerCount := len(PassCResults)
        
        var maxOcc int32
        roomRows, err := db.Query("select Maximum_Occupancy from guest_rooms where room_number = " + roomNum)
        defer roomRows.Close()
        for roomRows.Next() {
            var Maximum_Occupancy int32
            
            roomRows.Scan(&Maximum_Occupancy)
            maxOcc = Maximum_Occupancy
            
            break;
        }
        
        if maxOcc == int32(len(PassengerResults)) {
            InsertMessage = "Cannot add! Room is filled!"
        } else {
            ticketNum := passengerCount + 1
            _, err = db.Exec("insert into Passengers(Ticket_Number, Name, Room_Number) values (" + strconv.Itoa(ticketNum) + ", '" + passName + "', " + roomNum + ")")
            if err != nil {
                log.Fatal(err)
                InsertMessage = "Failed to insert!"
            } else {
                InsertMessage = "Successfully inserted!"
            }
        }
        
        p := &Page{ InsertMessage: InsertMessage }
        renderTemplate(w, "insert-into-tables", p)
        //http.Redirect(w, r, "/insert-into-tables/", 301)
    } else if r.Method == "GET" {
        
    }
}

func main() {
    fs := http.FileServer(http.Dir("css"))
    http.Handle("/css/", http.StripPrefix("/css/", fs))
    
    http.HandleFunc("/home/", makeHandler(basicHandler))
    http.HandleFunc("/view-tables/", makeHandler(basicHandler))
    http.HandleFunc("/insert-into-tables/", makeHandler(basicHandler))
    http.HandleFunc("/query-tables/", makeHandler(basicHandler))  
    
    http.HandleFunc("/view-tables/getTable", getTable)
    http.HandleFunc("/", insertIntoTable)

    fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", nil)
}

//Database models
type Crew struct {
    Employee_ID, Mans_Cannon, Fights_Sky_Pirates int32
    Name string
    Annual_Salary float64
}
type Crew_Role struct {
    Employee_ID int32
    Role string
}
type Crew_Assigned_Floor struct {
    Employee_ID, Floor_Number int32
}
type Cannon struct {
    Field_of_View string
    Floor_Number, Crew_Member int32
}
type Cannon_Ammo struct {
    Field_of_View string
    Floor_Number int32
    Ammunition_Type string
}
type Floor struct {
    Floor_Number int32
}
type Guest_Room struct {
    Room_Number int32
    Nightly_Rate float64
    Maximum_Occupancy, Floor_Number int32
}
type Passenger struct {
    Ticket_Number int32
    Name string
    Room_Number int32
}
