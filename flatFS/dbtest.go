package FlatFS

import "log"

func testFunc() {
	const dbpath = "foo.db"

	db := InitDB(dbpath)
	defer db.Close()
	CreateTable(db)

	items := []TestItem{
		TestItem{"10", "A", "213"},
		TestItem{"15", "B", "214"},
	}
	StoreItem(db, items)

	readItems := ReadItem(db)
	log.Println(readItems)
	items2 := []TestItem{
		TestItem{"1", "C", "215"},
		TestItem{"3", "D", "216"},
	}
	StoreItem(db, items2)
	readItems2 := ReadItem(db)
	log.Println(readItems2)
}
