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
)

type Page struct {
	Title string
	Body  []byte
    QueryResults queryResults
}

type queryResults struct {
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
    fmt.Println("load page " + title)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func basicHandler(w http.ResponseWriter, r *http.Request, title string) {
    fmt.Println("My title is: " + title)
	p := &Page{Title: title, Body: nil}
	renderTemplate(w, title, p)
}

func viewTableHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
    
	if err != nil {
		http.Redirect(w, r, "/view-table/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "table", p)
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

var templates = template.Must(template.ParseFiles("edit.html", "view.html", "home.html",
"view-tables.html", "table.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view|home|view-tables)/([a-zA-Z0-9]*)$")

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
        
        /*columns, _ := rows.Columns()
        count := len(columns)
        values := make([]interface{}, count)
        valuePtrs := make([]interface{}, count)*/
        
        CrewResults := make([]Crew, 0)
        CrewRoleResults := make([]Crew_Role, 0)
        CrewFloorResults := make([]Crew_Assigned_Floor, 0)
        CannonResults := make([]Cannon, 0)
        CannonAmmoResults := make([]Cannon_Ammo, 0)
        FloorResults := make([]Floor, 0)
        GuestRoomResults := make([]Guest_Room, 0)
        PassengerResults := make([]Passenger, 0)
        
        //Make slice of colum headers**************
        
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
                    
                    rows.Scan(&Field_of_View, &Ammunition_Type, &Floor_Number)
                    res := Cannon_Ammo{Field_of_View, Ammunition_Type, Floor_Number}
                    CannonAmmoResults = append(CannonAmmoResults, res)
                case "floors":
                    var Floor_Number int32
                    
                    rows.Scan(&Floor_Number)
                    res := Floor{Floor_Number}
                    FloorResults = append(FloorResults, res)
                case "guest_rooms":
                    var Room_Number, Maximum_Occupancy, Floor_Number int32
                    var Nightly_Rate float64
                    
                    rows.Scan(&Room_Number, &Maximum_Occupancy, &Floor_Number, &Nightly_Rate)
                    res := Guest_Room{Room_Number, Maximum_Occupancy, Floor_Number, Nightly_Rate}
                    GuestRoomResults = append(GuestRoomResults, res)
                case "passengers":
                    var Ticket_Number, Room_Number int32
                    var Name string
                    
                    rows.Scan(&Ticket_Number, &Room_Number, &Name)
                    res := Passenger{Ticket_Number, Room_Number, Name}
                    PassengerResults = append(PassengerResults, res)
                default:
                    panic("table " + title + " does not exist!")
            }
        }
        
        results := queryResults { CrewResults, CrewRoleResults, CrewFloorResults, CannonResults, CannonAmmoResults, FloorResults, GuestRoomResults, PassengerResults }

        /*for rows.Next() {

            for i, _ := range columns {
                valuePtrs[i] = &values[i]
            }

            rows.Scan(valuePtrs...)

            for i, col := range columns {

                var v interface{}

                val := values[i]

                b, ok := val.([]byte)

                if (ok) {
                    v = string(b)
                } else {
                    v = val
                }
                

                fmt.Println(col, v)
            }
            
            break;
        }*/
        
        p := &Page{Title: title, QueryResults: results}
        renderTemplate(w, "view-tables", p)
    }
}

func main() {
    http.HandleFunc("/home/", makeHandler(basicHandler))
    http.HandleFunc("/view-tables/", makeHandler(basicHandler))    
    
	http.HandleFunc("/view/", makeHandler(viewHandler))
    
    http.HandleFunc("/", getTable)

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
    Field_of_View, Ammunition_Type string
    Floor_Number int32
}
type Floor struct {
    Floor_Number int32
}
type Guest_Room struct {
    Room_Number, Maximum_Occupancy, Floor_Number int32
    Nightly_Rate float64
}
type Passenger struct {
    Ticket_Number, Room_Number int32
    Name string
}
